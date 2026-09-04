package pos

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"sync"
	"testing"
	"time"
)

func newTestService(t *testing.T) *OrderService {
	t.Helper()

	database, err := OpenDatabase(filepath.Join(t.TempDir(), "pos.sqlite"))
	if err != nil {
		t.Fatalf("open database: %v", err)
	}
	t.Cleanup(func() { database.Close() })

	return &OrderService{
		DB:                  database,
		Printer:             PrinterConfig{Enabled: false, WindowPort: "9100", KitchenPort: "9100"},
		StartingOrderNumber: 100,
		SystemAdminPIN:      "4242",
		LeaderPIN:           "1234",
		Now:                 time.Now,
	}
}

func postOrder(t *testing.T, service *OrderService, body any) (*httptest.ResponseRecorder, map[string]any) {
	t.Helper()

	encoded, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("encode request: %v", err)
	}

	request := httptest.NewRequest(http.MethodPost, "/api/orders", bytes.NewReader(encoded))
	recorder := httptest.NewRecorder()
	service.Handler().ServeHTTP(recorder, request)

	var decoded map[string]any
	if err := json.Unmarshal(recorder.Body.Bytes(), &decoded); err != nil {
		t.Fatalf("decode response %q: %v", recorder.Body.String(), err)
	}
	return recorder, decoded
}

func validOrder() map[string]any {
	return map[string]any{
		"clientOrderId": "client-1",
		"deviceId":      "device-1",
		"payment":       map[string]any{"method": "cash"},
		"items":         []map[string]any{{"menuItemId": "potato-pancake", "quantity": 2}},
	}
}

func TestPlaceOrderPersistsAndDisablesPrinting(t *testing.T) {
	service := newTestService(t)
	recorder, body := postOrder(t, service, validOrder())

	if recorder.Code != http.StatusCreated {
		t.Fatalf("status = %d, want 201", recorder.Code)
	}

	order := body["order"].(map[string]any)
	if order["orderNumber"] != float64(100) {
		t.Errorf("orderNumber = %v, want 100", order["orderNumber"])
	}
	if order["totalCents"] != float64(2000) {
		t.Errorf("totalCents = %v, want 2000", order["totalCents"])
	}

	print := body["print"].(map[string]any)
	if print["customer"] != "disabled" || print["kitchen"] != "disabled" {
		t.Errorf("print = %v, want both disabled", print)
	}
}

func TestPlaceOrderReplaysDuplicateClientIDs(t *testing.T) {
	service := newTestService(t)
	_, first := postOrder(t, service, validOrder())
	recorder, second := postOrder(t, service, validOrder())

	if recorder.Code != http.StatusCreated {
		t.Fatalf("status = %d, want 201", recorder.Code)
	}

	firstJSON, _ := json.Marshal(first)
	secondJSON, _ := json.Marshal(second)
	if !bytes.Equal(firstJSON, secondJSON) {
		t.Errorf("replay = %s, want %s", secondJSON, firstJSON)
	}
}

func TestPlaceOrderRejectsUnknownMenuItems(t *testing.T) {
	service := newTestService(t)
	request := validOrder()
	request["clientOrderId"] = "client-2"
	request["items"] = []map[string]any{{"menuItemId": "unknown-item", "quantity": 1}}

	recorder, body := postOrder(t, service, request)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", recorder.Code)
	}
	if body["error"] != "Unknown menu item: unknown-item" {
		t.Errorf("error = %v", body["error"])
	}
}

func TestPlaceOrderAcceptsAKnownSide(t *testing.T) {
	service := newTestService(t)
	request := validOrder()
	request["items"] = []map[string]any{{"menuItemId": "potato-pancake", "quantity": 1, "side": "applesauce"}}

	recorder, body := postOrder(t, service, request)

	if recorder.Code != http.StatusCreated {
		t.Fatalf("status = %d, want 201: %v", recorder.Code, body)
	}
}

func TestPlaceOrderRejectsAnUnknownSide(t *testing.T) {
	service := newTestService(t)
	request := validOrder()
	request["items"] = []map[string]any{{"menuItemId": "potato-pancake", "quantity": 1, "side": "hot-sauce"}}

	recorder, body := postOrder(t, service, request)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", recorder.Code)
	}
	if body["error"] != "Unknown side for Potato Pancake: hot-sauce" {
		t.Errorf("error = %v", body["error"])
	}
}

func TestPlaceOrderRejectsASideOnAnItemWithNoSides(t *testing.T) {
	service := newTestService(t)
	request := validOrder()
	request["items"] = []map[string]any{{"menuItemId": "og-toastie", "quantity": 1, "side": "ketchup"}}

	recorder, body := postOrder(t, service, request)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", recorder.Code)
	}
	if body["error"] != "Unknown side for OG Toastie: ketchup" {
		t.Errorf("error = %v", body["error"])
	}
}

func TestPlaceOrderRejectsFractionalQuantities(t *testing.T) {
	service := newTestService(t)
	request := validOrder()
	request["items"] = []map[string]any{{"menuItemId": "potato-pancake", "quantity": 2.5}}

	recorder, body := postOrder(t, service, request)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", recorder.Code)
	}
	if body["error"] != "Item quantity must be a positive integer" {
		t.Errorf("error = %v", body["error"])
	}
}

func TestOrderNumbersIncrementWithinOneBusinessDate(t *testing.T) {
	service := newTestService(t)

	for index, clientOrderID := range []string{"a", "b", "c"} {
		request := validOrder()
		request["clientOrderId"] = clientOrderID
		_, body := postOrder(t, service, request)

		order := body["order"].(map[string]any)
		if want := float64(100 + index); order["orderNumber"] != want {
			t.Errorf("orderNumber = %v, want %v", order["orderNumber"], want)
		}
	}
}

