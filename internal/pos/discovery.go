package pos

import (
	"net"
	"sort"
	"sync"
	"time"
)

const (
	discoveryDialTimeout = 300 * time.Millisecond
	discoveryConcurrency = 128
	// A subnet bigger than /22 (1022 hosts) is not a booth network, and
	// scanning it host by host would make the button hang. Skip it.
	discoveryMinPrefixBits = 22
)

// DiscoverPrinters scans the server's own local IPv4 subnets for hosts that
// answer as ESC/POS printers on port 9100, the raw port every ITPP047(P)
// listens on out of the box. ADR-0010: this is what makes the System Admin
// page's printer assignment usable without knowing an address in advance.
func DiscoverPrinters() []string {
	return scanForPrinters(localIPv4Subnets(), "9100")
}

// localIPv4Subnets lists the private IPv4 networks this machine has an
// address on, skipping loopback and anything too large to scan by hand.
func localIPv4Subnets() []*net.IPNet {
	addresses, err := net.InterfaceAddrs()
	if err != nil {
		return nil
	}

	var subnets []*net.IPNet
	for _, address := range addresses {
		ipNet, ok := address.(*net.IPNet)
		if !ok || ipNet.IP.IsLoopback() || ipNet.IP.To4() == nil {
			continue
		}
		ones, _ := ipNet.Mask.Size()
		if ones < discoveryMinPrefixBits {
			continue
		}
		subnets = append(subnets, &net.IPNet{IP: ipNet.IP.Mask(ipNet.Mask).To4(), Mask: ipNet.Mask})
	}
	return subnets
}

// hostsIn lists every usable host address in subnet: every address between
// the network and broadcast address, exclusive.
func hostsIn(subnet *net.IPNet) []string {
	ones, bits := subnet.Mask.Size()
	hostCount := 1 << (bits - ones)
	if hostCount < 3 {
		return nil
	}

	base := subnet.IP.Mask(subnet.Mask).To4()
	start := uint32(base[0])<<24 | uint32(base[1])<<16 | uint32(base[2])<<8 | uint32(base[3])

	hosts := make([]string, 0, hostCount-2)
	for offset := 1; offset < hostCount-1; offset++ {
		value := start + uint32(offset)
		ip := net.IPv4(byte(value>>24), byte(value>>16), byte(value>>8), byte(value))
		hosts = append(hosts, ip.String())
	}
	return hosts
}

// scanForPrinters dials every host in every subnet concurrently and keeps
// the ones that answer an ESC/POS status query on port.
func scanForPrinters(subnets []*net.IPNet, port string) []string {
	var hosts []string
	for _, subnet := range subnets {
		hosts = append(hosts, hostsIn(subnet)...)
	}

	var (
		mutex sync.Mutex
		found []string
		wait  sync.WaitGroup
	)
	semaphore := make(chan struct{}, discoveryConcurrency)

	for _, host := range hosts {
		wait.Add(1)
		semaphore <- struct{}{}
		go func(host string) {
			defer wait.Done()
			defer func() { <-semaphore }()
			if looksLikePrinter(host, port) {
				mutex.Lock()
				found = append(found, host)
				mutex.Unlock()
			}
		}(host)
	}
	wait.Wait()

	sort.Strings(found)
	return found
}

// looksLikePrinter dials host:port and asks its ESC/POS printer status (DLE
// EOT 1), the same query checkPrinter uses. A plain open port answers
// nothing recognizable and times out, so this tells a printer apart from
// whatever else happens to listen on 9100.
func looksLikePrinter(host string, port string) bool {
	address := net.JoinHostPort(host, port)
	connection, err := net.DialTimeout("tcp", address, discoveryDialTimeout)
	if err != nil {
		return false
	}
	defer connection.Close()

	_, ok := queryStatus(connection, 1)
	return ok
}
