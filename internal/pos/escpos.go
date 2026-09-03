package pos

import (
	"fmt"
	"strings"
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

// BuildCustomerReceipt makes the ESC/POS bytes of the customer receipt. A
// reprint carries a REPRINT header, so a second copy never looks like the
// original.
func BuildCustomerReceipt(order ReceiptOrder, reprint bool) []byte {
	lines := []string{
		"Apple Fest POS",
		fmt.Sprintf("Order #%d", order.OrderNumber),
		formatTimestamp(order.CreatedAt),
		"",
	}
	if reprint {
		lines = append([]string{"REPRINT", ""}, lines...)
	}

	for _, line := range order.Items {
		total := "$0.00"
		if item, found := MenuItemByID(line.MenuItemID); found {
			total = FormatCurrency(item.PriceCents * line.Quantity)
		}
		lines = append(lines,
			fmt.Sprintf("%d x %s", line.Quantity, MenuItemName(line.MenuItemID)),
			"  "+total,
		)
	}

	lines = append(lines,
		"",
		fmt.Sprintf("Total %s", FormatCurrency(order.TotalCents)),
		"Thank you!",
	)

	return concat(initializePrinter, encodeLines(lines), cutPaper)
}

// BuildKitchenTicket makes the ESC/POS bytes of the kitchen ticket. A reprint
// carries a REPRINT header, so the kitchen checks the order number instead of
// cooking the order again.
func BuildKitchenTicket(order ReceiptOrder, reprint bool) []byte {
	lines := []string{
		fmt.Sprintf("ORDER %d", order.OrderNumber),
		formatTimestamp(order.CreatedAt),
		"",
	}
	if reprint {
		lines = append([]string{"REPRINT", ""}, lines...)
	}

	for _, line := range order.Items {
		lines = append(lines, fmt.Sprintf("%d  %s", line.Quantity, strings.ToUpper(MenuItemName(line.MenuItemID))))
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
