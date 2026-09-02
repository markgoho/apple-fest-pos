package pos

import (
	"crypto/rand"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

// timestampLayout matches the JavaScript Date.toISOString() output exactly:
// UTC, three decimal places, a literal Z. The business date is the first ten
// characters of it, so the format must never drift.
const timestampLayout = "2006-01-02T15:04:05.000Z07:00"

// OrderService places orders and reads them back.
type OrderService struct {
	DB                  *sql.DB
	Printer             PrinterConfig
	StartingOrderNumber int
	Now                 func() time.Time
}

// ErrValidation marks a request the operator can correct. It maps to 400.
var ErrValidation = errors.New("validation")

type transactionRow struct {
	ID                  string
	ClientOrderID       string
	DeviceID            string
	OrderNumber         int
	BusinessDate        string
	Status              OrderStatus
	SubtotalCents       int
	TaxCents            int
	TotalCents          int
	PaymentMethod       string
	RequestJSON         string
	CustomerPrintStatus PrintStatus
	KitchenPrintStatus  PrintStatus
	CreatedAt           string
}

// PlaceOrder validates, stores, and prints one order. A repeated
// clientOrderId replays the stored order instead of selling twice.
func (service *OrderService) PlaceOrder(request PlaceOrderRequest) (PlaceOrderResponse, error) {
	if err := ValidateOrder(request); err != nil {
		return PlaceOrderResponse{}, err
	}

	if existing, found, err := service.findByClientOrderID(request.ClientOrderID); err != nil {
		return PlaceOrderResponse{}, err
	} else if found {
		return existing.response(existing.Status, PrintResult{
			Customer: existing.CustomerPrintStatus,
			Kitchen:  existing.KitchenPrintStatus,
		}), nil
	}

	row, err := service.insertOrder(request)
	if err != nil {
		return PlaceOrderResponse{}, err
	}

	print := PrintOrder(service.Printer, ReceiptOrder{
		OrderID:       row.ID,
		OrderNumber:   row.OrderNumber,
		CreatedAt:     row.CreatedAt,
		SubtotalCents: row.SubtotalCents,
		TotalCents:    row.TotalCents,
		Items:         request.Items,
	})
	status := orderStatusFor(print)

	if _, err := service.DB.Exec(
		`UPDATE transactions
		 SET status = ?, customer_print_status = ?, kitchen_print_status = ?
		 WHERE id = ?`,
		string(status), string(print.Customer), string(print.Kitchen), row.ID,
	); err != nil {
		return PlaceOrderResponse{}, fmt.Errorf("update print status: %w", err)
	}

	return row.response(status, print), nil
}

func (row transactionRow) response(status OrderStatus, print PrintResult) PlaceOrderResponse {
	return PlaceOrderResponse{
		Order: PlacedOrder{
			ID:            row.ID,
			OrderNumber:   row.OrderNumber,
			Status:        status,
			SubtotalCents: row.SubtotalCents,
			TaxCents:      row.TaxCents,
			TotalCents:    row.TotalCents,
			CreatedAt:     row.CreatedAt,
		},
		Print: print,
	}
}

func orderStatusFor(print PrintResult) OrderStatus {
	if print.Customer == PrintFailed || print.Kitchen == PrintFailed {
		return OrderPrintFailed
	}
	if print.Customer == PrintPrinted && print.Kitchen == PrintPrinted {
		return OrderPrinted
	}
	return OrderAccepted
}

// ValidateOrder checks the request in the same order as the TypeScript server,
// because the error message tells the operator what to correct.
func ValidateOrder(request PlaceOrderRequest) error {
	if request.ClientOrderID == "" || request.DeviceID == "" {
		return fmt.Errorf("%w: Missing order identifiers", ErrValidation)
	}
	if request.Payment.Method != "cash" {
		return fmt.Errorf("%w: Only cash payments are supported", ErrValidation)
	}
	if len(request.Items) == 0 {
		return fmt.Errorf("%w: Order must contain at least one item", ErrValidation)
	}
	if len(request.Notes) > 500 {
		return fmt.Errorf("%w: Order notes are too long", ErrValidation)
	}

	for _, line := range request.Items {
		if line.fractionalQuantity || line.Quantity < 1 {
			return fmt.Errorf("%w: Item quantity must be a positive integer", ErrValidation)
		}
		if len(line.Notes) > 200 {
			return fmt.Errorf("%w: Item notes are too long", ErrValidation)
		}
		if _, found := MenuItemByID(line.MenuItemID); !found {
			return fmt.Errorf("%w: Unknown menu item: %s", ErrValidation, line.MenuItemID)
		}
	}

	return nil
}

// newUUID makes a random version 4 UUID. The standard library has no UUID
// type, and the identifier is opaque, so ten lines beat a dependency.
func newUUID() string {
	var bytes [16]byte
	if _, err := rand.Read(bytes[:]); err != nil {
		panic(err)
	}
	bytes[6] = (bytes[6] & 0x0f) | 0x40
	bytes[8] = (bytes[8] & 0x3f) | 0x80
	return fmt.Sprintf("%x-%x-%x-%x-%x", bytes[0:4], bytes[4:6], bytes[6:8], bytes[8:10], bytes[10:16])
}

func mustMarshal(value any) string {
	encoded, err := json.Marshal(value)
	if err != nil {
		panic(err)
	}
	return string(encoded)
}
