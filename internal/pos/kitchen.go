package pos

// KitchenTicketLine is one item on a kitchen ticket.
type KitchenTicketLine struct {
	MenuItemID string `json:"menuItemId"`
	Name       string `json:"name"`
	Quantity   int    `json:"quantity"`
	Notes      string `json:"notes,omitempty"`
}

// KitchenTicket is one order on the kitchen display.
type KitchenTicket struct {
	OrderNumber int                 `json:"orderNumber"`
	CreatedAt   string              `json:"createdAt"`
	Notes       string              `json:"notes,omitempty"`
	Lines       []KitchenTicketLine `json:"lines"`
}

// KitchenBoard is the body of GET /api/kitchen.
type KitchenBoard struct {
	ServerTime string          `json:"serverTime"`
	Tickets    []KitchenTicket `json:"tickets"`
}

const maxKitchenTickets = 20

// GetKitchenBoard builds the kitchen display of today. It keeps the newest
// twenty orders and drops the items that only the customer receipt shows.
func (service *OrderService) GetKitchenBoard() (KitchenBoard, error) {
	rows, err := ReadSalesRows(service.DB, CurrentBusinessDate(service.Now))
	if err != nil {
		return KitchenBoard{}, err
	}
	if len(rows) > maxKitchenTickets {
		rows = rows[:maxKitchenTickets]
	}

	tickets := make([]KitchenTicket, 0, len(rows))
	for _, row := range rows {
		request := parseStoredRequest(row.RequestJSON)
		lines := make([]KitchenTicketLine, 0, len(request.Items))

		for _, line := range request.Items {
			if item, found := MenuItemByID(line.MenuItemID); found && item.PrintGroup == PrintGroupCustomer {
				continue
			}
			lines = append(lines, KitchenTicketLine{
				MenuItemID: line.MenuItemID,
				Name:       MenuItemName(line.MenuItemID),
				Quantity:   line.Quantity,
				Notes:      line.Notes,
			})
		}

		if len(lines) == 0 {
			continue
		}
		tickets = append(tickets, KitchenTicket{
			OrderNumber: row.OrderNumber,
			CreatedAt:   row.CreatedAt,
			Notes:       request.Notes,
			Lines:       lines,
		})
	}

	return KitchenBoard{
		ServerTime: service.Now().UTC().Format(timestampLayout),
		Tickets:    tickets,
	}, nil
}
