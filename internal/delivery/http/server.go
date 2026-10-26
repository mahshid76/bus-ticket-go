// Package http wires the HTML pages and the small JSON API.
package http

import (
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"bus-ticket-go/internal/application/search"
	"bus-ticket-go/internal/domain"
	"bus-ticket-go/internal/provider/flix"
	"bus-ticket-go/internal/view"

	"github.com/gin-gonic/gin"
)

type Server struct {
	flix   *flix.Provider
	search *search.UseCase
	pages  map[string]*template.Template
	now    func() time.Time
}

func New(f *flix.Provider, uc *search.UseCase) (*gin.Engine, error) {
	pages, err := view.Pages("home", "results")
	if err != nil {
		return nil, err
	}
	s := &Server{flix: f, search: uc, pages: pages, now: time.Now}

	r := gin.New()
	r.Use(gin.Recovery(), gin.Logger())
	r.StaticFS("/static", http.FS(view.Static()))
	r.GET("/", s.home)
	r.GET("/search", s.results)
	r.GET("/api/cities", s.cities)
	r.GET("/healthz", func(c *gin.Context) { c.String(200, "ok") })
	return r, nil
}

func (s *Server) render(c *gin.Context, code int, page string, data gin.H) {
	c.Status(code)
	c.Header("Content-Type", "text/html; charset=utf-8")
	if err := s.pages[page].ExecuteTemplate(c.Writer, "layout", data); err != nil {
		log.Printf("render %s: %v", page, err)
	}
}

type route struct{ From, To string }

var popular = []route{
	{"Berlin", "Prague"}, {"Vienna", "Munich"}, {"Paris", "Amsterdam"}, {"Barcelona", "Madrid"},
	{"Milan", "Zurich"}, {"Warsaw", "Krakow"}, {"Lisbon", "Porto"}, {"Budapest", "Vienna"},
}

// dateBounds returns the first and last bookable day: today until the end of the timetable.
func (s *Server) dateBounds() (time.Time, time.Time) {
	n := s.now()
	today := time.Date(n.Year(), n.Month(), n.Day(), 0, 0, 0, 0, time.UTC)
	return today, s.flix.FeedEnd()
}

func (s *Server) home(c *gin.Context) {
	min, max := s.dateBounds()
	tomorrow := min.AddDate(0, 0, 1)
	if tomorrow.After(max) {
		tomorrow = min
	}
	s.render(c, 200, "home", gin.H{
		"Title": "Busscanner · Bus trips across Europe",
		"From":  "Berlin", "To": "Prague",
		"Date": tomorrow.Format("2006-01-02"), "MinDate": min.Format("2006-01-02"), "MaxDate": max.Format("2006-01-02"),
		"Popular": popular,
	})
}

func (s *Server) cities(c *gin.Context) {
	c.Header("Cache-Control", "public, max-age=3600")
	c.JSON(200, s.flix.Cities(c.Query("q"), 8))
}

type chip struct {
	Label string
	URL   string
	On    bool
}

type offerView struct {
	Dep, Arr, Duration, FromStop, ToStop, Provider, Route, BookURL string
	NextDay                                                        bool
	Price                                                          string
}

func (s *Server) results(c *gin.Context) {
	from, to := strings.TrimSpace(c.Query("from")), strings.TrimSpace(c.Query("to"))
	min, max := s.dateBounds()
	date, err := time.Parse("2006-01-02", c.Query("date"))
	if err != nil {
		date = min.AddDate(0, 0, 1)
	}
	if date.Before(min) {
		date = min
	}
	if date.After(max) {
		date = max
	}
	if from == "" || to == "" {
		c.Redirect(http.StatusFound, "/")
		return
	}

	sortKey := domain.SortKey(c.DefaultQuery("sort", string(domain.SortEarliest)))
	tg := domain.TimeGroup(c.Query("time"))
	res, err := s.search.Execute(c.Request.Context(), domain.Query{From: from, To: to, Date: date},
		search.Options{Sort: sortKey, Time: tg})

	data := gin.H{
		"Title": from + " to " + to + " by bus · Busscanner",
		"From":  from, "To": to,
		"Date": date.Format("2006-01-02"), "DateLabel": date.Format("Mon, 2 Jan 2006"),
		"MinDate": min.Format("2006-01-02"), "MaxDate": max.Format("2006-01-02"),
	}
	if err != nil {
		code, msg := 500, "Search failed. Try again in a moment."
		if errors.Is(err, flix.ErrUnknownCity) {
			code, msg = 404, err.Error()+". Pick a city from the suggestions."
		}
		data["Error"], data["Offers"] = msg, []offerView{}
		data["Sorts"], data["Times"] = []chip{}, []chip{}
		s.render(c, code, "results", data)
		return
	}

	link := func(sort domain.SortKey, t domain.TimeGroup) string {
		v := url.Values{"from": {from}, "to": {to}, "date": {date.Format("2006-01-02")}}
		if sort != domain.SortEarliest {
			v.Set("sort", string(sort))
		}
		if t != "" {
			v.Set("time", string(t))
		}
		return "/search?" + v.Encode()
	}
	if sortKey == "" {
		sortKey = domain.SortEarliest
	}
	sorts := []chip{
		{"Earliest", link(domain.SortEarliest, tg), sortKey == domain.SortEarliest},
		{"Fastest", link(domain.SortFastest, tg), sortKey == domain.SortFastest},
	}
	if domain.HasPrices(res.Offers) {
		sorts = append([]chip{{"Cheapest", link(domain.SortCheapest, tg), sortKey == domain.SortCheapest}}, sorts...)
	}
	times := []chip{
		{"Any time", link(sortKey, ""), tg == ""},
		{"Morning", link(sortKey, domain.Morning), tg == domain.Morning},
		{"Afternoon", link(sortKey, domain.Afternoon), tg == domain.Afternoon},
		{"Evening", link(sortKey, domain.Evening), tg == domain.Evening},
	}

	views := make([]offerView, 0, len(res.Offers))
	for _, o := range res.Offers {
		v := offerView{
			Dep: o.Departure.Format("15:04"), Arr: o.Arrival.Format("15:04"), Duration: fmtDur(o.Duration()),
			FromStop: o.FromStop, ToStop: o.ToStop, Provider: o.Provider, Route: o.Route, BookURL: o.BookURL,
			NextDay: o.Arrival.Format("2006-01-02") != o.Departure.Format("2006-01-02"),
		}
		views = append(views, v)
	}
	data["Offers"], data["Sorts"], data["Times"] = views, sorts, times
	data["HasPrices"] = domain.HasPrices(res.Offers)
	s.render(c, 200, "results", data)
}

func fmtDur(d time.Duration) string {
	m := int(d.Minutes())
	if m < 0 {
		m = 0
	}
	return fmt.Sprintf("%dh %02dm", m/60, m%60)
}
