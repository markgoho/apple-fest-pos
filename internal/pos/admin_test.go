package pos

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func postForm(t *testing.T, service *OrderService, path string, form url.Values) *httptest.ResponseRecorder {
	t.Helper()

	request := httptest.NewRequest(http.MethodPost, path, strings.NewReader(form.Encode()))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	recorder := httptest.NewRecorder()
	service.Handler().ServeHTTP(recorder, request)
	return recorder
}

func TestSystemAdminAlwaysStartsLocked(t *testing.T) {
	service := newTestService(t)

	body := getPage(t, service, "/system-admin")
	if strings.Contains(body, "Wipe all orders") {
		t.Error("/system-admin shows the reset tool before a PIN is entered")
	}
	if !strings.Contains(body, `name="pin"`) {
		t.Error("/system-admin does not show a PIN field")
	}
}

func TestSystemAdminWrongPINStaysLocked(t *testing.T) {
	service := newTestService(t)

	recorder := postForm(t, service, "/system-admin", url.Values{"pin": {"0000"}})
	body := recorder.Body.String()
	if strings.Contains(body, "Wipe all orders") {
		t.Error("a wrong PIN unlocked the reset tool")
	}
	if !strings.Contains(body, "Wrong PIN") {
		t.Error("a wrong PIN does not show an error")
	}
}

func TestSystemAdminCorrectPINUnlocksTheResetTool(t *testing.T) {
	service := newTestService(t)

	recorder := postForm(t, service, "/system-admin", url.Values{"pin": {service.SystemAdminPIN}})
	body := recorder.Body.String()
	if !strings.Contains(body, "Wipe all orders") {
		t.Error("the correct PIN does not unlock the reset tool")
	}
	if !strings.Contains(body, `value="`+service.SystemAdminPIN+`"`) {
		t.Error("the unlocked page does not carry the PIN into its action forms")
	}
}

func TestSystemAdminResetWipesOrdersAndRestartsNumbering(t *testing.T) {
	service := newTestService(t)
	postOrder(t, service, PlaceOrderRequest{
		ClientOrderID: "reset-order",
		DeviceID:      "tablet-1",
		Payment:       Payment{Method: "cash"},
		Items:         []CartLine{{MenuItemID: "potato-pancake", Quantity: 1, Side: "applesauce"}},
	})

	recorder := postForm(t, service, "/system-admin/reset", url.Values{"pin": {service.SystemAdminPIN}})
	if !strings.Contains(recorder.Body.String(), "All orders wiped") {
		t.Error("the reset page does not confirm the wipe")
	}

	sales, err := service.GetAdminSales("")
	if err != nil {
		t.Fatalf("read sales after reset: %v", err)
	}
	if sales.Summary.OrderCount != 0 {
		t.Errorf("order count after reset = %d, want 0", sales.Summary.OrderCount)
	}

	_, body := postOrder(t, service, PlaceOrderRequest{
		ClientOrderID: "post-reset-order",
		DeviceID:      "tablet-1",
		Payment:       Payment{Method: "cash"},
		Items:         []CartLine{{MenuItemID: "potato-pancake", Quantity: 1, Side: "applesauce"}},
	})
	order := body["order"].(map[string]any)
	if int(order["orderNumber"].(float64)) != service.StartingOrderNumber {
		t.Errorf("order number after reset = %v, want %d", order["orderNumber"], service.StartingOrderNumber)
	}
}

func TestSystemAdminEmptyPINNeverUnlocks(t *testing.T) {
	service := newTestService(t)
	service.SystemAdminPIN = ""

	recorder := postForm(t, service, "/system-admin/reset", url.Values{"pin": {""}})
	if strings.Contains(recorder.Body.String(), "Wipe all orders") {
		t.Error("an unset SYSTEM_ADMIN_PIN unlocked the reset tool on an empty submit")
	}
}

func TestSystemAdminStartEventLocksOutTheResetTool(t *testing.T) {
	service := newTestService(t)

	postForm(t, service, "/system-admin/start-event", url.Values{"pin": {service.SystemAdminPIN}})

	recorder := postForm(t, service, "/system-admin/reset", url.Values{"pin": {service.SystemAdminPIN}})
	body := recorder.Body.String()
	if !strings.Contains(body, "Locked") {
		t.Error("the reset tool does not report itself locked after Start Event")
	}

	started, err := service.EventStarted()
	if err != nil {
		t.Fatalf("read EventStarted: %v", err)
	}
	if !started {
		t.Error("EventStarted() is false after Start Event was set")
	}
}
