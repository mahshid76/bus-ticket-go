package flix

import (
	"context"
	"errors"
	"testing"
	"time"

	"bus-ticket-go/internal/domain"
	"bus-ticket-go/internal/gtfs"
)

// Two cities, one Saturday service, one trip that leaves A at 23:30 and
// arrives in B at 27:00 (03:00 the next day), plus a night trip that starts
// after midnight on the service day.
func fixture() *gtfs.Index {
	var sat gtfs.Service
	sat.Days[time.Saturday] = true
	sat.Start, sat.End = 20261101, 20261130
	return &gtfs.Index{
		FeedEnd: 20261130,
		Stops: []gtfs.Stop{
			{ID: "a1", Name: "A central", TZ: "Europe/Berlin", City: 0},
			{ID: "b1", Name: "B central", TZ: "Europe/Prague", City: 1},
		},
		Cities:   []gtfs.City{{Name: "A", Stops: []int32{0}, Trips: []int32{0, 1}}, {Name: "B", Stops: []int32{1}, Trips: []int32{0, 1}}},
		Routes:   []gtfs.Route{{Name: "FlixBus 1", Agency: "FlixBus"}},
		Services: []gtfs.Service{sat},
		Trips: []gtfs.Trip{
			{Service: 0, Stops: []gtfs.StopTime{{Stop: 0, Dep: 23*3600 + 1800, Pickup: true, DropOff: true}, {Stop: 1, Arr: 27 * 3600, DropOff: true}}},
			{Service: 0, Stops: []gtfs.StopTime{{Stop: 0, Dep: 25 * 3600, Pickup: true}, {Stop: 1, Arr: 28 * 3600, DropOff: true}}},
		},
	}
}

func TestSearchHandlesTimesPastMidnight(t *testing.T) {
	p := New(fixture())
	sat := time.Date(2026, 11, 14, 0, 0, 0, 0, time.UTC)

	offers, err := p.Search(context.Background(), domain.Query{From: "a", To: "B", Date: sat})
	if err != nil {
		t.Fatal(err)
	}
	if len(offers) != 1 || offers[0].Departure.Format("15:04") != "23:30" {
		t.Fatalf("Saturday: want the 23:30 trip only, got %+v", offers)
	}
	if got := offers[0].Duration(); got != 3*time.Hour+30*time.Minute {
		t.Errorf("duration = %v, want 3h30m", got)
	}

	// The 25:00 trip belongs to Saturday's service but leaves at 01:00 on Sunday.
	sun := sat.AddDate(0, 0, 1)
	offers, _ = p.Search(context.Background(), domain.Query{From: "A", To: "B", Date: sun})
	if len(offers) != 1 || offers[0].Departure.Format("15:04") != "01:00" {
		t.Errorf("Sunday: want the 01:00 trip, got %+v", offers)
	}
}

func TestSearchUnknownCity(t *testing.T) {
	_, err := New(fixture()).Search(context.Background(), domain.Query{From: "Nowhere", To: "B"})
	if !errors.Is(err, ErrUnknownCity) {
		t.Errorf("err = %v, want ErrUnknownCity", err)
	}
}

func TestCitiesAutocomplete(t *testing.T) {
	got := New(fixture()).Cities("a", 5)
	if len(got) != 1 || got[0] != "A" {
		t.Errorf("got %v", got)
	}
}
