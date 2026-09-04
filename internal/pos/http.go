package pos

import (
	"encoding/json"
	"errors"
	"fmt"
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
	mux.HandleFunc("GET /leader", service.handleLeaderScreen)
	mux.HandleFunc("POST /leader", service.handleLeaderUnlock)
	mux.HandleFunc("POST /leader/orders/{id}/void", service.handleLeaderVoidOrder)
	mux.HandleFunc("GET /system-admin", service.handleSystemAdminScreen)
	mux.HandleFunc("POST /system-admin", service.handleSystemAdminUnlock)
	mux.HandleFunc("POST /system-admin/start-event", service.handleSystemAdminStartEvent)
	mux.HandleFunc("POST /system-admin/reset", service.handleSystemAdminReset)
	mux.HandleFunc("POST /api/orders", service.handlePlaceOrder)
	mux.HandleFunc("POST /api/orders/{id}/reprint", service.handleReprintOrder)
	mux.HandleFunc("GET /api/kitchen", service.handleKitchen)
	return mux
}

func handleHome(writer http.ResponseWriter, request *http.Request) {
	render(writer, "home.html", page{Title: "Apple Fest POS", BodyClass: "theme", Kiosk: true})
}

// handlePOSScreen draws the menu grid. The cart is client-side state, so the
// server sends the menu once and the script does the rest.
func handlePOSScreen(writer http.ResponseWriter, request *http.Request) {
	var sections []menuSection
	for _, item := range MenuItems {
		var tiles []menuTile
		if len(item.Sides) == 0 {
			tiles = append(tiles, menuTile{MenuItemID: item.ID, Name: item.Name, Label: item.Label(), PriceCents: item.PriceCents})
		} else {
			tiles = append(tiles, menuTile{MenuItemID: item.ID, Name: item.Name, Label: item.Label(), PriceCents: item.PriceCents, SideLabel: "Plain"})
			for _, side := range item.Sides {
				tiles = append(tiles, menuTile{MenuItemID: item.ID, Name: item.Name, Label: item.Label(), PriceCents: item.PriceCents, SideID: side.ID, SideLabel: side.Label})
			}
		}

		if len(sections) > 0 && sections[len(sections)-1].Category == item.Category {
			sections[len(sections)-1].Tiles = append(sections[len(sections)-1].Tiles, tiles...)
			continue
		}
		sections = append(sections, menuSection{Category: item.Category, Tiles: tiles})
	}
	render(writer, "pos.html", posPage{
		page:         page{Title: "Cashier POS", BodyClass: "theme", Kiosk: true},
		MenuSections: sections,
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
		page:         page{Title: "Kitchen display", BodyClass: "theme", Kiosk: true},
		KitchenBoard: board,
	})
}

// handleLeaderScreen always shows the blank PIN form: ADR-0006 keeps no
// session, so a fresh visit never remembers a prior unlock.
func (service *OrderService) handleLeaderScreen(writer http.ResponseWriter, request *http.Request) {
	render(writer, "leader.html", leaderPage{page: page{Title: "Leader", BodyClass: "theme"}})
}

func (service *OrderService) handleLeaderUnlock(writer http.ResponseWriter, request *http.Request) {
	service.renderLeader(writer, request.FormValue("pin"), "", "figures")
}

// handleLeaderVoidOrder voids a Placed order (CONTEXT.md, issue #45) and
// redraws the Leader page on the Orders tab, where every void starts.
func (service *OrderService) handleLeaderVoidOrder(writer http.ResponseWriter, request *http.Request) {
	pin := request.FormValue("pin")
	var message string
	if service.leaderUnlocked(pin) {
		switch row, err := service.VoidOrder(request.PathValue("id")); {
		case errors.Is(err, ErrNotFound):
			message = "Order not found."
		case errors.Is(err, ErrValidation):
			message = validationMessage(err)
		case err != nil:
			log.Printf("void order: %v", err)
			message = "Could not void the order."
		default:
			message = fmt.Sprintf("Order #%d voided.", row.OrderNumber)
		}
	}
	service.renderLeader(writer, pin, message, "orders")
}

