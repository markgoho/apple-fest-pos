package pos

import (
	"bytes"
	"strings"
	"testing"
)

var receiptOrder = ReceiptOrder{
	OrderID:       "order-1",
	OrderNumber:   101,
	CreatedAt:     "2026-05-07T12:34:00.000Z",
	SubtotalCents: 2000,
	TotalCents:    2000,
	Items:         []CartLine{{MenuItemID: "potato-pancake", Quantity: 2}},
}

func TestBuildCustomerReceiptHasTheCutCommand(t *testing.T) {
	payload := BuildCustomerReceipt(receiptOrder, HeaderNone)

	if !bytes.HasPrefix(payload, []byte{0x1b, 0x40}) {
		t.Errorf("payload does not start with the initialize command")
	}
	if !bytes.HasSuffix(payload, []byte{0x1d, 0x56, 0x42, 0x08}) {
		t.Errorf("payload does not end with the cut command")
	}
	if !strings.Contains(string(payload), "\r\n") {
		t.Errorf("payload has no CRLF line endings")
	}
	if !strings.Contains(string(payload), "Total $20.00") {
		t.Errorf("payload has no formatted total: %q", string(payload))
	}
	if strings.Contains(string(payload), "REPRINT") {
		t.Errorf("a first print should not carry a REPRINT header")
	}
}

func TestBuildCustomerReceiptReprintHasTheReprintHeader(t *testing.T) {
	payload := BuildCustomerReceipt(receiptOrder, HeaderReprint)

	if !strings.Contains(string(payload), "REPRINT") {
		t.Errorf("reprint payload has no REPRINT header: %q", string(payload))
	}
}

func TestBuildCustomerReceiptTestHasTheTestHeaderAndNoOrderNumber(t *testing.T) {
	payload := BuildCustomerReceipt(receiptOrder, HeaderTest)

	if !strings.Contains(string(payload), "TEST") {
		t.Errorf("test payload has no TEST header: %q", string(payload))
	}
	if strings.Contains(string(payload), "Order #") {
		t.Errorf("test payload should carry no order number: %q", string(payload))
	}
}

func TestBuildKitchenTicketHasEmphasisAndTheCutCommand(t *testing.T) {
	payload := BuildKitchenTicket(receiptOrder, HeaderNone)

	if !bytes.HasPrefix(payload, []byte{0x1b, 0x40}) {
		t.Errorf("payload does not start with the initialize command")
	}
	if !bytes.HasSuffix(payload, []byte{0x1d, 0x56, 0x42, 0x08}) {
		t.Errorf("payload does not end with the cut command")
	}
	if !bytes.Contains(payload, []byte{0x1d, 0x21, 0x11}) {
		t.Errorf("payload has no double size command")
	}
	if !strings.Contains(string(payload), "2  POTATO PANCAKE") {
		t.Errorf("payload has no upper case item line: %q", string(payload))
	}
	if strings.Contains(string(payload), "REPRINT") {
		t.Errorf("a first print should not carry a REPRINT header")
	}
}

func TestBuildKitchenTicketReprintHasTheReprintHeader(t *testing.T) {
	payload := BuildKitchenTicket(receiptOrder, HeaderReprint)

	if !strings.Contains(string(payload), "REPRINT") {
		t.Errorf("reprint payload has no REPRINT header: %q", string(payload))
	}
}

func TestBuildKitchenTicketTestHasTheTestHeaderAndNoOrderNumber(t *testing.T) {
	payload := BuildKitchenTicket(receiptOrder, HeaderTest)

	if !strings.Contains(string(payload), "TEST") {
		t.Errorf("test payload has no TEST header: %q", string(payload))
	}
	if strings.Contains(string(payload), "ORDER ") {
		t.Errorf("test payload should carry no order number: %q", string(payload))
	}
}

func TestFormatCurrency(t *testing.T) {
	cases := map[int]string{0: "$0.00", 5: "$0.05", 1000: "$10.00", 123456: "$1234.56"}
	for cents, want := range cases {
		if got := FormatCurrency(cents); got != want {
			t.Errorf("FormatCurrency(%d) = %q, want %q", cents, got, want)
		}
	}
}

var sideOrder = ReceiptOrder{
	OrderID:       "order-2",
	OrderNumber:   102,
	CreatedAt:     "2026-05-07T12:34:00.000Z",
	SubtotalCents: 1000,
	TotalCents:    1000,
	Items:         []CartLine{{MenuItemID: "potato-pancake", Quantity: 1, Side: "sour-cream"}},
}

func TestBuildCustomerReceiptPrintsTheSideUnderItsLine(t *testing.T) {
	payload := string(BuildCustomerReceipt(sideOrder, HeaderNone))

	if !strings.Contains(payload, "1 x Potato Pancake\r\n  Sour Cream\r\n  $10.00") {
		t.Errorf("the side does not sit between the item and its price: %q", payload)
	}
}

func TestBuildKitchenTicketPrintsTheSideUnderItsLine(t *testing.T) {
	payload := string(BuildKitchenTicket(sideOrder, HeaderNone))

	if !strings.Contains(payload, "1  POTATO PANCAKE\r\n   SOUR CREAM") {
		t.Errorf("the side does not sit under the item line: %q", payload)
	}
}

func TestAPlainLinePrintsNoSide(t *testing.T) {
	receipt := string(BuildCustomerReceipt(receiptOrder, HeaderNone))
	if !strings.Contains(receipt, "2 x Potato Pancake\r\n  $20.00") {
		t.Errorf("a plain line must add no side line to the receipt: %q", receipt)
	}

	ticket := string(BuildKitchenTicket(receiptOrder, HeaderNone))
	if !strings.Contains(ticket, "2  POTATO PANCAKE\r\n\r\n") {
		t.Errorf("a plain line must add no side line to the kitchen ticket: %q", ticket)
	}
}
