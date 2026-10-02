package pos

import (
	"bytes"
	_ "embed"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"strings"
	"sync"
	"time"
)

// The ESC/POS byte sequences are copied from the TypeScript server. They print
// correctly on the real hardware. Do not derive them again.
var (
	initializePrinter = []byte{0x1b, 0x40}
	doubleSizeOn      = []byte{0x1d, 0x21, 0x11}
	doubleSizeOff     = []byte{0x1d, 0x21, 0x00}
	cutPaper          = []byte{0x1d, 0x56, 0x42, 0x08}
)

// The alignment and raster commands were proved on the Window Printer, an
// ITPP047(P), on 2 Oct 2026.
var (
	alignCentre = []byte{0x1b, 0x61, 0x01}
	alignLeft   = []byte{0x1b, 0x61, 0x00}
)

// paperDots is the print width of the 80mm roll.
const paperDots = 576

const (
	receiptTagline  = "Spencerport, NY | facebook.com/troop813"
	receiptThankYou = "Thank you for supporting Troop 813!"
)

// The Troop 813 logo, 448 dots square, black and white only. It was cut from
// the Apple Fest flyer, and its two rings were drawn again as exact circles.
//
//go:embed receipt-logo.png
var receiptLogoPNG []byte

// receiptLogo is the raster of the logo, made once. If the PNG does not
// decode, the receipt prints with no logo.
var receiptLogo = sync.OnceValue(func() []byte {
	logo, err := png.Decode(bytes.NewReader(receiptLogoPNG))
	if err != nil {
		return nil
	}
	return rasterImage(logo)
})

// rasterImage makes the GS v 0 raster command that prints an image in the
// centre of the paper. A pixel darker than mid-grey prints black. The command
// always covers the full paper width, so the centring does not depend on the
// alignment state of the printer.
func rasterImage(source image.Image) []byte {
	bounds := source.Bounds()
	width, height := min(bounds.Dx(), paperDots), bounds.Dy()
	rowBytes := paperDots / 8
	left := (paperDots - width) / 2

	raster := make([]byte, 8, 8+rowBytes*height)
	copy(raster, []byte{0x1d, 0x76, 0x30, 0x00, byte(rowBytes), byte(rowBytes >> 8), byte(height), byte(height >> 8)})
	raster = raster[:cap(raster)]
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			grey := color.GrayModel.Convert(source.At(bounds.Min.X+x, bounds.Min.Y+y)).(color.Gray)
			if grey.Y < 128 {
				dot := left + x
				raster[8+y*rowBytes+dot/8] |= 0x80 >> (dot % 8)
			}
		}
	}
	return raster
}

// ReceiptHeader marks a document as something other than a first-time print
// of a real order, so a reader never mistakes it for one.
type ReceiptHeader string

const (
	HeaderNone    ReceiptHeader = ""
	HeaderReprint ReceiptHeader = "REPRINT"
	HeaderTest    ReceiptHeader = "TEST"
)

// BuildCustomerReceipt makes the ESC/POS bytes of the customer receipt. A
// reprint carries a REPRINT header, so a second copy never looks like the
// original. A test ticket (CONTEXT.md's Test ticket) carries a TEST header
// instead of an order number, so it never looks like a real sale.
//
// The top is centred: the logo, the tagline, then the header and the order
// number at double size, because the order number is how the customer
// collects the food. The item list and the total are on the left.
func BuildCustomerReceipt(order ReceiptOrder, header ReceiptHeader) []byte {
	var number []string
	if header != HeaderNone {
		number = append(number, string(header))
	}
	if header != HeaderTest {
		number = append(number, fmt.Sprintf("Order #%d", order.OrderNumber))
	}

	var lines []string
	for _, line := range order.Items {
		total := "$0.00"
		if item, found := MenuItemByID(line.MenuItemID); found {
			total = FormatCurrency(item.PriceCents * line.Quantity)
		}
		lines = append(lines,
			fmt.Sprintf("%d x %s", line.Quantity, MenuItemName(line.MenuItemID)),
		)
		if labels := SideLabels(line.MenuItemID, line.Sides); len(labels) > 0 {
			lines = append(lines, "  "+strings.Join(labels, ", "))
		}
		lines = append(lines, "  "+total)
	}

	lines = append(lines,
		"",
		fmt.Sprintf("Total %s", FormatCurrency(order.TotalCents)),
	)

	return concat(
		initializePrinter,
		receiptLogo(),
		alignCentre,
		[]byte(receiptTagline+"\r\n\r\n"),
		doubleSizeOn,
		[]byte(strings.Join(number, "\r\n")+"\r\n"),
		doubleSizeOff,
		[]byte(formatTimestamp(order.CreatedAt)+"\r\n\r\n"),
		alignLeft,
		[]byte(strings.Join(lines, "\r\n")+"\r\n\r\n"),
		alignCentre,
		[]byte(receiptThankYou+"\r\n\r\n\r\n\r\n"),
		alignLeft,
		cutPaper,
	)
}

// BuildKitchenTicket makes the ESC/POS bytes of the kitchen ticket. A reprint
// carries a REPRINT header, so the kitchen checks the order number instead of
// cooking the order again. A test ticket carries a TEST header instead of an
// order number, so the kitchen never mistakes it for a real order.
func BuildKitchenTicket(order ReceiptOrder, header ReceiptHeader) []byte {
	var lines []string
	if header != HeaderTest {
		lines = append(lines, fmt.Sprintf("ORDER %d", order.OrderNumber))
	}
	lines = append(lines, formatTimestamp(order.CreatedAt), "")
	if header != HeaderNone {
		lines = append([]string{string(header), ""}, lines...)
	}

	for _, line := range order.Items {
		lines = append(lines, fmt.Sprintf("%d  %s", line.Quantity, strings.ToUpper(MenuItemName(line.MenuItemID))))
		// One topping per line. The kitchen ticket prints at double width, so
		// an 80mm roll holds about 24 characters and nothing here wraps them:
		// "SOUR CREAM, APPLESAUCE" already overflows, and a break mid-word is
		// the one thing the cook must not read.
		for _, label := range SideLabels(line.MenuItemID, line.Sides) {
			lines = append(lines, "   "+strings.ToUpper(label))
		}
	}

	return concat(initializePrinter, doubleSizeOn, encodeLines(lines), doubleSizeOff, cutPaper)
}

// FormatCurrency writes cents as dollars. Integer arithmetic only.
func FormatCurrency(cents int) string {
	sign := ""
	if cents < 0 {
		sign = "-"
		cents = -cents
	}
	return fmt.Sprintf("%s$%d.%02d", sign, cents/100, cents%100)
}

func encodeLines(lines []string) []byte {
	return []byte(strings.Join(lines, "\r\n") + "\r\n\r\n\r\n")
}

func concat(chunks ...[]byte) []byte {
	length := 0
	for _, chunk := range chunks {
		length += len(chunk)
	}

	output := make([]byte, 0, length)
	for _, chunk := range chunks {
		output = append(output, chunk...)
	}
	return output
}

// formatTimestamp writes the receipt time in the local time zone of the Pi.
func formatTimestamp(value string) string {
	parsed, err := time.Parse(timestampLayout, value)
	if err != nil {
		return value
	}
	return parsed.Local().Format("1/2, 3:04 PM")
}
