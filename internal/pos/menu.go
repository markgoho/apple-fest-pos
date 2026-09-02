package pos

// PrintGroup says which printer a menu item belongs to.
type PrintGroup string

const (
	PrintGroupKitchen  PrintGroup = "kitchen"
	PrintGroupCustomer PrintGroup = "customer"
)

// MenuItem is one sellable product. The menu is hard-coded in the binary.
type MenuItem struct {
	ID         string     `json:"id"`
	Name       string     `json:"name"`
	Category   string     `json:"category"`
	PriceCents int        `json:"priceCents"`
	SortOrder  int        `json:"sortOrder"`
	PrintGroup PrintGroup `json:"printGroup"`
}

// MenuItems holds the menu in sort order.
var MenuItems = []MenuItem{
	{ID: "potato-pancake", Name: "Potato Pancake", Category: "Menu", PriceCents: 1000, SortOrder: 10, PrintGroup: PrintGroupKitchen},
	{ID: "og-toastie", Name: "The OG Toastie", Category: "Grilled Cheese", PriceCents: 500, SortOrder: 20, PrintGroup: PrintGroupKitchen},
	{ID: "pizza-toastie", Name: "The Pizza Toastie", Category: "Grilled Cheese", PriceCents: 600, SortOrder: 30, PrintGroup: PrintGroupKitchen},
	{ID: "harvest-toastie", Name: "The Harvest Toastie", Category: "Grilled Cheese", PriceCents: 800, SortOrder: 40, PrintGroup: PrintGroupKitchen},
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
