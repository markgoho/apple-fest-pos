package pos

import (
	"net"
	"strconv"
	"testing"
)

func TestLooksLikePrinterTrueWhenItAnswersStatus(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	defer listener.Close()
	serveStatus(t, listener, map[byte]byte{1: 0x12})

	port := strconv.Itoa(listener.Addr().(*net.TCPAddr).Port)
	if !looksLikePrinter("127.0.0.1", port) {
		t.Errorf("looksLikePrinter = false, want true for a listener answering DLE EOT")
	}
}

func TestLooksLikePrinterFalseWhenNothingListens(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	address := listener.Addr().(*net.TCPAddr)
	listener.Close()

	if looksLikePrinter("127.0.0.1", strconv.Itoa(address.Port)) {
		t.Errorf("looksLikePrinter = true, want false when nothing is listening")
	}
}

func TestLooksLikePrinterFalseWhenTheConnectionGivesNoStatus(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	defer listener.Close()
	go func() {
		connection, err := listener.Accept()
		if err != nil {
			return
		}
		defer connection.Close()
		// Accepts the connection but never answers: not an ESC/POS printer.
		buffer := make([]byte, 3)
		connection.Read(buffer)
	}()

	port := strconv.Itoa(listener.Addr().(*net.TCPAddr).Port)
	if looksLikePrinter("127.0.0.1", port) {
		t.Errorf("looksLikePrinter = true, want false when the peer never answers")
	}
}

func TestHostsInListsEveryAddressBetweenNetworkAndBroadcast(t *testing.T) {
	_, subnet, err := net.ParseCIDR("192.168.8.8/29")
	if err != nil {
		t.Fatalf("parse CIDR: %v", err)
	}

	hosts := hostsIn(subnet)

	want := []string{"192.168.8.9", "192.168.8.10", "192.168.8.11", "192.168.8.12", "192.168.8.13", "192.168.8.14"}
	if len(hosts) != len(want) {
		t.Fatalf("hosts = %v, want %v", hosts, want)
	}
	for i, host := range want {
		if hosts[i] != host {
			t.Errorf("hosts[%d] = %q, want %q", i, hosts[i], host)
		}
	}
}

func TestScanForPrintersFindsOnlyTheHostThatAnswers(t *testing.T) {
	// 127.0.0.0/30 has two usable hosts: .1 and .2. Only .1 has a listener,
	// so the scan of both must come back with just .1.
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	defer listener.Close()
	serveStatus(t, listener, map[byte]byte{1: 0x12})

	port := strconv.Itoa(listener.Addr().(*net.TCPAddr).Port)
	_, subnet, err := net.ParseCIDR("127.0.0.0/30")
	if err != nil {
		t.Fatalf("parse CIDR: %v", err)
	}

	found := scanForPrinters([]*net.IPNet{subnet}, port)

	if len(found) != 1 || found[0] != "127.0.0.1" {
		t.Errorf("found = %v, want [127.0.0.1]", found)
	}
}
