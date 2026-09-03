package pos

import (
	"path/filepath"
	"testing"
)

func TestOpenDatabaseRenamesOldPrintedStatusToSent(t *testing.T) {
	path := filepath.Join(t.TempDir(), "pos.sqlite")

	database, err := OpenDatabase(path)
	if err != nil {
		t.Fatalf("open database: %v", err)
	}
	if _, err := database.Exec(
		`INSERT INTO transactions (id, client_order_id, device_id, order_number, business_date, status,
			subtotal_cents, tax_cents, total_cents, payment_method, request_json,
			customer_print_status, kitchen_print_status, created_at)
		 VALUES ('t1', 'c1', 'd1', 100, '2026-05-07', 'printed', 2000, 0, 2000, 'cash', '{}', 'printed', 'printed', '2026-05-07T00:00:00.000Z')`,
	); err != nil {
		t.Fatalf("insert legacy row: %v", err)
	}
	database.Close()

	reopened, err := OpenDatabase(path)
	if err != nil {
		t.Fatalf("reopen database: %v", err)
	}
	defer reopened.Close()

	var customer, kitchen string
	if err := reopened.QueryRow(
		`SELECT customer_print_status, kitchen_print_status FROM transactions WHERE id = 't1'`,
	).Scan(&customer, &kitchen); err != nil {
		t.Fatalf("read row: %v", err)
	}
	if customer != "sent" || kitchen != "sent" {
		t.Errorf("customer = %q, kitchen = %q, want both %q", customer, kitchen, "sent")
	}
}
