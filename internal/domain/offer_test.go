package domain

import (
	"testing"
	"time"
)

func offer(dep string, minutes int, price *Price) Offer {
	d, _ := time.Parse("15:04", dep)
	return Offer{Departure: d, Arrival: d.Add(time.Duration(minutes) * time.Minute), Price: price}
}

func TestSortOffers(t *testing.T) {
	cheap, dear := &Price{Amount: 9}, &Price{Amount: 30}
	offers := []Offer{offer("10:00", 300, dear), offer("08:00", 400, nil), offer("12:00", 250, cheap)}

	SortOffers(offers, SortEarliest)
	if offers[0].Departure.Hour() != 8 {
		t.Errorf("earliest first: got %v", offers[0].Departure)
	}
	SortOffers(offers, SortFastest)
	if offers[0].Duration() != 250*time.Minute {
		t.Errorf("fastest first: got %v", offers[0].Duration())
	}
	SortOffers(offers, SortCheapest)
	if offers[0].Price != cheap || offers[2].Price != nil {
		t.Errorf("cheapest should lead and unpriced should trail: %+v", offers)
	}
}

func TestFilterByTimeGroup(t *testing.T) {
	offers := []Offer{offer("06:00", 60, nil), offer("13:00", 60, nil), offer("22:00", 60, nil), offer("00:30", 60, nil)}
	if n := len(FilterByTimeGroup(offers, Morning)); n != 2 {
		t.Errorf("morning: %d, want 2 (06:00 and 00:30)", n)
	}
	if n := len(FilterByTimeGroup(offers, Evening)); n != 1 {
		t.Errorf("evening: %d, want 1", n)
	}
	if n := len(FilterByTimeGroup(offers, "")); n != 4 {
		t.Errorf("no filter: %d, want 4", n)
	}
}
