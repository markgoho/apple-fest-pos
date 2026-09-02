package pos

import (
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
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
	"cents": FormatCents,
	"clock": FormatClock,
	"asset": func(path string) string { return path + "?v=" + buildVersion },
}

// pageTemplates holds one parsed template per screen, keyed by the template
// file name. html/template finds a missing field only when it runs, so every
// screen has a test that renders it.
var pageTemplates = map[string]*template.Template{
	"home.html":    parsePage("home.html"),
	"pos.html":     parsePage("pos.html"),
	"kitchen.html": parsePage("kitchen.html"),
	"admin.html":   parsePage("admin.html"),
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

// page holds what every screen puts in the shared layout.
type page struct {
	Title     string
	BodyClass string
}

// menuButton is one menu item as the /pos grid needs it.
type menuButton struct {
	MenuItem
	SideList string
}

type posPage struct {
	page
	MenuItems []menuButton
}

type kitchenPage struct {
	page
	KitchenBoard
}

type adminPage struct {
	page
	AdminSalesResponse
}

func render(writer http.ResponseWriter, name string, data any) {
	writer.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := pageTemplates[name].ExecuteTemplate(writer, "base", data); err != nil {
		log.Printf("render %s: %v", name, err)
	}
}
