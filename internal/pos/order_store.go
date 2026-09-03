package pos

import (
	"database/sql"
	"errors"
	"fmt"
)

const transactionColumns = `id, client_order_id, device_id, order_number, business_date, status,
	subtotal_cents, tax_cents, total_cents, payment_method, request_json,
	customer_print_status, kitchen_print_status, created_at`

// findByClientOrderID reads a stored order inside the open transaction, so
// that the check and the insert are one atomic step. A double tap on the
// tablet must replay the order, not hit the UNIQUE constraint.
func findByClientOrderID(transaction *sql.Tx, clientOrderID string) (transactionRow, bool, error) {
	row := transaction.QueryRow(`SELECT `+transactionColumns+` FROM transactions WHERE client_order_id = ?`, clientOrderID)
	return scanTransactionRow(row)
}

// findByID reads a stored order by its id, for the reprint endpoint. It is a
// plain read outside any transaction: a reprint carries no dedup requirement
// the way placing an order does.
func findByID(database *sql.DB, id string) (transactionRow, bool, error) {
	row := database.QueryRow(`SELECT `+transactionColumns+` FROM transactions WHERE id = ?`, id)
	return scanTransactionRow(row)
}

func scanTransactionRow(row *sql.Row) (transactionRow, bool, error) {
	var transaction transactionRow
	err := row.Scan(&transaction.ID, &transaction.ClientOrderID, &transaction.DeviceID, &transaction.OrderNumber,
		&transaction.BusinessDate, &transaction.Status, &transaction.SubtotalCents, &transaction.TaxCents,
		&transaction.TotalCents, &transaction.PaymentMethod, &transaction.RequestJSON,
		&transaction.CustomerPrintStatus, &transaction.KitchenPrintStatus, &transaction.CreatedAt)

	if errors.Is(err, sql.ErrNoRows) {
		return transactionRow{}, false, nil
	}
	if err != nil {
		return transactionRow{}, false, fmt.Errorf("read order: %w", err)
	}
	return transaction, true, nil
}

// insertOrder takes the next order number and writes the row in one
// transaction. Every read and write inside it goes through tx, because the
// pool holds one connection and a db call here would deadlock.
func (service *OrderService) insertOrder(request PlaceOrderRequest) (transactionRow, bool, error) {
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
		return transactionRow{}, false, fmt.Errorf("begin order transaction: %w", err)
	}
	defer transaction.Rollback()

	existing, found, err := findByClientOrderID(transaction, request.ClientOrderID)
	if err != nil {
		return transactionRow{}, false, err
	}
	if found {
		return existing, true, nil
	}

	storedBusinessDate, err := metadataValue(transaction, "business_date", businessDate)
	if err != nil {
		return transactionRow{}, false, err
	}
	nextOrderNumber, err := metadataValue(transaction, "next_order_number", "")
	if err != nil {
		return transactionRow{}, false, err
	}

	row.OrderNumber = service.StartingOrderNumber
	if storedBusinessDate == businessDate {
		if parsed, err := parseOrderNumber(nextOrderNumber); err == nil {
			row.OrderNumber = parsed
		}
	}

	if err := setMetadataValue(transaction, "business_date", businessDate); err != nil {
		return transactionRow{}, false, err
	}
	if err := setMetadataValue(transaction, "next_order_number", fmt.Sprint(row.OrderNumber+1)); err != nil {
		return transactionRow{}, false, err
	}

	if _, err := transaction.Exec(
		`INSERT INTO transactions (`+transactionColumns+`)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		row.ID, row.ClientOrderID, row.DeviceID, row.OrderNumber, row.BusinessDate,
		string(row.Status), row.SubtotalCents, row.TaxCents, row.TotalCents, row.PaymentMethod,
		row.RequestJSON, string(row.CustomerPrintStatus), string(row.KitchenPrintStatus), row.CreatedAt,
	); err != nil {
		return transactionRow{}, false, fmt.Errorf("insert order: %w", err)
	}

	if err := transaction.Commit(); err != nil {
		return transactionRow{}, false, fmt.Errorf("commit order: %w", err)
	}
	return row, false, nil
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