// leaderUnlocked reports whether pin matches the configured PIN. An unset
// LEADER_PIN never unlocks, so a missing deploy config fails closed instead
// of leaving sales and void open to an empty submit.
func (service *OrderService) leaderUnlocked(pin string) bool {
	return service.LeaderPIN != "" && pin == service.LeaderPIN
}

// renderLeader checks pin against the configured PIN and draws the Leader
// page: the blank form again on a wrong PIN, or the unlocked sales figures
// and order list with pin carried into every void form on a correct one.
func (service *OrderService) renderLeader(writer http.ResponseWriter, pin string, message string, tab string) {
	if !service.leaderUnlocked(pin) {
		render(writer, "leader.html", leaderPage{
			page:  page{Title: "Leader", BodyClass: "theme"},
			Error: "Wrong PIN.",
		})
		return
	}

	sales, err := service.GetAdminSales("")
	if err != nil {
		log.Printf("leader sales: %v", err)
		http.Error(writer, "Could not read the sales", http.StatusInternalServerError)
		return
	}
	render(writer, "leader.html", leaderPage{
		page:               page{Title: "Leader", BodyClass: "theme"},
		Unlocked:           true,
		PIN:                pin,
		Message:            message,
		Tab:                tab,
		AdminSalesResponse: sales,
	})
}

// handleSystemAdminScreen always shows the blank PIN form: ADR-0006 keeps no
// session, so a fresh visit never remembers a prior unlock.
func (service *OrderService) handleSystemAdminScreen(writer http.ResponseWriter, request *http.Request) {
	render(writer, "system-admin.html", systemAdminPage{page: page{Title: "System Admin", BodyClass: "theme"}})
}

func (service *OrderService) handleSystemAdminUnlock(writer http.ResponseWriter, request *http.Request) {
	service.renderSystemAdmin(writer, request.FormValue("pin"), "")
}

func (service *OrderService) handleSystemAdminStartEvent(writer http.ResponseWriter, request *http.Request) {
	pin := request.FormValue("pin")
	var message string
	if service.systemAdminUnlocked(pin) {
		if err := service.StartEvent(); err != nil {
			log.Printf("start event: %v", err)
			message = "Could not set Start Event."
		}
	}
	service.renderSystemAdmin(writer, pin, message)
}

func (service *OrderService) handleSystemAdminReset(writer http.ResponseWriter, request *http.Request) {
	pin := request.FormValue("pin")
	var message string
	if service.systemAdminUnlocked(pin) {
		switch err := service.ResetAllOrders(); {
		case errors.Is(err, ErrEventStarted):
			message = "Locked: Start Event is set."
		case err != nil:
			log.Printf("reset all orders: %v", err)
			message = "Could not wipe orders."
		default:
			message = "All orders wiped."
		}
	}
	service.renderSystemAdmin(writer, pin, message)
}

// systemAdminUnlocked reports whether pin matches the configured PIN. An
// unset SYSTEM_ADMIN_PIN never unlocks, so a missing deploy config fails
// closed instead of leaving the reset tool open to an empty submit.
func (service *OrderService) systemAdminUnlocked(pin string) bool {
	return service.SystemAdminPIN != "" && pin == service.SystemAdminPIN
}

// renderSystemAdmin checks pin against the configured PIN and draws the
// System Admin page: the blank form again on a wrong PIN, or the unlocked
// tool with pin carried into its action forms on a correct one.
func (service *OrderService) renderSystemAdmin(writer http.ResponseWriter, pin string, message string) {
	if !service.systemAdminUnlocked(pin) {
		render(writer, "system-admin.html", systemAdminPage{
			page:  page{Title: "System Admin", BodyClass: "theme"},
			Error: "Wrong PIN.",
		})
		return
	}

	started, err := service.EventStarted()
	if err != nil {
		log.Printf("event started: %v", err)
		http.Error(writer, "Could not read Start Event", http.StatusInternalServerError)
		return
	}
	render(writer, "system-admin.html", systemAdminPage{
		page:         page{Title: "System Admin", BodyClass: "theme"},
		Unlocked:     true,
		PIN:          pin,
		EventStarted: started,
		Message:      message,
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
