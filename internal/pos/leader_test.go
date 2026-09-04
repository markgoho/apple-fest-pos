package pos

import (
	"net/url"
	"strings"
	"testing"
	"time"
)

func TestLeaderAlwaysStartsLocked(t *testing.T) {
	service := newTestService(t)

	body := getPage(t, service, "/leader")
	if strings.Contains(body, "Items sold") {
		t.Error("/leader shows sales figures before a PIN is entered")
	}
	if !strings.Contains(body, `name="pin"`) {
		t.Error("/leader does not show a PIN field")
	}
}

func TestLeaderWrongPINStaysLocked(t *testing.T) {
	service := newTestService(t)

	recorder := postForm(t, service, "/leader", url.Values{"pin": {"0000"}})
	body := recorder.Body.String()
	if strings.Contains(body, "Items sold") {
		t.Error("a wrong PIN unlocked the Leader page")
	}
	if !strings.Contains(body, "Wrong PIN") {
		t.Error("a wrong PIN does not show an error")
	}
}

func TestLeaderCorrectPINShowsSalesAndCarriesThePINIntoVoidForms(t *testing.T) {
	service := newTestService(t)
	postOrder(t, service, validOrder())

	recorder := postForm(t, service, "/leader", url.Values{"pin": {service.LeaderPIN}})
	body := recorder.Body.String()
	if !strings.Contains(body, "Items sold") {
		t.Error("the correct PIN does not show the sales figures")
	}
	if !strings.Contains(body, `value="`+service.LeaderPIN+`"`) {
		t.Error("the unlocked page does not carry the PIN into its void forms")
	}
	if !strings.Contains(body, `action="/leader/orders/`) {
		t.Error("the unlocked page does not draw a void form for the order")
	}
}

func TestLeaderEmptyPINNeverUnlocks(t *testing.T) {
	service := newTestService(t)
	service.LeaderPIN = ""

	recorder := postForm(t, service, "/leader", url.Values{"pin": {""}})
	if strings.Contains(recorder.Body.String(), "Items sold") {
		t.Error("an unset LEADER_PIN unlocked the Leader page on an empty submit")
	}
}

func TestLeaderVoidRequiresTheCorrectPIN(t *testing.T) {
	service := newTestService(t)
	_, body := postOrder(t, service, validOrder())
	orderID := body["order"].(map[string]any)["id"].(string)

	postForm(t, service, "/leader/orders/"+orderID+"/void", url.Values{"pin": {"0000"}})

	sales, err := service.GetAdminSales("")
	if err != nil {
		t.Fatalf("admin sales: %v", err)
	}
	if sales.Summary.OrderCount != 1 {
		t.Error("a wrong PIN voided the order")
	}
}

func TestLeaderVoidDropsTheOrderFromTheSummary(t *testing.T) {
	service := newTestService(t)
	_, body := postOrder(t, service, validOrder())
	orderID := body["order"].(map[string]any)["id"].(string)

	recorder := postForm(t, service, "/leader/orders/"+orderID+"/void", url.Values{"pin": {service.LeaderPIN}})
	responseBody := recorder.Body.String()
	if !strings.Contains(responseBody, "voided") {
		t.Errorf("the void response does not confirm the void: %s", responseBody)
	}
	if !strings.Contains(responseBody, "Voided") {
		t.Error("the voided order is not marked Voided in the order list")
	}

	sales, err := service.GetAdminSales("")
	if err != nil {
		t.Fatalf("admin sales: %v", err)
	}
	if sales.Summary.OrderCount != 0 {
		t.Errorf("orderCount after void = %d, want 0", sales.Summary.OrderCount)
	}
}

func TestLeaderVoidOfAnUnknownOrderReportsNotFound(t *testing.T) {
	service := newTestService(t)

	recorder := postForm(t, service, "/leader/orders/does-not-exist/void", url.Values{"pin": {service.LeaderPIN}})
	if !strings.Contains(recorder.Body.String(), "Order not found") {
		t.Error("voiding an unknown order does not report Order not found")
	}
}

func TestLeaderShowsTheDayToggleAndEventTotal(t *testing.T) {
	service := newTestService(t)
	eastern, err := time.LoadLocation("America/New_York")
	if err != nil {
		t.Fatalf("load America/New_York: %v", err)
	}
	original := time.Local
	time.Local = eastern
	defer func() { time.Local = original }()
	service.Now = func() time.Time { return time.Date(2026, 10, 3, 9, 0, 0, 0, eastern) }
	postOrder(t, service, validOrder())

	body := postForm(t, service, "/leader", url.Values{"pin": {service.LeaderPIN}}).Body.String()
	if !strings.Contains(body, `formaction="/leader?date=2026-10-03"`) {
		t.Error("the day toggle does not offer Saturday, 2026-10-03")
	}
	if !strings.Contains(body, `formaction="/leader?date=2026-10-04"`) {
		t.Error("the day toggle does not offer Sunday, 2026-10-04")
	}
	if !strings.Contains(body, "Event total") {
		t.Error("the Figures tab does not show an event total")
	}
}

func TestLeaderDayToggleSwitchesTheViewedDate(t *testing.T) {
	service := newTestService(t)
	eastern, err := time.LoadLocation("America/New_York")
	if err != nil {
		t.Fatalf("load America/New_York: %v", err)
	}
	original := time.Local
	time.Local = eastern
	defer func() { time.Local = original }()
	service.Now = func() time.Time { return time.Date(2026, 10, 3, 9, 0, 0, 0, eastern) }
	postOrder(t, service, validOrder())

	saturday := postForm(t, service, "/leader?date=2026-10-03", url.Values{"pin": {service.LeaderPIN}}).Body.String()
	if !strings.Contains(saturday, "Potato Pancake") {
		t.Error("2026-10-03 does not show the order placed that day")
	}

	sunday := postForm(t, service, "/leader?date=2026-10-04", url.Values{"pin": {service.LeaderPIN}}).Body.String()
	if strings.Contains(sunday, "Potato Pancake") {
		t.Error("2026-10-04 shows an order placed on 2026-10-03")
	}
	if !strings.Contains(sunday, "No items sold yet today") {
		t.Error("2026-10-04 does not show an empty day")
	}
	if !strings.Contains(sunday, "Event total") {
		t.Error("2026-10-04 does not show the event total")
	}
}

func TestLeaderVoidPreservesTheViewedDate(t *testing.T) {
	service := newTestService(t)
	eastern, err := time.LoadLocation("America/New_York")
	if err != nil {
		t.Fatalf("load America/New_York: %v", err)
	}
	original := time.Local
	time.Local = eastern
	defer func() { time.Local = original }()
	service.Now = func() time.Time { return time.Date(2026, 10, 3, 9, 0, 0, 0, eastern) }
	_, body := postOrder(t, service, validOrder())
	orderID := body["order"].(map[string]any)["id"].(string)

	recorder := postForm(t, service, "/leader/orders/"+orderID+"/void", url.Values{
		"pin": {service.LeaderPIN}, "date": {"2026-10-03"},
	})
	if !strings.Contains(recorder.Body.String(), "<span>2026-10-03</span>") {
		t.Error("voiding does not keep showing the date the Leader was viewing")
	}
}
