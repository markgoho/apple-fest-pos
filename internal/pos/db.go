package pos

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

var schemaStatements = []string{
	`CREATE TABLE IF NOT EXISTS metadata (
    key TEXT PRIMARY KEY,
    value TEXT NOT NULL
  )`,
	`CREATE TABLE IF NOT EXISTS transactions (
    id TEXT PRIMARY KEY,
    client_order_id TEXT NOT NULL UNIQUE,
    device_id TEXT NOT NULL,
    order_number INTEGER NOT NULL,
    business_date TEXT NOT NULL,
    status TEXT NOT NULL,
    subtotal_cents INTEGER NOT NULL,
    tax_cents INTEGER NOT NULL,
    total_cents INTEGER NOT NULL,
    payment_method TEXT NOT NULL,
    request_json TEXT NOT NULL,
    customer_print_status TEXT NOT NULL,
    kitchen_print_status TEXT NOT NULL,
    created_at TEXT NOT NULL
  )`,
	`CREATE INDEX IF NOT EXISTS idx_transactions_business_date_order_number
    ON transactions (business_date, order_number)`,
	`CREATE INDEX IF NOT EXISTS idx_transactions_created_at
    ON transactions (created_at)`,
}

// OpenDatabase opens the SQLite file at path, makes its directory, and applies
// the schema. The pool holds one connection, so the PRAGMAs stay applied and
// two tablets that post at the same time are serialized.
func OpenDatabase(path string) (*sql.DB, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, fmt.Errorf("make database directory: %w", err)
	}

	database, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}
	database.SetMaxOpenConns(1)

	pragmas := []string{
		"PRAGMA journal_mode = WAL",
		"PRAGMA synchronous = NORMAL",
		"PRAGMA busy_timeout = 5000",
	}
	for _, statement := range append(pragmas, schemaStatements...) {
		if _, err := database.Exec(statement); err != nil {
			database.Close()
			return nil, fmt.Errorf("apply %q: %w", statement, err)
		}
	}

	return database, nil
}