func TestOrderNumbersRestartOnANewBusinessDate(t *testing.T) {
	service := newTestService(t)
	yesterday := time.Now().UTC().AddDate(0, 0, -1)
	service.Now = func() time.Time { return yesterday }

	first := validOrder()
	first["clientOrderId"] = "yesterday-1"
	postOrder(t, service, first)

	second := validOrder()
	second["clientOrderId"] = "yesterday-2"
	_, body := postOrder(t, service, second)
	if body["order"].(map[string]any)["orderNumber"] != float64(101) {
		t.Fatalf("second order number = %v, want 101", body["order"].(map[string]any)["orderNumber"])
	}

	service.Now = time.Now
	today := validOrder()
	today["clientOrderId"] = "today-1"
	_, todayBody := postOrder(t, service, today)
	if todayBody["order"].(map[string]any)["orderNumber"] != float64(100) {
		t.Errorf("new day order number = %v, want 100", todayBody["order"].(map[string]any)["orderNumber"])
	}
}

func TestCreatedAtMatchesTheJavaScriptISOFormat(t *testing.T) {
	service := newTestService(t)
	_, body := postOrder(t, service, validOrder())

	createdAt := body["order"].(map[string]any)["createdAt"].(string)
	if _, err := time.Parse(timestampLayout, createdAt); err != nil {
		t.Fatalf("createdAt %q does not parse: %v", createdAt, err)
	}
	if got := len(createdAt); got != 24 || createdAt[23] != 'Z' {
		t.Errorf("createdAt = %q, want a 24 character string that ends in Z", createdAt)
	}
}

func TestReprintOrderReturns404ForAnUnknownID(t *testing.T) {
	service := newTestService(t)

	_, err := service.ReprintOrder("does-not-exist")
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("err = %v, want ErrNotFound", err)
	}
}

func TestReprintOrderResendsWithPrinterDisabled(t *testing.T) {
	service := newTestService(t)
	_, body := postOrder(t, service, validOrder())
	orderID := body["order"].(map[string]any)["id"].(string)

	response, err := service.ReprintOrder(orderID)
	if err != nil {
		t.Fatalf("reprint order: %v", err)
	}
	if response.Print.Customer != PrintDisabled || response.Print.Kitchen != PrintDisabled {
		t.Errorf("print = %+v, want both disabled", response.Print)
	}
	if response.Order.OrderNumber != 100 {
		t.Errorf("orderNumber = %d, want 100", response.Order.OrderNumber)
	}
}

func TestReprintOrderRouteReturns200ForAKnownID(t *testing.T) {
	service := newTestService(t)
	_, body := postOrder(t, service, validOrder())
	orderID := body["order"].(map[string]any)["id"].(string)

	request := httptest.NewRequest(http.MethodPost, "/api/orders/"+orderID+"/reprint", nil)
	recorder := httptest.NewRecorder()
	service.Handler().ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200: %s", recorder.Code, recorder.Body.String())
	}
}

func TestReprintOrderRouteReturns404ForAnUnknownID(t *testing.T) {
	service := newTestService(t)

	request := httptest.NewRequest(http.MethodPost, "/api/orders/does-not-exist/reprint", nil)
	recorder := httptest.NewRecorder()
	service.Handler().ServeHTTP(recorder, request)

	if recorder.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", recorder.Code)
	}
}

func TestVoidOrderReturns404ForAnUnknownID(t *testing.T) {
	service := newTestService(t)

	_, err := service.VoidOrder("does-not-exist")
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("err = %v, want ErrNotFound", err)
	}
}

func TestVoidOrderRejectsAnAlreadyVoidedOrder(t *testing.T) {
	service := newTestService(t)
	_, body := postOrder(t, service, validOrder())
	orderID := body["order"].(map[string]any)["id"].(string)

	if _, err := service.VoidOrder(orderID); err != nil {
		t.Fatalf("void order: %v", err)
	}

	_, err := service.VoidOrder(orderID)
	if !errors.Is(err, ErrValidation) {
		t.Fatalf("err = %v, want ErrValidation", err)
	}
}

func TestConcurrentDuplicateClientIDsSellOnce(t *testing.T) {
	service := newTestService(t)

	const attempts = 8
	responses := make(chan PlaceOrderResponse, attempts)
	errs := make(chan error, attempts)

	var start sync.WaitGroup
	start.Add(1)

	var finished sync.WaitGroup
	for range attempts {
		finished.Add(1)
		go func() {
			defer finished.Done()
			start.Wait()

			response, err := service.PlaceOrder(PlaceOrderRequest{
				ClientOrderID: "double-tap",
				DeviceID:      "device-1",
				Payment:       Payment{Method: "cash"},
				Items:         []CartLine{{MenuItemID: "potato-pancake", Quantity: 2}},
			})
			if err != nil {
				errs <- err
				return
			}
			responses <- response
		}()
	}

	start.Done()
	finished.Wait()
	close(responses)
	close(errs)

	for err := range errs {
		t.Fatalf("place order: %v", err)
	}

	ids := map[string]bool{}
	for response := range responses {
		ids[response.Order.ID] = true
		if response.Order.OrderNumber != 100 {
			t.Errorf("orderNumber = %d, want 100", response.Order.OrderNumber)
		}
	}
	if len(ids) != 1 {
		t.Errorf("stored %d orders, want 1", len(ids))
	}

	var rows int
	if err := service.DB.QueryRow(`SELECT COUNT(*) FROM transactions`).Scan(&rows); err != nil {
		t.Fatalf("count rows: %v", err)
	}
	if rows != 1 {
		t.Errorf("transactions = %d, want 1", rows)
	}
}
