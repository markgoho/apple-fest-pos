package pos

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
)

// Handler builds the screens, the static files, and the JSON API of the POS.
func (service *OrderService) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.Handle("GET /static/", http.FileServerFS(StaticFiles))
	mux.HandleFunc("GET /{$}", handleHome)
	mux.HandleFunc("GET /pos", handlePOSScreen)
	mux.HandleFunc("GET /kitchen", service.handleKitchenScreen)
	mux.HandleFunc("GET /admin", service.handleAdminScreen)
	mux.HandleFunc("POST /api/orders", service.handlePlaceOrder)
	mux.HandleFunc("POST /api/orders/{id}/reprint", service.handleReprintOrder)
	mux.HandleFunc("GET /api/kitchen", service.handleKitchen)
	mux.HandleFunc("GET /api/admin/sales", service.handleAdminSales)
	return mux
}

func handleHome(writer http.ResponseWriter, request *http.Request) {
	render(writer, "home.html", page{Title: "Apple Fest POS", BodyClass: "theme"})
}

// handlePOSScreen draws the menu grid. The cart is client-side state, so the
// server sends the menu once and the script does the rest.
func handlePOSScreen(writer http.ResponseWriter, request *http.Request) {
	buttons := make([]menuButton, 0, len(MenuItems))
	for _, item := range MenuItems {
		buttons = append(buttons, menuButton{MenuItem: item, SideList: strings.Join(item.Sides, "|")})
	}
	render(writer, "pos.html", posPage{
		page:      page{Title: "Cashier POS", BodyClass: "theme"},
		MenuItems: buttons,
	})
}

// handleKitchenScreen draws the open orders. The screen is read-only; a
// script polls /api/kitchen and redraws in place, because a meta refresh's
// navigation would drop the tablet out of full screen every cycle.
func (service *OrderService) handleKitchenScreen(writer http.ResponseWriter, request *http.Request) {
	board, err := service.GetKitchenBoard()
	if err != nil {
		log.Printf("kitchen board: %v", err)
		http.Error(writer, "Could not read the kitchen board", http.StatusInternalServerError)
		return
	}
	render(writer, "kitchen.html", kitchenPage{
		page:         page{Title: "Kitchen display", BodyClass: "theme"},
		KitchenBoard: board,
	})
}

// handleAdminScreen draws the sales of one event day. An empty date means
// today.
func (service *OrderService) handleAdminScreen(writer http.ResponseWriter, request *http.Request) {
	sales, err := service.GetAdminSales(request.URL.Query().Get("date"))
	if err != nil {
		log.Printf("admin sales: %v", err)
		http.Error(writer, "Could not read the sales", http.StatusInternalServerError)
		return
	}
	render(writer, "admin.html", adminPage{
		page:               page{Title: "Sales", BodyClass: "theme"},
		AdminSalesResponse: sales,
	})
}

func (service *OrderService) handlePlaceOrder(writer http.ResponseWriter, request *http.Request) {
	var body PlaceOrderRequest
	if err := json.NewDecoder(request.Body).Decode(&body); err != nil {
		writeError(writer, http.StatusBadRequest, "Invalid JSON")
		return
	}

	response, err := service.PlaceOrder(body)
	if errors.Is(err, ErrValidation) {
		writeError(writer, http.StatusBadRequest, validationMessage(err))
		return
	}
	if err != nil {
		log.Printf("place order: %v", err)
		writeError(writer, http.StatusInternalServerError, "Could not save the order")
		return
	}

	writeJSON(writer, http.StatusCreated, response)
}

func (service *OrderService) handleReprintOrder(writer http.ResponseWriter, request *http.Request) {
	response, err := service.ReprintOrder(request.PathValue("id"))
	if errors.Is(err, ErrNotFound) {
		writeError(writer, http.StatusNotFound, "Order not found")
		return
	}
	if err != nil {
		log.Printf("reprint order: %v", err)
		writeError(writer, http.StatusInternalServerError, "Could not reprint the order")
		return
	}

	writeJSON(writer, http.StatusOK, response)
}

func (service *OrderService) handleKitchen(writer http.ResponseWriter, request *http.Request) {
	board, err := service.GetKitchenBoard()
	if err != nil {
		log.Printf("kitchen board: %v", err)
		writeError(writer, http.StatusInternalServerError, "Could not read the kitchen board")
		return
	}
	writeJSON(writer, http.StatusOK, board)
}

func (service *OrderService) handleAdminSales(writer http.ResponseWriter, request *http.Request) {
	sales, err := service.GetAdminSales(request.URL.Query().Get("date"))
	if err != nil {
		log.Printf("admin sales: %v", err)
		writeError(writer, http.StatusInternalServerError, "Could not read the sales")
		return
	}
	writeJSON(writer, http.StatusOK, sales)
}

// validationMessage removes the sentinel prefix, so the operator reads only
// the message the TypeScript server sent.
func validationMessage(err error) string {
	message := err.Error()
	if _, rest, found := strings.Cut(message, ": "); found {
		return rest
	}
	return message
}

func writeJSON(writer http.ResponseWriter, status int, body any) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(status)
	if err := json.NewEncoder(writer).Encode(body); err != nil {
		log.Printf("write response: %v", err)
	}
}

func writeError(writer http.ResponseWriter, status int, message string) {
	writeJSON(writer, status, map[string]string{"error": message})
}
