package pos

import (
	"io"
	"net"
	"strconv"
	"strings"
	"testing"
)

func TestPrintOrderReportsDisabledWhenThePrinterIsOff(t *testing.T) {
	result := PrintOrder(PrinterConfig{Enabled: false, WindowHost: "127.0.0.1", WindowPort: "9100", KitchenHost: "127.0.0.1", KitchenPort: "9100"}, receiptOrder)

	if result.Customer != PrintDisabled || result.Kitchen != PrintDisabled {
		t.Errorf("result = %+v, want both disabled", result)
	}
}

func TestPrintOrderReportsDisabledWhenAHostIsEmpty(t *testing.T) {
	result := PrintOrder(PrinterConfig{Enabled: true, WindowHost: "", KitchenHost: "127.0.0.1", KitchenPort: "9100"}, receiptOrder)

	if result.Customer != PrintDisabled {
		t.Errorf("customer = %v, want disabled", result.Customer)
	}
}

func TestPrintOrderReportsFailedWhenNoPrinterAnswers(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	address := listener.Addr().(*net.TCPAddr)
	listener.Close()

	port := strconv.Itoa(address.Port)
	result := PrintOrder(PrinterConfig{
		Enabled:     true,
		WindowHost:  "127.0.0.1",
		WindowPort:  port,
		KitchenHost: "127.0.0.1",
		KitchenPort: port,
	}, receiptOrder)

	if result.Customer != PrintFailed || result.Kitchen != PrintFailed {
		t.Errorf("result = %+v, want both failed", result)
	}
}

func acceptTwo(t *testing.T, listener net.Listener) chan []byte {
	t.Helper()
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
	return received
}

func TestPrintOrderSendsBothPayloads(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	defer listener.Close()

	received := acceptTwo(t, listener)
	port := strconv.Itoa(listener.Addr().(*net.TCPAddr).Port)

	result := PrintOrder(PrinterConfig{
		Enabled:     true,
		WindowHost:  "127.0.0.1",
		WindowPort:  port,
		KitchenHost: "127.0.0.1",
		KitchenPort: port,
	}, receiptOrder)

	if result.Customer != PrintSent || result.Kitchen != PrintSent {
		t.Fatalf("result = %+v, want both sent", result)
	}

	for range 2 {
		if payload := <-received; len(payload) == 0 {
			t.Errorf("the printer got an empty payload")
		}
	}
}

func TestPrintReprintSendsOnlyWhatFailed(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	defer listener.Close()

	received := make(chan []byte, 1)
	go func() {
		connection, err := listener.Accept()
		if err != nil {
			return
		}
		payload, _ := io.ReadAll(connection)
		connection.Close()
		received <- payload
	}()

	config := PrinterConfig{
		Enabled:     true,
		WindowHost:  "127.0.0.1",
		WindowPort:  strconv.Itoa(listener.Addr().(*net.TCPAddr).Port),
		KitchenHost: "",
	}

	result := PrintReprint(config, receiptOrder, PrintResult{Customer: PrintFailed, Kitchen: PrintSent})

	if result.Customer != PrintSent {
		t.Errorf("customer = %v, want sent", result.Customer)
	}
	if result.Kitchen != PrintSent {
		t.Errorf("kitchen = %v, want left as sent, unresent", result.Kitchen)
	}

	payload := <-received
	if !strings.Contains(string(payload), "REPRINT") {
		t.Errorf("reprinted payload has no REPRINT header: %q", payload)
	}
}

func TestSendTestTicketMarksBothDocumentsTest(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	defer listener.Close()

	received := acceptTwo(t, listener)
	port := strconv.Itoa(listener.Addr().(*net.TCPAddr).Port)

	result := SendTestTicket(PrinterConfig{
		Enabled:     true,
		WindowHost:  "127.0.0.1",
		WindowPort:  port,
		KitchenHost: "127.0.0.1",
		KitchenPort: port,
	}, receiptOrder)

	if result.Customer != PrintSent || result.Kitchen != PrintSent {
		t.Fatalf("result = %+v, want both sent", result)
	}
	for range 2 {
		if payload := <-received; !strings.Contains(string(payload), "TEST") {
			t.Errorf("test ticket payload has no TEST header: %q", payload)
		}
	}
}

