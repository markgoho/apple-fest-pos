package pos

import (
	"database/sql"
	"errors"
	"fmt"
)

const transactionColumns = `id, client_order_id, device_id, order_number, business_date, status,
	subtotal_cents, tax_cents, total_cents, payment_method, request_json,
	customer_print_status, kitchen_print_status, created_at`

func (service *OrderService) findByClientOrderID(clientOrderID string) (transactionRow, bool, error) {
	var row transactionRow
	err := service.DB.
		QueryRow(`SELECT `+transactionColumns+` FROM transactions WHERE client_order_id = ?`, clientOrderID).
		Scan(&row.ID, &row.ClientOrderID, &row.DeviceID, &row.OrderNumber, &row.BusinessDate,
			&row.Status, &row.SubtotalCents, &row.TaxCents, &row.TotalCents, &row.PaymentMethod,
			&row.RequestJSON, &row.CustomerPrintStatus, &row.KitchenPrintStatus, &row.CreatedAt)

	if errors.Is(err, sql.ErrNoRows) {
		return transactionRow{}, false, nil
	}
	if err != nil {
		return transactionRow{}, false, fmt.Errorf("read order by client id: %w", err)
	}
	return row, true, nil
}

// insertOrder takes the next order number and writes the row in one
// transaction. Every read and write inside it goes through tx, because the
// pool holds one connection and a db call here would deadlock.
func (service *OrderService) insertOrder(request PlaceOrderRequest) (transactionRow, error) {
	createdAt := service.Now().UTC().Format(timestampLayout)
	businessDate := createdAt[:10]

	subtotalCents := 0
	for _, line := range request.Items {
		item, _ := MenuItemByID(line.MenuItemID)
		subtotalCents += item.PriceCents * line.Quantity
	}

	row := transactionRow{
		ID:                  newUUID(),
		ClientOrderID:       request.ClientOrderID,
		DeviceID:            request.DeviceID,
		BusinessDate:        businessDate,
		Status:              OrderAccepted,
		SubtotalCents:       subtotalCents,
		TaxCents:            0,
		TotalCents:          subtotalCents,
		PaymentMethod:       request.Payment.Method,
		RequestJSON:         mustMarshal(request),
		CustomerPrintStatus: PrintQueued,
		KitchenPrintStatus:  PrintQueued,
		CreatedAt:           createdAt,
	}

	transaction, err := service.DB.Begin()
	if err != nil {
		return transactionRow{}, fmt.Errorf("begin order transaction: %w", err)
	}
	defer transaction.Rollback()

	storedBusinessDate, err := metadataValue(transaction, "business_date", businessDate)
	if err != nil {
		return transactionRow{}, err
	}
	nextOrderNumber, err := metadataValue(transaction, "next_order_number", "")
	if err != nil {
		return transactionRow{}, err
	}

	row.OrderNumber = service.StartingOrderNumber
	if storedBusinessDate == businessDate {
		if parsed, err := parseOrderNumber(nextOrderNumber); err == nil {
			row.OrderNumber = parsed
		}
	}

	if err := setMetadataValue(transaction, "business_date", businessDate); err != nil {
		return transactionRow{}, err
	}
	if err := setMetadataValue(transaction, "next_order_number", fmt.Sprint(row.OrderNumber+1)); err != nil {
		return transactionRow{}, err
	}

	if _, err := transaction.Exec(
		`INSERT INTO transactions (`+transactionColumns+`)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		row.ID, row.ClientOrderID, row.DeviceID, row.OrderNumber, row.BusinessDate,
		string(row.Status), row.SubtotalCents, row.TaxCents, row.TotalCents, row.PaymentMethod,
		row.RequestJSON, string(row.CustomerPrintStatus), string(row.KitchenPrintStatus), row.CreatedAt,
	); err != nil {
		return transactionRow{}, fmt.Errorf("insert order: %w", err)
	}

	if err := transaction.Commit(); err != nil {
		return transactionRow{}, fmt.Errorf("commit order: %w", err)
	}
	return row, nil
}

func metadataValue(transaction *sql.Tx, key string, fallback string) (string, error) {
	var value string
	err := transaction.QueryRow(`SELECT value FROM metadata WHERE key = ?`, key).Scan(&value)
	if errors.Is(err, sql.ErrNoRows) {
		return fallback, nil
	}
	if err != nil {
		return "", fmt.Errorf("read metadata %q: %w", key, err)
	}
	return value, nil
}

func setMetadataValue(transaction *sql.Tx, key string, value string) error {
	_, err := transaction.Exec(
		`INSERT INTO metadata (key, value) VALUES (?, ?)
		 ON CONFLICT(key) DO UPDATE SET value = excluded.value`, key, value)
	if err != nil {
		return fmt.Errorf("write metadata %q: %w", key, err)
	}
	return nil
}
