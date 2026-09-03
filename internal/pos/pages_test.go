package pos

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func getPage(t *testing.T, service *OrderService, path string) string {
	t.Helper()

	request := httptest.NewRequest(http.MethodGet, path, nil)
	recorder := httptest.NewRecorder()
	service.Handler().ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("GET %s: got %d, want 200", path, recorder.Code)
	}
	return recorder.Body.String()
}

func TestFormatCents(t *testing.T) {
	cases := map[int]string{0: "$0", 500: "$5", 1000: "$10", 1250: "$12.50", 1205: "$12.05"}
	for cents, want := range cases {
		if got := FormatCents(cents); got != want {
			t.Errorf("FormatCents(%d) = %q, want %q", cents, got, want)
		}
	}
}

func TestPOSScreenShowsTheMenu(t *testing.T) {
	body := getPage(t, newTestService(t), "/pos")

	for _, item := range MenuItems {
		if !strings.Contains(body, item.Name) {
			t.Errorf("/pos does not show %q", item.Name)
		}
	}
	// The pancake draws one tile per side plus Plain: 4 pancake tiles + 3
	// toasties = 7 tiles for 4 menu items, per issue #19.
	for _, tag := range []string{"Plain", "Sour Cream", "Applesauce", "Ketchup"} {
		if !strings.Contains(body, `data-side-label="`+tag+`"`) {
			t.Errorf("/pos does not draw a %q tile", tag)
		}
	}
	if strings.Count(body, `data-menu-item-id="potato-pancake"`) != 4 {
		t.Error("/pos does not draw 4 potato-pancake tiles")
	}
	if !strings.Contains(body, "/static/pos.js") {
		t.Error("/pos does not load the cart script")
	}
}

// TestScreensRenderWithAnOrder catches a template field that does not exist,
// because html/template finds that only when the template runs.
func TestScreensRenderWithAnOrder(t *testing.T) {
	service := newTestService(t)
	recorder, _ := postOrder(t, service, PlaceOrderRequest{
		ClientOrderID: "screen-order",
		DeviceID:      "tablet-1",
		Payment:       Payment{Method: "cash"},
		Items:         []CartLine{{MenuItemID: "potato-pancake", Quantity: 1, Notes: "Applesauce"}},
	})
	if recorder.Code != http.StatusCreated {
		t.Fatalf("place order: got %d, want 201", recorder.Code)
	}

	kitchen := getPage(t, service, "/kitchen")
	if !strings.Contains(kitchen, "#100") || !strings.Contains(kitchen, "Applesauce") {
		t.Errorf("/kitchen misses the order: %s", kitchen)
	}

	sales := getPage(t, service, "/admin")
	if !strings.Contains(sales, "$10") || !strings.Contains(sales, "Potato Pancake") {
		t.Errorf("/admin misses the order: %s", sales)
	}

	home := getPage(t, service, "/")
	if !strings.Contains(home, "Kitchen display") {
		t.Error("/ does not link the kitchen display")
	}
}

func TestStaticFilesAreServed(t *testing.T) {
	body := getPage(t, newTestService(t), "/static/pos.css")
	if !strings.Contains(body, "--apple-red") {
		t.Error("/static/pos.css is not the stylesheet")
	}
}
