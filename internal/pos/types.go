package pos

import (
	"encoding/json"
	"strconv"
)

// PrintStatus is the result of one print attempt.
type PrintStatus string

const (
	PrintQueued   PrintStatus = "queued"
	PrintSent     PrintStatus = "sent"
	PrintFailed   PrintStatus = "failed"
	PrintDisabled PrintStatus = "disabled"
)

// OrderStatus is the state of a placed order.
type OrderStatus string

const (
	OrderAccepted    OrderStatus = "accepted"
	OrderPrinted     OrderStatus = "printed"
	OrderPrintFailed OrderStatus = "print_failed"
)

// CartLine is one line of a cart.
type CartLine struct {
	MenuItemID string `json:"menuItemId"`
	Quantity   int    `json:"quantity"`
	Side       string `json:"side,omitempty"`
	Notes      string `json:"notes,omitempty"`

	// fractionalQuantity is true when the request sent a quantity that is not
	// a whole number. The zero value therefore keeps hand-built lines valid.
	fractionalQuantity bool
}

// Payment holds the payment method. Only cash is supported.
type Payment struct {
	Method string `json:"method"`
}

// PlaceOrderRequest is the body of POST /api/orders.
type PlaceOrderRequest struct {
	ClientOrderID string     `json:"clientOrderId"`
	DeviceID      string     `json:"deviceId"`
	Payment       Payment    `json:"payment"`
	Notes         string     `json:"notes,omitempty"`
	Items         []CartLine `json:"items"`
}

// PlacedOrder is the order part of the place-order response.
type PlacedOrder struct {
	ID            string      `json:"id"`
	OrderNumber   int         `json:"orderNumber"`
	Status        OrderStatus `json:"status"`
	SubtotalCents int         `json:"subtotalCents"`
	TaxCents      int         `json:"taxCents"`
	TotalCents    int         `json:"totalCents"`
	CreatedAt     string      `json:"createdAt"`
}

// PrintResult holds the status of both printers.
type PrintResult struct {
	Customer PrintStatus `json:"customer"`
	Kitchen  PrintStatus `json:"kitchen"`
}

// PlaceOrderResponse is the body of a successful POST /api/orders.
type PlaceOrderResponse struct {
	Order PlacedOrder `json:"order"`
	Print PrintResult `json:"print"`
}

// ReceiptOrder is what the ESC/POS builders need.
type ReceiptOrder struct {
	OrderID       string
	OrderNumber   int
	CreatedAt     string
	SubtotalCents int
	TotalCents    int
	Items         []CartLine
}

// cartLineJSON mirrors CartLine but takes the quantity as a raw JSON number,
// so a fractional quantity gives the validation message instead of a decode
// error. The TypeScript server checked Number.isInteger for the same reason.
type cartLineJSON struct {
	MenuItemID string      `json:"menuItemId"`
	Quantity   json.Number `json:"quantity"`
	Side       string      `json:"side,omitempty"`
	Notes      string      `json:"notes,omitempty"`
}

// UnmarshalJSON reads a cart line and marks a non-integer quantity as invalid.
func (line *CartLine) UnmarshalJSON(data []byte) error {
	var raw cartLineJSON
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	line.MenuItemID = raw.MenuItemID
	line.Side = raw.Side
	line.Notes = raw.Notes
	quantity, err := strconv.ParseInt(raw.Quantity.String(), 10, 32)
	line.Quantity = int(quantity)
	line.fractionalQuantity = err != nil

	return nil
}
