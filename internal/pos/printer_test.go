package pos

import (
	"io"
	"net"
	"strconv"
	"testing"
)

func TestPrintOrderReportsDisabledWhenThePrinterIsOff(t *testing.T) {
	result := PrintOrder(PrinterConfig{Enabled: false, Host: "127.0.0.1", Port: "9100"}, receiptOrder)

	if result.Customer != PrintDisabled || result.Kitchen != PrintDisabled {
		t.Errorf("result = %+v, want both disabled", result)
	}
}

func TestPrintOrderReportsFailedWhenNoPrinterAnswers(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	address := listener.Addr().(*net.TCPAddr)
	listener.Close()

	result := PrintOrder(PrinterConfig{
		Enabled: true,
		Host:    "127.0.0.1",
		Port:    strconv.Itoa(address.Port),
	}, receiptOrder)

	if result.Customer != PrintFailed || result.Kitchen != PrintFailed {
		t.Errorf("result = %+v, want both failed", result)
	}
}

func TestPrintOrderSendsBothPayloads(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	defer listener.Close()

	received := make(chan []byte, 2)
	go func() {
		for range 2 {
			connection, err := listener.Accept()
			if err != nil {
				return
			}
			payload, _ := io.ReadAll(connection)
			connection.Close()
			received <- payload
		}
	}()

	result := PrintOrder(PrinterConfig{
		Enabled: true,
		Host:    "127.0.0.1",
		Port:    strconv.Itoa(listener.Addr().(*net.TCPAddr).Port),
	}, receiptOrder)

	if result.Customer != PrintPrinted || result.Kitchen != PrintPrinted {
		t.Fatalf("result = %+v, want both printed", result)
	}

	for range 2 {
		if payload := <-received; len(payload) == 0 {
			t.Errorf("the printer got an empty payload")
		}
	}
}
