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

// AdminSalesResponse is one business date's sales report, drawn on the
// Leader page's Figures tab.
type AdminSalesResponse struct {
	BusinessDate string               `json:"businessDate"`
	ServerTime   string               `json:"serverTime"`
	Summary      AdminSalesSummary    `json:"summary"`
	Items        []AdminSalesItemLine `json:"items"`
	Orders       []AdminSalesOrder    `json:"orders"`
	Chart        AdminSalesChart      `json:"-"`
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
		Chart:        buildHourlyChart(orders),
	}, nil
}

// AdminSalesEventTotal is the two-event-day total shown alongside whichever
// single day the Leader is viewing (issue #12).
type AdminSalesEventTotal struct {
	OrderCount int
	TotalCents int
}

// GetEventTotal sums the order count and revenue of the two event dates.
func (service *OrderService) GetEventTotal(saturday, sunday string) (AdminSalesEventTotal, error) {
	var total AdminSalesEventTotal
	for _, date := range []string{saturday, sunday} {
		sales, err := service.GetAdminSales(date)
		if err != nil {
			return AdminSalesEventTotal{}, err
		}
		total.OrderCount += sales.Summary.OrderCount
		total.TotalCents += sales.Summary.TotalCents
	}
	return total, nil
}

// eventDayPair gives the two consecutive business dates of the event
// (CONTEXT.md's Event day: "a Saturday and a Sunday"), derived from today's
// local weekday rather than a configured date, so no deploy config carries
// the actual 2026 calendar dates. During the real event "today" genuinely
// falls on one of the two days; outside it (development, or before the
// event starts) today is neither, and the pair defaults to (today,
// tomorrow) so the toggle still has two distinct dates to show.
func eventDayPair(now func() time.Time) (saturday, sunday string) {
	local := now().In(time.Local)
	if local.Weekday() == time.Sunday {
		return local.AddDate(0, 0, -1).Format("2006-01-02"), local.Format("2006-01-02")
	}
	return local.Format("2006-01-02"), local.AddDate(0, 0, 1).Format("2006-01-02")
}

// Hourly chart geometry (issue #12, Variant A from issue #11): a stacked bar
// per hour, dollars by menu item. html/template only interpolates the
// numbers computed here; no arithmetic happens in the template.
const (
	chartWidth         = 300.0
	chartHeight        = 110.0
	chartBarGap        = 4.0
	chartLabelMargin   = 14.0
	chartLabelBaseline = chartHeight + 12
)

// AdminSalesChart is the hourly revenue-by-item chart drawn on the Figures
// tab. Bars is empty when the day has no non-Voided sales, so the template
// can skip drawing an all-zero chart.
type AdminSalesChart struct {
	ViewWidth  float64
	ViewHeight float64
	LabelY     float64
	Bars       []AdminSalesChartBar
}

// AdminSalesChartBar is one hour's stacked bar.
type AdminSalesChartBar struct {
	X, Width float64
	LabelX   float64
	Label    string
	Segments []AdminSalesChartSegment
}

// AdminSalesChartSegment is one menu item's slice of one hour's bar.
type AdminSalesChartSegment struct {
	Y, Height float64
	ColorVar  string
}

// buildHourlyChart buckets each non-Voided order's items by the local hour
// its CreatedAt (a UTC ISO timestamp) converts to — the same conversion
// pages.go's clock template func uses. The bucket range runs from the
// earliest to the latest hour that actually has a sale, so the chart never
// hard-codes the booth's 8am-6pm hours (CONTEXT.md's Event day).
func buildHourlyChart(orders []AdminSalesOrder) AdminSalesChart {
	revenueByHour := map[int]map[string]int{}
	minHour, maxHour := 0, -1

	for _, order := range orders {
		if order.Status == OrderVoided {
			continue
		}
		moment, err := time.Parse(timestampLayout, order.CreatedAt)
		if err != nil {
			continue
		}
		hour := moment.In(time.Local).Hour()
		if revenueByHour[hour] == nil {
			revenueByHour[hour] = map[string]int{}
		}
		for _, line := range order.Items {
			item, found := MenuItemByID(line.MenuItemID)
			if !found {
				continue
			}
			revenueByHour[hour][line.MenuItemID] += line.Quantity * item.PriceCents
		}
		if maxHour < minHour {
			minHour, maxHour = hour, hour
		} else if hour < minHour {
			minHour = hour
		} else if hour > maxHour {
			maxHour = hour
		}
	}
	if len(revenueByHour) == 0 {
		return AdminSalesChart{}
	}

	hourCount := maxHour - minHour + 1
	maxRevenue := 0
	for hour := minHour; hour <= maxHour; hour++ {
		total := 0
		for _, cents := range revenueByHour[hour] {
			total += cents
		}
		if total > maxRevenue {
			maxRevenue = total
		}
	}
	if maxRevenue == 0 {
		return AdminSalesChart{}
	}

	barWidth := (chartWidth - chartBarGap*float64(hourCount-1)) / float64(hourCount)
	bars := make([]AdminSalesChartBar, 0, hourCount)
	for index := 0; index < hourCount; index++ {
		hour := minHour + index
		x := float64(index) * (barWidth + chartBarGap)
		y := chartHeight
		var segments []AdminSalesChartSegment
		for _, item := range MenuItems {
			cents := revenueByHour[hour][item.ID]
			if cents == 0 {
				continue
			}
			height := float64(cents) / float64(maxRevenue) * chartHeight
			y -= height
			segments = append(segments, AdminSalesChartSegment{Y: y, Height: height, ColorVar: item.ChartColorVar})
		}
		bars = append(bars, AdminSalesChartBar{
			X: x, Width: barWidth, LabelX: x + barWidth/2, Label: hourLabel(hour), Segments: segments,
		})
	}

	return AdminSalesChart{
		ViewWidth:  chartWidth,
		ViewHeight: chartHeight + chartLabelMargin,
		LabelY:     chartLabelBaseline,
		Bars:       bars,
	}
}

// hourLabel writes a 24-hour hour the way a phone screen reads a clock: "8a",
// "2p", "12p".
func hourLabel(hour int) string {
	period := "a"
	if hour >= 12 {
		period = "p"
	}
	display := hour % 12
	if display == 0 {
		display = 12
	}
	return fmt.Sprintf("%d%s", display, period)
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
