package pos

// PrintGroup says which printer a menu item belongs to.
type PrintGroup string

const (
	PrintGroupKitchen  PrintGroup = "kitchen"
	PrintGroupCustomer PrintGroup = "customer"
)

// Side is one condiment choice a menu item can carry. See CONTEXT.md.
type Side struct {
	ID    string `json:"id"`
	Label string `json:"label"`
}

// MenuItem is one sellable product. The menu is hard-coded in the binary.
type MenuItem struct {
	ID         string     `json:"id"`
	Name       string     `json:"name"`
	Category   string     `json:"category"`
	PriceCents int        `json:"priceCents"`
	SortOrder  int        `json:"sortOrder"`
	PrintGroup PrintGroup `json:"printGroup"`

	// Sides is the fixed set of condiments the item can carry. The Operator
	// chooses the Side on the menu tile at add time, so an item with Sides
	// draws one tile per Side plus one Plain tile; the side rides on the cart
	// line, not on a separate step.
	Sides []Side `json:"sides,omitempty"`
}

// HasSide reports whether id names one of the item's Sides.
func (item MenuItem) HasSide(id string) bool {
	for _, side := range item.Sides {
		if side.ID == id {
			return true
		}
	}
	return false
}

// MenuItems holds the menu in sort order.
var MenuItems = []MenuItem{
	{ID: "potato-pancake", Name: "Potato Pancake", Category: "Potato Pancakes", PriceCents: 1000, SortOrder: 10, PrintGroup: PrintGroupKitchen,
		Sides: []Side{
			{ID: "sour-cream", Label: "Sour Cream"},
			{ID: "applesauce", Label: "Applesauce"},
			{ID: "ketchup", Label: "Ketchup"},
		}},
	{ID: "og-toastie", Name: "OG Toastie", Category: "Grilled Cheese", PriceCents: 500, SortOrder: 20, PrintGroup: PrintGroupKitchen},
	{ID: "pizza-toastie", Name: "Pizza Toastie", Category: "Grilled Cheese", PriceCents: 600, SortOrder: 30, PrintGroup: PrintGroupKitchen},
	{ID: "harvest-toastie", Name: "Harvest Toastie", Category: "Grilled Cheese", PriceCents: 800, SortOrder: 40, PrintGroup: PrintGroupKitchen},
}

var menuItemsByID = func() map[string]MenuItem {
	byID := make(map[string]MenuItem, len(MenuItems))
	for _, item := range MenuItems {
		byID[item.ID] = item
	}
	return byID
}()

// MenuItemByID finds a menu item by its id.
func MenuItemByID(id string) (MenuItem, bool) {
	item, found := menuItemsByID[id]
	return item, found
}

// MenuItemName gives the item name, or the id when the item is unknown.
func MenuItemName(id string) string {
	if item, found := menuItemsByID[id]; found {
		return item.Name
	}
	return id
}
