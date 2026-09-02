package pos

import (
	"net"
	"time"
)

// PrinterConfig holds where the thermal printer is, and whether to use it.
type PrinterConfig struct {
	Enabled bool
	Host    string
	Port    string
}

const (
	printerDialTimeout  = 4 * time.Second
	printerWriteTimeout = 4 * time.Second
)

// PrintOrder sends the customer receipt and the kitchen ticket to the printer.
// A printer that is off or unplugged gives "failed", not an error: the sale
// still completes.
func PrintOrder(config PrinterConfig, order ReceiptOrder) PrintResult {
	if !config.Enabled || config.Host == "" {
		return PrintResult{Customer: PrintDisabled, Kitchen: PrintDisabled}
	}

	return PrintResult{
		Customer: printStatus(sendEscPos(config, BuildCustomerReceipt(order))),
		Kitchen:  printStatus(sendEscPos(config, BuildKitchenTicket(order))),
	}
}

func printStatus(sent bool) PrintStatus {
	if sent {
		return PrintPrinted
	}
	return PrintFailed
}

// sendEscPos opens TCP to the printer, writes the payload, and closes. The
// timeouts stop an unplugged printer from holding the order POST open for the
// whole operating-system TCP timeout.
func sendEscPos(config PrinterConfig, payload []byte) bool {
	address := net.JoinHostPort(config.Host, config.Port)
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
