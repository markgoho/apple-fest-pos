// Command pos serves the Apple Fest POS from one binary.
package main

import (
	"log"
	"net"
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
			Enabled:     os.Getenv("PRINTER_ENABLED") == "true",
			WindowHost:  os.Getenv("WINDOW_PRINTER_HOST"),
			WindowPort:  envOrDefault("WINDOW_PRINTER_PORT", "9100"),
			KitchenHost: os.Getenv("KITCHEN_PRINTER_HOST"),
			KitchenPort: envOrDefault("KITCHEN_PRINTER_PORT", "9100"),
		},
		StartingOrderNumber: envInt("ORDER_NUMBER_START", 100),
		SystemAdminPIN:      os.Getenv("SYSTEM_ADMIN_PIN"),
		LeaderPIN:           os.Getenv("LEADER_PIN"),
		Now:                 time.Now,
	}

	certificate := os.Getenv("TLS_CERT")
	key := os.Getenv("TLS_KEY")
	if certificate == "" || key == "" {
		address := ":" + envOrDefault("PORT", "3000")
		log.Printf("Apple Fest POS listening on %s without TLS", address)
		if err := http.ListenAndServe(address, service.Handler()); err != nil {
			log.Fatalf("serve: %v", err)
		}
		return
	}

	go func() {
		address := ":" + envOrDefault("HTTP_PORT", "80")
		log.Printf("Apple Fest POS redirects %s to HTTPS", address)
		if err := http.ListenAndServe(address, http.HandlerFunc(redirectToHTTPS)); err != nil {
			log.Printf("redirect listener stopped: %v", err)
		}
	}()

	address := ":" + envOrDefault("PORT", "443")
	log.Printf("Apple Fest POS listening on %s with TLS", address)
	if err := http.ListenAndServeTLS(address, certificate, key, service.Handler()); err != nil {
		log.Fatalf("serve: %v", err)
	}
}

// redirectToHTTPS sends a plain HTTP request to the same path on HTTPS, so a
// person who types the bare booth name still gets the POS.
func redirectToHTTPS(writer http.ResponseWriter, request *http.Request) {
	host := request.Host
	if name, _, err := net.SplitHostPort(host); err == nil {
		host = name
	}
	http.Redirect(writer, request, "https://"+host+request.URL.RequestURI(), http.StatusMovedPermanently)
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
