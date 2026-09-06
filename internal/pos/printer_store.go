package pos

import (
	"database/sql"
	"fmt"
	"strconv"
)

// Metadata keys under which the System Admin page's printer assignment is
// stored (ADR-0010). Env vars only seed the very first boot; once any of
// these keys exists, it wins over the env var every later boot.
const (
	metadataWindowHost  = "window_printer_host"
	metadataWindowPort  = "window_printer_port"
	metadataKitchenHost = "kitchen_printer_host"
	metadataKitchenPort = "kitchen_printer_port"
	metadataPrinterOn   = "printer_enabled"
)

// LoadPrinterConfig overlays any printer assignment saved from the System
// Admin page onto the env-derived config main.go built, and makes the result
// the live config. Call it once, right after constructing the OrderService.
func (service *OrderService) LoadPrinterConfig() error {
	base := service.printerConfig()

	transaction, err := service.DB.Begin()
	if err != nil {
		return fmt.Errorf("begin load printer config: %w", err)
	}
	defer transaction.Rollback()

	config, err := readPrinterConfig(transaction, base)
	if err != nil {
		return err
	}
	if err := transaction.Commit(); err != nil {
		return fmt.Errorf("commit load printer config: %w", err)
	}

	service.setPrinterConfig(config)
	return nil
}

// SavePrinterConfig persists config's host and port assignment and makes it
// the live config immediately, with no restart. Enabled is not taken from
// config: it is derived here, true whenever either host is set, so clearing
// both hosts is the only way to turn printing off from the page.
func (service *OrderService) SavePrinterConfig(config PrinterConfig) error {
	config.Enabled = config.WindowHost != "" || config.KitchenHost != ""

	transaction, err := service.DB.Begin()
	if err != nil {
		return fmt.Errorf("begin save printer config: %w", err)
	}
	defer transaction.Rollback()

	writes := map[string]string{
		metadataWindowHost:  config.WindowHost,
		metadataWindowPort:  config.WindowPort,
		metadataKitchenHost: config.KitchenHost,
		metadataKitchenPort: config.KitchenPort,
		metadataPrinterOn:   strconv.FormatBool(config.Enabled),
	}
	for key, value := range writes {
		if err := setMetadataValue(transaction, key, value); err != nil {
			return err
		}
	}
	if err := transaction.Commit(); err != nil {
		return fmt.Errorf("commit save printer config: %w", err)
	}

	service.setPrinterConfig(config)
	return nil
}

func readPrinterConfig(transaction *sql.Tx, base PrinterConfig) (PrinterConfig, error) {
	windowHost, err := metadataValue(transaction, metadataWindowHost, base.WindowHost)
	if err != nil {
		return PrinterConfig{}, err
	}
	windowPort, err := metadataValue(transaction, metadataWindowPort, base.WindowPort)
	if err != nil {
		return PrinterConfig{}, err
	}
	kitchenHost, err := metadataValue(transaction, metadataKitchenHost, base.KitchenHost)
	if err != nil {
		return PrinterConfig{}, err
	}
	kitchenPort, err := metadataValue(transaction, metadataKitchenPort, base.KitchenPort)
	if err != nil {
		return PrinterConfig{}, err
	}
	enabledText, err := metadataValue(transaction, metadataPrinterOn, strconv.FormatBool(base.Enabled))
	if err != nil {
		return PrinterConfig{}, err
	}
	enabled, err := strconv.ParseBool(enabledText)
	if err != nil {
		enabled = base.Enabled
	}

	return PrinterConfig{
		Enabled:     enabled,
		WindowHost:  windowHost,
		WindowPort:  windowPort,
		KitchenHost: kitchenHost,
		KitchenPort: kitchenPort,
	}, nil
}
