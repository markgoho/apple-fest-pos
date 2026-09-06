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

	// TileLabel is what the /pos tile prints instead of Name. It exists so a
	// tile can drop a word the section header already says ("Toastie" under
	// GRILLED CHEESE) without shortening the name on the kitchen ticket or
	// the cart line. Empty means the tile prints Name.
	TileLabel string `json:"tileLabel,omitempty"`

	// Sides is the fixed set of condiments the item can carry. The Operator
	// chooses any number of them, none to all, on the cart line after the
	// tile tap (ADR-0009), so an item with Sides draws one tile and the
	// chosen Sides ride on the cart line.
	Sides []Side `json:"sides,omitempty"`

	// Minor marks an item whose tile is a short strip under the item it
	// belongs to, rather than a full tile of its own. Extra Sour Cream is an
	// add-on to a pancake, not a thing anyone comes to the booth to buy, and
	// a tile the size of the pancake's would say otherwise.
	Minor bool `json:"minor,omitempty"`

	// ChartColorVar names the pos.css custom property (without the leading
	// "--") that draws this item's segment of the Leader Figures tab's hourly
	// revenue chart (issue #12), so the chart and the rest of the theme share
	// one palette.
	ChartColorVar string `json:"-"`
}

// Label gives the text a /pos tile prints for the item.
func (item MenuItem) Label() string {
	if item.TileLabel != "" {
		return item.TileLabel
	}
	return item.Name
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
		ChartColorVar: "apple-red",
		Sides: []Side{
			{ID: "sour-cream", Label: "Sour Cream"},
			{ID: "applesauce", Label: "Applesauce"},
			{ID: "ketchup", Label: "Ketchup"},
		}},
	{ID: "extra-sour-cream", Name: "Extra Sour Cream", TileLabel: "Extra Sour Cream", Category: "Potato Pancakes", PriceCents: 100, SortOrder: 15, PrintGroup: PrintGroupKitchen, Minor: true, ChartColorVar: "cider-gold"},
	{ID: "og-toastie", Name: "OG Toastie", TileLabel: "OG", Category: "Grilled Cheese", PriceCents: 500, SortOrder: 20, PrintGroup: PrintGroupKitchen, ChartColorVar: "gold-tan"},
	{ID: "pizza-toastie", Name: "Pizza Toastie", TileLabel: "Pizza", Category: "Grilled Cheese", PriceCents: 600, SortOrder: 30, PrintGroup: PrintGroupKitchen, ChartColorVar: "leaf-green"},
	{ID: "harvest-toastie", Name: "Harvest Toastie", TileLabel: "Harvest", Category: "Grilled Cheese", PriceCents: 800, SortOrder: 40, PrintGroup: PrintGroupKitchen, ChartColorVar: "gold-tan-ink"},
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

// SideLabels gives the labels of the Sides of one cart line, in the order the
// menu lists them, so a receipt and a kitchen ticket read the same way whatever
// order the Operator tapped.
func SideLabels(menuItemID string, sideIDs []string) []string {
	if len(sideIDs) == 0 {
		return nil
	}
	labels := make([]string, 0, len(sideIDs))
	item, found := menuItemsByID[menuItemID]
	if !found {
		return append(labels, sideIDs...)
	}
	chosen := make(map[string]bool, len(sideIDs))
	for _, id := range sideIDs {
		chosen[id] = true
	}
	for _, side := range item.Sides {
		if chosen[side.ID] {
			labels = append(labels, side.Label)
			delete(chosen, side.ID)
		}
	}
	for _, id := range sideIDs {
		if chosen[id] {
			labels = append(labels, id)
			delete(chosen, id)
		}
	}
	return labels
}
