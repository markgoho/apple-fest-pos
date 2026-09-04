package pos

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"time"
)

func parseOrderNumber(value string) (int, error) {
	return strconv.Atoi(value)
}

// AdminSalesItemLine is one aggregated menu-item line of the sales report.
type AdminSalesItemLine struct {
	MenuItemID   string `json:"menuItemId"`
	Name         string `json:"name"`
	Quantity     int    `json:"quantity"`
	RevenueCents int    `json:"revenueCents"`
}

// AdminSalesOrderLine is one item of one order in the sales report.
type AdminSalesOrderLine struct {
	MenuItemID string `json:"menuItemId"`
	Name       string `json:"name"`
	Quantity   int    `json:"quantity"`
}

// AdminSalesOrder is one order in the sales report.
type AdminSalesOrder struct {
	ID                  string                `json:"id"`
	OrderNumber         int                   `json:"orderNumber"`
	Status              OrderStatus           `json:"status"`
	SubtotalCents       int                   `json:"subtotalCents"`
	TaxCents            int                   `json:"taxCents"`
	TotalCents          int                   `json:"totalCents"`
	PaymentMethod       string                `json:"paymentMethod"`
	CustomerPrintStatus PrintStatus           `json:"customerPrintStatus"`
	KitchenPrintStatus  PrintStatus           `json:"kitchenPrintStatus"`
	CreatedAt           string                `json:"createdAt"`
	Items               []AdminSalesOrderLine `json:"items"`
}

// AdminSalesSummary is the day total at the top of the sales report.
type AdminSalesSummary struct {
	OrderCount    int `json:"orderCount"`
	TotalCents    int `json:"totalCents"`
	PrintFailures int `json:"printFailures"`
}

// AdminSalesResponse is the body of GET /api/admin/sales.
type AdminSalesResponse struct {
	BusinessDate string               `json:"businessDate"`
	ServerTime   string               `json:"serverTime"`
	Summary      AdminSalesSummary    `json:"summary"`
	Items        []AdminSalesItemLine `json:"items"`
	Orders       []AdminSalesOrder    `json:"orders"`
}

type storedRequest struct {
	Items []CartLine `json:"items"`
	Notes string     `json:"notes"`
}

// CurrentBusinessDate gives today's business date, read in local time (the
// deploy unit sets TZ=America/New_York; time.Local honours it, the same way
// pages.go's clock template func does). Written business dates stay UTC
// (order_store.go's insertOrder): the event runs 8am-6pm Eastern, well
// inside one UTC day, so a write-time UTC date always agrees with the local
// date during those hours. Reading "today" in local time only matters for a
// Leader who opens the page after the UTC day has rolled over but the local
// day has not, around 8pm Eastern.
func CurrentBusinessDate(now func() time.Time) string {
	return now().In(time.Local).Format("2006-01-02")
}

// ReadSalesRows reads every order of one business date, newest order first.
func ReadSalesRows(database *sql.DB, businessDate string) ([]transactionRow, error) {
	rows, err := database.Query(
		`SELECT id, order_number, business_date, status, subtotal_cents, tax_cents,
		        total_cents, payment_method, request_json,
		        customer_print_status, kitchen_print_status, created_at
		 FROM transactions
		 WHERE business_date = ?
		 ORDER BY order_number DESC`, businessDate)
	if err != nil {
		return nil, fmt.Errorf("read sales rows: %w", err)
	}
	defer rows.Close()

	var sales []transactionRow
	for rows.Next() {
		var row transactionRow
		if err := rows.Scan(&row.ID, &row.OrderNumber, &row.BusinessDate, &row.Status,
			&row.SubtotalCents, &row.TaxCents, &row.TotalCents, &row.PaymentMethod,
			&row.RequestJSON, &row.CustomerPrintStatus, &row.KitchenPrintStatus,
			&row.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan sales row: %w", err)
		}
		sales = append(sales, row)
	}
	return sales, rows.Err()
}

func parseStoredRequest(raw string) storedRequest {
	var parsed storedRequest
	if err := json.Unmarshal([]byte(raw), &parsed); err != nil {
		return storedRequest{}
	}
	return parsed
}

// GetAdminSales builds the sales report of one business date. An empty date
// means today.
func (service *OrderService) GetAdminSales(date string) (AdminSalesResponse, error) {
	businessDate := date
	if businessDate == "" {
		businessDate = CurrentBusinessDate(service.Now)
	}

	rows, err := ReadSalesRows(service.DB, businessDate)
	if err != nil {
		return AdminSalesResponse{}, err
	}

	orders := make([]AdminSalesOrder, 0, len(rows))
	summary := AdminSalesSummary{}
	quantities := map[string]int{}

	for _, row := range rows {
		request := parseStoredRequest(row.RequestJSON)
		items := make([]AdminSalesOrderLine, 0, len(request.Items))
		for _, line := range request.Items {
			if row.Status != OrderVoided {
				quantities[line.MenuItemID] += line.Quantity
			}
			items = append(items, AdminSalesOrderLine{
				MenuItemID: line.MenuItemID,
				Name:       MenuItemName(line.MenuItemID),
				Quantity:   line.Quantity,
			})
		}

		// A Voided order drops out of the summary as if it were never sold,
		// but stays in Orders below, marked Voided (CONTEXT.md, issue #45).
		if row.Status != OrderVoided {
			summary.OrderCount++
			summary.TotalCents += row.TotalCents
		}
		if row.Status == OrderPrintFailed {
			summary.PrintFailures++
		}

		orders = append(orders, AdminSalesOrder{
			ID: row.ID, OrderNumber: row.OrderNumber, Status: row.Status,
			SubtotalCents: row.SubtotalCents, TaxCents: row.TaxCents, TotalCents: row.TotalCents,
			PaymentMethod: row.PaymentMethod, CustomerPrintStatus: row.CustomerPrintStatus,
			KitchenPrintStatus: row.KitchenPrintStatus, CreatedAt: row.CreatedAt, Items: items,
		})
	}

	return AdminSalesResponse{
		BusinessDate: businessDate,
		ServerTime:   service.Now().UTC().Format(timestampLayout),
		Summary:      summary,
		Items:        aggregateItems(quantities),
		Orders:       orders,
	}, nil
}

// aggregateItems totals each menu item, most sold first.
func aggregateItems(quantities map[string]int) []AdminSalesItemLine {
	lines := make([]AdminSalesItemLine, 0, len(MenuItems))
	for _, item := range MenuItems {
		quantity := quantities[item.ID]
		if quantity == 0 {
			continue
		}
		lines = append(lines, AdminSalesItemLine{
			MenuItemID:   item.ID,
			Name:         item.Name,
			Quantity:     quantity,
			RevenueCents: quantity * item.PriceCents,
		})
	}

	sort.SliceStable(lines, func(first, second int) bool {
		return lines[first].Quantity > lines[second].Quantity
	})
	return lines
}
