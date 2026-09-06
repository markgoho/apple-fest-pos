package pos

import (
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

//go:embed templates/*.html
var templateFiles embed.FS

// StaticFiles holds the one stylesheet and the one cart script. They ship
// inside the binary, so the Pi gets one file on a deploy.
//
//go:embed static
var StaticFiles embed.FS

// buildVersion changes every time the process starts, which is every deploy.
// Static files carry no cache-control header, so a tablet's browser can keep
// serving a stale script after a redeploy unless the URL itself changes.
var buildVersion = fmt.Sprintf("%d", time.Now().Unix())

var templateFuncs = template.FuncMap{
	"cents":     FormatCents,
	"clock":     FormatClock,
	"join":      strings.Join,
	"asset":     func(path string) string { return path + "?v=" + buildVersion },
	"menuItems": func() []MenuItem { return MenuItems },
	"svgnum": func(value float64) string {
		// Stacking segments by repeated subtraction can leave a segment a hair
		// below zero (float division rarely lands exactly on the total), which
		// rounds to "-0.0" here; SVG reads that fine, but it looks like a bug.
		if formatted := strconv.FormatFloat(value, 'f', 1, 64); formatted != "-0.0" {
			return formatted
		}
		return "0.0"
	},
}

// pageTemplates holds one parsed template per screen, keyed by the template
// file name. html/template finds a missing field only when it runs, so every
// screen has a test that renders it.
var pageTemplates = map[string]*template.Template{
	"home.html":         parsePage("home.html"),
	"pos.html":          parsePage("pos.html"),
	"kitchen.html":      parsePage("kitchen.html"),
	"leader.html":       parsePage("leader.html"),
	"system-admin.html": parsePage("system-admin.html"),
}

func parsePage(name string) *template.Template {
	return template.Must(template.New(name).Funcs(templateFuncs).
		ParseFS(templateFiles, "templates/base.html", "templates/"+name))
}

// FormatCents writes money the way the booth reads it: whole dollars where the
// price is whole dollars.
func FormatCents(cents int) string {
	if cents%100 == 0 {
		return fmt.Sprintf("$%d", cents/100)
	}
	return fmt.Sprintf("$%d.%02d", cents/100, cents%100)
}

// FormatClock turns a stored UTC timestamp into a local clock time. The Pi
// therefore needs the event timezone set.
func FormatClock(timestamp string) string {
	moment, err := time.Parse(timestampLayout, timestamp)
	if err != nil {
		return timestamp
	}
	return moment.In(time.Local).Format("3:04 PM")
}

// page holds what every screen puts in the shared layout. Kiosk gates the
// full-screen/wake-lock lock-down (issue #6): only the Operator's tablets
// (Home, /pos, /kitchen) run it, not pages meant for a personal phone or
// laptop.
type page struct {
	Title     string
	BodyClass string
	Kiosk     bool
}

// menuTile is one tile of the /pos grid. Every item draws exactly one tile:
// the Sides of an item are chosen on the cart line after the tap (ADR-0009),
// so the tile only carries the set the line may offer.
type menuTile struct {
	MenuItemID string
	Name       string
	Label      string
	PriceCents int
	Minor      bool
	Sides      []Side
}

// SidesAttribute packs the tile's Sides into one data attribute, as
// "id:Label" pairs separated by "|". The ids and labels are hard-coded in
// menu.go and hold neither character, so the browser splits it back with no
// escaping and the page needs no inline JSON.
func (tile menuTile) SidesAttribute() string {
	pairs := make([]string, 0, len(tile.Sides))
	for _, side := range tile.Sides {
		pairs = append(pairs, side.ID+":"+side.Label)
	}
	return strings.Join(pairs, "|")
}

// menuSection is one labeled group of tiles on the /pos grid, so the Operator
// sees a clear break between Potato Pancakes and Grilled Cheese.
type menuSection struct {
	Category string
	Tiles    []menuTile
}

type posPage struct {
	page
	MenuSections []menuSection
}

type kitchenPage struct {
	page
	KitchenBoard
}

// leaderPage draws the Leader PIN gate and, once unlocked, the Figures/Orders
// tabs. PIN carries the entered PIN back into every void form on the page, so
// the visit stays unlocked with no session (ADR-0006): nothing is stored
// server-side, but the browser resubmits it with each void.
type leaderPage struct {
	page
	Unlocked bool
	PIN      string
	Error    string
	Message  string
	Tab      string
	// Saturday and Sunday are the event's two business dates (see
	// eventDayPair), the day toggle's two choices.
	Saturday   string
	Sunday     string
	EventTotal AdminSalesEventTotal
	AdminSalesResponse
}

// systemAdminPage draws the System Admin PIN gate and, once unlocked, the
// data-reset tool. PIN carries the entered PIN back into the page's own
// action forms, so the visit stays unlocked with no session (ADR-0006):
// nothing is stored server-side, but the browser resubmits it with each
// action until the page is left and reloaded.
type systemAdminPage struct {
	page
	Unlocked     bool
	PIN          string
	EventStarted bool
	Error        string
	Message      string
}

func render(writer http.ResponseWriter, name string, data any) {
	writer.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := pageTemplates[name].ExecuteTemplate(writer, "base", data); err != nil {
		log.Printf("render %s: %v", name, err)
	}
}
