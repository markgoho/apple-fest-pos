package pos

import (
	"testing"
)

func TestLoadPrinterConfigKeepsTheEnvSeedWhenNothingIsSaved(t *testing.T) {
	service := newTestService(t)
	service.Printer = PrinterConfig{Enabled: false, WindowHost: "192.168.8.200", WindowPort: "9100"}

	if err := service.LoadPrinterConfig(); err != nil {
		t.Fatalf("load printer config: %v", err)
	}

	if got := service.printerConfig(); got.WindowHost != "192.168.8.200" {
		t.Errorf("WindowHost = %q, want the env seed unchanged", got.WindowHost)
	}
}

func TestSaveThenLoadPrinterConfigRoundTrips(t *testing.T) {
	service := newTestService(t)

	err := service.SavePrinterConfig(PrinterConfig{
		WindowHost:  "192.168.8.200",
		WindowPort:  "9100",
		KitchenHost: "192.168.8.201",
		KitchenPort: "9100",
	})
	if err != nil {
		t.Fatalf("save printer config: %v", err)
	}

	live := service.printerConfig()
	if !live.Enabled {
		t.Errorf("Enabled = false, want true once both hosts are set")
	}

	reloaded := &OrderService{DB: service.DB}
	if err := reloaded.LoadPrinterConfig(); err != nil {
		t.Fatalf("load printer config: %v", err)
	}
	got := reloaded.printerConfig()
	if got.WindowHost != "192.168.8.200" || got.KitchenHost != "192.168.8.201" || !got.Enabled {
		t.Errorf("reloaded config = %+v, want the saved assignment", got)
	}
}

func TestSavePrinterConfigDisablesWhenBothHostsAreCleared(t *testing.T) {
	service := newTestService(t)
	if err := service.SavePrinterConfig(PrinterConfig{WindowHost: "192.168.8.200", WindowPort: "9100"}); err != nil {
		t.Fatalf("save printer config: %v", err)
	}

	if err := service.SavePrinterConfig(PrinterConfig{}); err != nil {
		t.Fatalf("save printer config: %v", err)
	}

	if got := service.printerConfig(); got.Enabled {
		t.Errorf("Enabled = true, want false once both hosts are cleared")
	}
}
