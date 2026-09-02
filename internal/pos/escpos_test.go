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
	payload := BuildCustomerReceipt(receiptOrder)

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
}

func TestBuildKitchenTicketHasEmphasisAndTheCutCommand(t *testing.T) {
	payload := BuildKitchenTicket(receiptOrder)

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
}

func TestFormatCurrency(t *testing.T) {
	cases := map[int]string{0: "$0.00", 5: "$0.05", 1000: "$10.00", 123456: "$1234.56"}
	for cents, want := range cases {
		if got := FormatCurrency(cents); got != want {
			t.Errorf("FormatCurrency(%d) = %q, want %q", cents, got, want)
		}
	}
}
