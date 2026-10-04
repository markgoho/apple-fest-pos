package pos

import (
	"bytes"
	"image"
	"image/color"
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

// The order number is how the customer collects the food, so it prints at
// double size, and the size goes back to normal before the item list.
func TestBuildCustomerReceiptPrintsTheOrderNumberAtDoubleSize(t *testing.T) {
	payload := BuildCustomerReceipt(receiptOrder, HeaderNone)

	if !bytes.Contains(payload, []byte("\x1d\x21\x11Order #101\r\n\x1d\x21\x00")) {
		t.Errorf("the order number is not between the double size commands")
	}
	if strings.Contains(string(payload), "Apple Fest POS") {
		t.Errorf("the receipt should not carry the name of the software")
	}
	if !strings.Contains(string(payload), receiptTagline) || !strings.Contains(string(payload), receiptThankYou) {
		t.Errorf("the receipt has no tagline or no thank-you line")
	}
}

func TestBuildCustomerReceiptStartsWithTheLogo(t *testing.T) {
	payload := BuildCustomerReceipt(receiptOrder, HeaderNone)

	// GS v 0, mode 0, 72 bytes for each row, 448 rows.
	header := []byte{0x1b, 0x40, 0x1d, 0x76, 0x30, 0x00, 0x48, 0x00, 0xc0, 0x01}
	if !bytes.HasPrefix(payload, header) {
		t.Fatalf("payload does not start with the logo raster: % x", payload[:10])
	}
	if got, want := len(receiptLogo()), 8+72*448; got != want {
		t.Errorf("logo raster is %d bytes, want %d", got, want)
	}
	afterLogo := payload[2+len(receiptLogo()):]
	if !bytes.HasPrefix(afterLogo, append([]byte{0x1b, 0x61, 0x01}, receiptTagline...)) {
		t.Errorf("the centred tagline does not follow the logo: %q", afterLogo[:20])
	}
}

func TestRasterImagePutsTheImageInTheCentreOfThePaper(t *testing.T) {
	source := image.NewGray(image.Rect(0, 0, 8, 2))
	for i := range source.Pix {
		source.Pix[i] = 255
	}
	source.SetGray(0, 0, color.Gray{Y: 0})
	source.SetGray(7, 1, color.Gray{Y: 100})

	raster := rasterImage(source)

	if want := []byte{0x1d, 0x76, 0x30, 0x00, 72, 0, 2, 0}; !bytes.Equal(raster[:8], want) {
		t.Fatalf("raster header = % x, want % x", raster[:8], want)
	}
	// Dot 284 is the left edge of an 8-dot image on 576 dots: byte 35, bit 4.
	want := make([]byte, 2*72)
	want[35] = 0x08
	want[72+36] = 0x10
	if !bytes.Equal(raster[8:], want) {
		t.Errorf("raster data has the wrong dots set: % x", raster[8:])
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
	if !strings.Contains(string(payload), "POTATO PANCAKE") {
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
	Items:         []CartLine{{MenuItemID: "potato-pancake", Quantity: 1, Sides: []string{"ketchup", "sour-cream"}}},
}

// The sides print in menu order (Sour Cream, Applesauce, Ketchup), not in the
// order the Operator tapped, so the paper always reads the same way.
func TestBuildCustomerReceiptPrintsEverySideUnderItsLine(t *testing.T) {
	payload := string(BuildCustomerReceipt(sideOrder, HeaderNone))

	if !strings.Contains(payload, "1 x Potato Pancake\r\n  Sour Cream, Ketchup\r\n  $10.00") {
		t.Errorf("the sides do not sit between the item and its price: %q", payload)
	}
}

func TestBuildKitchenTicketPrintsEverySideUnderItsLine(t *testing.T) {
	payload := string(BuildKitchenTicket(sideOrder, HeaderNone))

	// One per line: the ticket prints at double width and nothing wraps it.
	if !strings.Contains(payload, "POTATO PANCAKE\r\n   SOUR CREAM\r\n   KETCHUP") {
		t.Errorf("the sides do not sit under the item line, one each: %q", payload)
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

// A young cook reads "2  POTATO PANCAKE / SOUR CREAM" as sour cream on one
// pancake only. A line with Sides prints once per unit, with its Sides under
// each one. A line with no Sides keeps its count, shown only above one.
func TestBuildKitchenTicketRepeatsOnlyLinesWithSides(t *testing.T) {
	order := ReceiptOrder{
		OrderNumber: 103,
		CreatedAt:   "2026-10-03T16:00:00.000Z",
		Items: []CartLine{
			{MenuItemID: "potato-pancake", Quantity: 2, Sides: []string{"sour-cream"}},
			{MenuItemID: "harvest-toastie", Quantity: 2},
			{MenuItemID: "og-toastie", Quantity: 1},
		},
	}
	payload := string(BuildKitchenTicket(order, HeaderNone))

	want := "POTATO PANCAKE\r\n   SOUR CREAM\r\nPOTATO PANCAKE\r\n   SOUR CREAM\r\n2  HARVEST TOASTIE\r\nOG TOASTIE\r\n"
	if !strings.Contains(payload, want) {
		t.Errorf("ticket = %q, want it to contain %q", payload, want)
	}
}
