package pos

import (
	"io"
	"net"
	"time"
)

// PrinterConfig holds where each of the booth's two thermal printers is, and
// whether printing is on at all. The Customer Receipt goes to the Window
// Printer, and the Kitchen Ticket goes to the Kitchen Printer.
type PrinterConfig struct {
	Enabled     bool
	WindowHost  string
	WindowPort  string
	KitchenHost string
	KitchenPort string
}

const (
	printerDialTimeout  = 2 * time.Second
	printerWriteTimeout = 2 * time.Second
)

// PrintOrder sends the customer receipt to the Window Printer and the kitchen
// ticket to the Kitchen Printer. A printer that is off or unplugged gives
// "failed", not an error: the sale still completes.
func PrintOrder(config PrinterConfig, order ReceiptOrder) PrintResult {
	return PrintResult{
		Customer: sendDocument(config, config.WindowHost, config.WindowPort, BuildCustomerReceipt(order, HeaderNone)),
		Kitchen:  sendDocument(config, config.KitchenHost, config.KitchenPort, BuildKitchenTicket(order, HeaderNone)),
	}
}

// PrintReprint resends the documents of an order that already went through
// PrintOrder once. It sends only the documents that failed, or both when
// nothing failed, and marks every document it sends with REPRINT so the
// kitchen does not cook the order twice.
func PrintReprint(config PrinterConfig, order ReceiptOrder, previous PrintResult) PrintResult {
	sendCustomer, sendKitchen := true, true
	if previous.Customer == PrintFailed || previous.Kitchen == PrintFailed {
		sendCustomer = previous.Customer == PrintFailed
		sendKitchen = previous.Kitchen == PrintFailed
	}

	result := previous
	if sendCustomer {
		result.Customer = sendDocument(config, config.WindowHost, config.WindowPort, BuildCustomerReceipt(order, HeaderReprint))
	}
	if sendKitchen {
		result.Kitchen = sendDocument(config, config.KitchenHost, config.KitchenPort, BuildKitchenTicket(order, HeaderReprint))
	}
	return result
}

// TestOrder is the fixed, synthetic order behind "Send test ticket"
// (ADR-0008): real enough to exercise both builders' side-line and total
// formatting, disposable enough that no operator mistakes it for a sale.
func TestOrder(now time.Time) ReceiptOrder {
	return ReceiptOrder{
		CreatedAt:     now.UTC().Format(timestampLayout),
		SubtotalCents: 1000,
		TotalCents:    1000,
		Items:         []CartLine{{MenuItemID: "potato-pancake", Quantity: 1, Side: "sour-cream"}},
	}
}

// SendTestTicket prints order through both printers as a real, TEST-marked
// document (ADR-0008). It writes no transactions row: this is a manual
// troubleshooting action, never a sale.
func SendTestTicket(config PrinterConfig, order ReceiptOrder) PrintResult {
	return PrintResult{
		Customer: sendDocument(config, config.WindowHost, config.WindowPort, BuildCustomerReceipt(order, HeaderTest)),
		Kitchen:  sendDocument(config, config.KitchenHost, config.KitchenPort, BuildKitchenTicket(order, HeaderTest)),
	}
}

// CheckPrinters dials the Window and Kitchen printers and reads back their
// status (ADR-0008), for the System Admin page's "Check printers" action.
func CheckPrinters(config PrinterConfig) PrinterCheckResult {
	return PrinterCheckResult{
		Window:  checkPrinter(config, config.WindowHost, config.WindowPort),
		Kitchen: checkPrinter(config, config.KitchenHost, config.KitchenPort),
	}
}

// checkPrinter dials one printer and reads its DLE EOT status bytes over the
// same connection. Byte layout: Epson ESC/POS Command Manual, "DLE EOT n" (n=1
// printer status, n=2 off-line status, n=4 paper roll sensor status). The
// paper and cover checks run before the generic off-line check, so a printer
// that is off-line because its cover is open or its paper is out reports the
// named cause, not just "Offline".
func checkPrinter(config PrinterConfig, host string, port string) PrinterStatus {
	if !config.Enabled || host == "" {
		return StatusNotConfigured
	}

	address := net.JoinHostPort(host, port)
	connection, err := net.DialTimeout("tcp", address, printerDialTimeout)
	if err != nil {
		return StatusNotReachable
	}
	defer connection.Close()

	paper, ok := queryStatus(connection, 4)
	if !ok {
		return StatusNotReachable
	}
	if paper&0x60 == 0x60 {
		return StatusPaperOut
	}

	offline, ok := queryStatus(connection, 2)
	if !ok {
		return StatusNotReachable
	}
	if offline&0x04 == 0x04 {
		return StatusCoverOpen
	}

	printer, ok := queryStatus(connection, 1)
	if !ok {
		return StatusNotReachable
	}
	if printer&0x08 == 0x08 {
		return StatusOffline
	}

	return StatusReady
}

// queryStatus sends DLE EOT n and reads back the one status byte the printer
// answers with.
func queryStatus(connection net.Conn, n byte) (byte, bool) {
	if err := connection.SetDeadline(time.Now().Add(printerWriteTimeout)); err != nil {
		return 0, false
	}
	if _, err := connection.Write([]byte{0x10, 0x04, n}); err != nil {
		return 0, false
	}

	var response [1]byte
	if _, err := io.ReadFull(connection, response[:]); err != nil {
		return 0, false
	}
	return response[0], true
}

func sendDocument(config PrinterConfig, host string, port string, payload []byte) PrintStatus {
	if !config.Enabled || host == "" {
		return PrintDisabled
	}
	return printStatus(sendEscPos(host, port, payload))
}

func printStatus(sent bool) PrintStatus {
	if sent {
		return PrintSent
	}
	return PrintFailed
}

// sendEscPos opens TCP to the printer, writes the payload, and closes. The
// timeouts stop an unplugged printer from holding the order POST open for the
// whole operating-system TCP timeout.
func sendEscPos(host string, port string, payload []byte) bool {
	address := net.JoinHostPort(host, port)
	connection, err := net.DialTimeout("tcp", address, printerDialTimeout)
	if err != nil {
		return false
	}
	defer connection.Close()

	if err := connection.SetWriteDeadline(time.Now().Add(printerWriteTimeout)); err != nil {
		return false
	}

	_, err = connection.Write(payload)
	return err == nil
}
