export const schemaStatements = [
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
    ON transactions (created_at)`
] as const;
