package pos

import (
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
		Customer: sendDocument(config, config.WindowHost, config.WindowPort, BuildCustomerReceipt(order, false)),
		Kitchen:  sendDocument(config, config.KitchenHost, config.KitchenPort, BuildKitchenTicket(order, false)),
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
		result.Customer = sendDocument(config, config.WindowHost, config.WindowPort, BuildCustomerReceipt(order, true))
	}
	if sendKitchen {
		result.Kitchen = sendDocument(config, config.KitchenHost, config.KitchenPort, BuildKitchenTicket(order, true))
	}
	return result
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
