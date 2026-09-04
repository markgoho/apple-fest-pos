package pos

import "testing"

func TestAdminSalesAggregatesTheDay(t *testing.T) {
	service := newTestService(t)

	first := validOrder()
	first["clientOrderId"] = "sales-1"
	postOrder(t, service, first)

	second := validOrder()
	second["clientOrderId"] = "sales-2"
	second["items"] = []map[string]any{
		{"menuItemId": "og-toastie", "quantity": 3},
		{"menuItemId": "potato-pancake", "quantity": 1},
	}
	postOrder(t, service, second)

	sales, err := service.GetAdminSales("")
	if err != nil {
		t.Fatalf("admin sales: %v", err)
	}

	if sales.Summary.OrderCount != 2 {
		t.Errorf("orderCount = %d, want 2", sales.Summary.OrderCount)
	}
	if want := 2000 + 3*500 + 1000; sales.Summary.TotalCents != want {
		t.Errorf("totalCents = %d, want %d", sales.Summary.TotalCents, want)
	}
	if sales.Summary.PrintFailures != 0 {
		t.Errorf("printFailures = %d, want 0", sales.Summary.PrintFailures)
	}

	// Newest order first.
	if sales.Orders[0].OrderNumber != 101 {
		t.Errorf("first order = %d, want 101", sales.Orders[0].OrderNumber)
	}

	// Most sold item first: three OG Toasties beat three potato pancakes only
	// on a tie, so check both lines.
	if len(sales.Items) != 2 {
		t.Fatalf("items = %+v, want two lines", sales.Items)
	}
	if sales.Items[0].MenuItemID != "potato-pancake" || sales.Items[0].Quantity != 3 {
		t.Errorf("first item line = %+v", sales.Items[0])
	}
	if sales.Items[0].RevenueCents != 3000 {
		t.Errorf("revenueCents = %d, want 3000", sales.Items[0].RevenueCents)
	}
}

func TestAdminSalesDropsAVoidedOrderFromTheSummary(t *testing.T) {
	service := newTestService(t)

	first := validOrder()
	first["clientOrderId"] = "voided-sales-1"
	_, body := postOrder(t, service, first)
	orderID := body["order"].(map[string]any)["id"].(string)

	second := validOrder()
	second["clientOrderId"] = "voided-sales-2"
	postOrder(t, service, second)

	if _, err := service.VoidOrder(orderID); err != nil {
		t.Fatalf("void order: %v", err)
	}

	sales, err := service.GetAdminSales("")
	if err != nil {
		t.Fatalf("admin sales: %v", err)
	}

	if sales.Summary.OrderCount != 1 {
		t.Errorf("orderCount = %d, want 1", sales.Summary.OrderCount)
	}
	if sales.Summary.TotalCents != 2000 {
		t.Errorf("totalCents = %d, want 2000", sales.Summary.TotalCents)
	}
	if len(sales.Items) != 1 || sales.Items[0].Quantity != 2 {
		t.Errorf("items = %+v, want one line of quantity 2", sales.Items)
	}

	if len(sales.Orders) != 2 {
		t.Fatalf("orders = %d, want 2", len(sales.Orders))
	}
	var voided AdminSalesOrder
	for _, order := range sales.Orders {
		if order.ID == orderID {
			voided = order
		}
	}
	if voided.Status != OrderVoided {
		t.Errorf("voided order status = %q, want %q", voided.Status, OrderVoided)
	}
}

func TestAdminSalesOfAnEmptyDate(t *testing.T) {
	service := newTestService(t)

	sales, err := service.GetAdminSales("2020-01-01")
	if err != nil {
		t.Fatalf("admin sales: %v", err)
	}
	if sales.BusinessDate != "2020-01-01" || sales.Summary.OrderCount != 0 {
		t.Errorf("sales = %+v, want an empty 2020-01-01", sales.Summary)
	}
	if sales.Orders == nil || sales.Items == nil {
		t.Errorf("orders and items must encode as [], not null")
	}
}

func TestKitchenBoardShowsTheNewestOrders(t *testing.T) {
	service := newTestService(t)

	for _, clientOrderID := range []string{"kitchen-1", "kitchen-2"} {
		request := validOrder()
		request["clientOrderId"] = clientOrderID
		request["notes"] = "no onions"
		postOrder(t, service, request)
	}

	board, err := service.GetKitchenBoard()
	if err != nil {
		t.Fatalf("kitchen board: %v", err)
	}

	if len(board.Tickets) != 2 {
		t.Fatalf("tickets = %d, want 2", len(board.Tickets))
	}
	if board.Tickets[0].OrderNumber != 101 {
		t.Errorf("first ticket = %d, want 101", board.Tickets[0].OrderNumber)
	}
	if board.Tickets[0].Notes != "no onions" {
		t.Errorf("notes = %q, want \"no onions\"", board.Tickets[0].Notes)
	}
	if board.Tickets[0].Lines[0].Name != "Potato Pancake" {
		t.Errorf("line name = %q", board.Tickets[0].Lines[0].Name)
	}
	if board.Tickets[0].Lines[0].Side != "" {
		t.Errorf("a plain line must carry no side, got %q", board.Tickets[0].Lines[0].Side)
	}
}

func TestKitchenBoardShowsTheSideOfALine(t *testing.T) {
	service := newTestService(t)

	request := validOrder()
	request["items"] = []map[string]any{{"menuItemId": "potato-pancake", "quantity": 1, "side": "applesauce"}}
	postOrder(t, service, request)

	board, err := service.GetKitchenBoard()
	if err != nil {
		t.Fatalf("kitchen board: %v", err)
	}
	if board.Tickets[0].Lines[0].Side != "Applesauce" {
		t.Errorf("side = %q, want \"Applesauce\"", board.Tickets[0].Lines[0].Side)
	}
}
