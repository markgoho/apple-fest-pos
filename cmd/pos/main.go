// Command pos serves the Apple Fest POS from one binary.
package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/markgoho/apple-fest-pos/internal/pos"
)

func main() {
	database, err := pos.OpenDatabase(sqlitePath())
	if err != nil {
		log.Fatalf("open database: %v", err)
	}
	defer database.Close()

	service := &pos.OrderService{
		DB: database,
		Printer: pos.PrinterConfig{
			Enabled: os.Getenv("PRINTER_ENABLED") == "true",
			Host:    os.Getenv("PRINTER_HOST"),
			Port:    envOrDefault("PRINTER_PORT", "9100"),
		},
		StartingOrderNumber: envInt("ORDER_NUMBER_START", 100),
		Now:                 time.Now,
	}

	address := ":" + envOrDefault("PORT", "3000")
	log.Printf("Apple Fest POS listening on %s", address)
	if err := http.ListenAndServe(address, service.Handler()); err != nil {
		log.Fatalf("serve: %v", err)
	}
}

func sqlitePath() string {
	if path := os.Getenv("SQLITE_PATH"); path != "" {
		return path
	}
	working, err := os.Getwd()
	if err != nil {
		log.Fatalf("find working directory: %v", err)
	}
	return filepath.Join(working, "data", "pos.sqlite")
}

func envOrDefault(name string, fallback string) string {
	if value := os.Getenv(name); value != "" {
		return value
	}
	return fallback
}

func envInt(name string, fallback int) int {
	value, err := strconv.Atoi(os.Getenv(name))
	if err != nil {
		return fallback
	}
	return value
}