func TestCheckPrintersReportsNotConfiguredWhenAHostIsEmpty(t *testing.T) {
	result := CheckPrinters(PrinterConfig{Enabled: true, WindowHost: "", KitchenHost: "127.0.0.1", KitchenPort: "9100"})

	if result.Window != StatusNotConfigured {
		t.Errorf("window = %v, want not configured", result.Window)
	}
}

func TestCheckPrintersReportsNotReachableWhenNoPrinterAnswers(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	address := listener.Addr().(*net.TCPAddr)
	listener.Close()

	port := strconv.Itoa(address.Port)
	result := CheckPrinters(PrinterConfig{Enabled: true, WindowHost: "127.0.0.1", WindowPort: port, KitchenHost: "127.0.0.1", KitchenPort: port})

	if result.Window != StatusNotReachable || result.Kitchen != StatusNotReachable {
		t.Errorf("result = %+v, want both not reachable", result)
	}
}

// serveStatus answers every DLE EOT n query on one connection with the given
// byte for n, so a test can script a printer's status without real hardware.
func serveStatus(t *testing.T, listener net.Listener, statusFor map[byte]byte) {
	t.Helper()
	go func() {
		connection, err := listener.Accept()
		if err != nil {
			return
		}
		defer connection.Close()
		for {
			var query [3]byte
			if _, err := io.ReadFull(connection, query[:]); err != nil {
				return
			}
			if _, err := connection.Write([]byte{statusFor[query[2]]}); err != nil {
				return
			}
		}
	}()
}

func TestCheckPrintersReportsReady(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	defer listener.Close()
	serveStatus(t, listener, map[byte]byte{1: 0x12, 2: 0x12, 4: 0x12})

	port := strconv.Itoa(listener.Addr().(*net.TCPAddr).Port)
	result := CheckPrinters(PrinterConfig{Enabled: true, WindowHost: "127.0.0.1", WindowPort: port, KitchenHost: "", KitchenPort: ""})

	if result.Window != StatusReady {
		t.Errorf("window = %v, want ready", result.Window)
	}
}

func TestCheckPrintersReportsPaperOut(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	defer listener.Close()
	serveStatus(t, listener, map[byte]byte{1: 0x12, 2: 0x12, 4: 0x72})

	port := strconv.Itoa(listener.Addr().(*net.TCPAddr).Port)
	result := CheckPrinters(PrinterConfig{Enabled: true, WindowHost: "127.0.0.1", WindowPort: port, KitchenHost: "", KitchenPort: ""})

	if result.Window != StatusPaperOut {
		t.Errorf("window = %v, want paper out", result.Window)
	}
}

func TestCheckPrintersReportsCoverOpen(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	defer listener.Close()
	serveStatus(t, listener, map[byte]byte{1: 0x12, 2: 0x16, 4: 0x12})

	port := strconv.Itoa(listener.Addr().(*net.TCPAddr).Port)
	result := CheckPrinters(PrinterConfig{Enabled: true, WindowHost: "127.0.0.1", WindowPort: port, KitchenHost: "", KitchenPort: ""})

	if result.Window != StatusCoverOpen {
		t.Errorf("window = %v, want cover open", result.Window)
	}
}

func TestCheckPrintersReportsOffline(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	defer listener.Close()
	serveStatus(t, listener, map[byte]byte{1: 0x1a, 2: 0x12, 4: 0x12})

	port := strconv.Itoa(listener.Addr().(*net.TCPAddr).Port)
	result := CheckPrinters(PrinterConfig{Enabled: true, WindowHost: "127.0.0.1", WindowPort: port, KitchenHost: "", KitchenPort: ""})

	if result.Window != StatusOffline {
		t.Errorf("window = %v, want offline", result.Window)
	}
}

func TestPrintReprintSendsBothWhenNothingFailed(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	defer listener.Close()

	received := acceptTwo(t, listener)
	port := strconv.Itoa(listener.Addr().(*net.TCPAddr).Port)

	config := PrinterConfig{
		Enabled:     true,
		WindowHost:  "127.0.0.1",
		WindowPort:  port,
		KitchenHost: "127.0.0.1",
		KitchenPort: port,
	}

	result := PrintReprint(config, receiptOrder, PrintResult{Customer: PrintSent, Kitchen: PrintSent})

	if result.Customer != PrintSent || result.Kitchen != PrintSent {
		t.Fatalf("result = %+v, want both sent", result)
	}
	for range 2 {
		payload := <-received
		if !strings.Contains(string(payload), "REPRINT") {
			t.Errorf("reprinted payload has no REPRINT header: %q", payload)
		}
	}
}
