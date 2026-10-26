// Package flix answers searches from the public FlixBus/FlixTrain GTFS feed.
// GTFS has schedules but no fares, so offers carry no price and link to the
// Flix booking page instead.
package flix

import (
	"context"
	"fmt"
	"net/url"
	"sort"
	"strings"
	"time"
	_ "time/tzdata" // stop time zones must resolve in minimal containers

	"bus-ticket-go/internal/domain"
	"bus-ticket-go/internal/gtfs"
)

type Provider struct {
	ix     *gtfs.Index
	byName map[string][]int32 // lower-case city name -> city indexes
	locs   map[string]*time.Location
}

func New(ix *gtfs.Index) *Provider {
	p := &Provider{ix: ix, byName: map[string][]int32{}, locs: map[string]*time.Location{}}
	for i, c := range ix.Cities {
		k := strings.ToLower(c.Name)
		p.byName[k] = append(p.byName[k], int32(i))
	}
	return p
}

func (p *Provider) Name() string { return "FlixBus" }

// FeedEnd is the last date the loaded timetable covers.
func (p *Provider) FeedEnd() time.Time { return dateOf(p.ix.FeedEnd) }

func dateOf(d int32) time.Time {
	return time.Date(int(d/10000), time.Month(d/100%100), int(d%100), 0, 0, 0, 0, time.UTC)
}

// Cities returns city names starting with prefix (case-insensitive), the
// busiest first.
func (p *Provider) Cities(prefix string, limit int) []string {
	prefix = strings.ToLower(strings.TrimSpace(prefix))
	type hit struct {
		name  string
		trips int
	}
	var hits []hit
	for _, c := range p.ix.Cities {
		if strings.HasPrefix(strings.ToLower(c.Name), prefix) {
			hits = append(hits, hit{c.Name, len(c.Trips)})
		}
	}
	sort.Slice(hits, func(i, j int) bool {
		if hits[i].trips != hits[j].trips {
			return hits[i].trips > hits[j].trips
		}
		return hits[i].name < hits[j].name
	})
	if len(hits) > limit {
		hits = hits[:limit]
	}
	out := make([]string, len(hits))
	for i, h := range hits {
		out[i] = h.name
	}
	return out
}

// resolve picks the busiest city with this name.
func (p *Provider) resolve(name string) (int32, bool) {
	best, ok := int32(0), false
	for _, c := range p.byName[strings.ToLower(strings.TrimSpace(name))] {
		if !ok || len(p.ix.Cities[c].Trips) > len(p.ix.Cities[best].Trips) {
			best, ok = c, true
		}
	}
	return best, ok
}

func (p *Provider) loc(name string) *time.Location {
	if l, ok := p.locs[name]; ok {
		return l
	}
	l, err := time.LoadLocation(name)
	if err != nil {
		l = time.UTC
	}
	p.locs[name] = l
	return l
}

// at turns a GTFS time (seconds after "noon minus 12h" of the service day)
// into an instant in the stop's own time zone.
func (p *Provider) at(day time.Time, secs int32, tz string) time.Time {
	l := p.loc(tz)
	base := time.Date(day.Year(), day.Month(), day.Day(), 12, 0, 0, 0, l).Add(-12 * time.Hour)
	return base.Add(time.Duration(secs) * time.Second).In(l)
}

var ErrUnknownCity = fmt.Errorf("unknown city")

func (p *Provider) Search(ctx context.Context, q domain.Query) ([]domain.Offer, error) {
	from, ok := p.resolve(q.From)
	if !ok {
		return nil, fmt.Errorf("%w: %q", ErrUnknownCity, q.From)
	}
	to, ok := p.resolve(q.To)
	if !ok {
		return nil, fmt.Errorf("%w: %q", ErrUnknownCity, q.To)
	}
	if from == to {
		return nil, nil
	}
	ix := p.ix
	var offers []domain.Offer
	for n, ti := range ix.Cities[from].Trips {
		if n%2048 == 0 {
			if err := ctx.Err(); err != nil {
				return nil, err
			}
		}
		t := &ix.Trips[ti]
		// find the first boarding stop in the origin city, then a later stop in the destination city
		i := -1
		for k, st := range t.Stops {
			if st.Pickup && ix.Stops[st.Stop].City == from {
				i = k
				break
			}
		}
		if i < 0 {
			continue
		}
		j := -1
		for k := i + 1; k < len(t.Stops); k++ {
			if t.Stops[k].DropOff && ix.Stops[t.Stops[k].Stop].City == to {
				j = k
				break
			}
		}
		if j < 0 {
			continue
		}
		// A trip's service day can be the day before when its times pass midnight.
		day := q.Date
		if t.Stops[i].Dep >= 86400 {
			day = q.Date.AddDate(0, 0, -1)
		}
		if !ix.Services[t.Service].Active(day) {
			continue
		}
		a, b := ix.Stops[t.Stops[i].Stop], ix.Stops[t.Stops[j].Stop]
		offers = append(offers, domain.Offer{
			Provider:  ix.Routes[t.Route].Agency,
			Route:     ix.Routes[t.Route].Name,
			Departure: p.at(day, t.Stops[i].Dep, a.TZ),
			Arrival:   p.at(day, t.Stops[j].Arr, b.TZ),
			FromStop:  a.Name,
			ToStop:    b.Name,
			BookURL:   bookURL(a, b, q.Date),
		})
	}
	return offers, nil
}

// bookURL deep-links into the Flix shop. Flix GTFS stop ids are the station
// UUIDs used by the shop.
func bookURL(from, to gtfs.Stop, day time.Time) string {
	v := url.Values{}
	v.Set("departureStation", from.ID)
	v.Set("arrivalStation", to.ID)
	v.Set("rideDate", day.Format("02.01.2006"))
	v.Set("adult", "1")
	return "https://shop.flixbus.com/search?" + v.Encode()
}
