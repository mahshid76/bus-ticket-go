package gtfs

import (
	"testing"
	"time"
)

func TestCityName(t *testing.T) {
	cases := map[string]string{
		"Prague (Bus Terminal Roztyly)":      "Prague",
		"Brodnica, Sądowa street":            "Brodnica",
		"Berlin central bus station":         "Berlin",
		"Hamburg Airport":                    "Hamburg",
		"Rome Tiburtina":                     "Rome",
		"Berlin-Spandau":                     "Berlin",
		"Beravci 1 aut. stajalište":          "Beravci",
		"Porto San Giorgio (Largo Carducci)": "Porto San Giorgio",
		"Schwerin central train station":     "Schwerin",
		"Bilbao (Airport)":                   "Bilbao",
	}
	for in, want := range cases {
		if got := cityName(in); got != want {
			t.Errorf("cityName(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestGroupCitiesKeepsDistantNamesApart(t *testing.T) {
	stops := []Stop{
		{Name: "Porto (TIC)", Lat: 41.15, Lon: -8.61},
		{Name: "Porto (Airport)", Lat: 41.24, Lon: -8.67},
		{Name: "Porto, Via Roma", Lat: 43.2, Lon: 13.8}, // Porto in Italy
	}
	cities := groupCities(stops)
	if len(cities) != 2 {
		t.Fatalf("got %d cities, want 2 (Portugal and Italy kept apart)", len(cities))
	}
}

func TestServiceActive(t *testing.T) {
	s := Service{Start: 20261101, End: 20261130}
	s.Days[6] = true // Saturdays
	s.Removed = []int32{20261114}
	s.Added = []int32{20261118}
	check := func(y, m, day int, want bool) {
		t.Helper()
		tm := dateFor(y, m, day)
		if got := s.Active(tm); got != want {
			t.Errorf("Active(%d-%02d-%02d) = %v, want %v", y, m, day, got, want)
		}
	}
	check(2026, 11, 7, true)   // a Saturday
	check(2026, 11, 14, false) // removed
	check(2026, 11, 18, true)  // added on a Wednesday
	check(2026, 11, 11, false) // plain Wednesday
	check(2026, 12, 5, false)  // after the range
}

func dateFor(y, m, d int) time.Time { return time.Date(y, time.Month(m), d, 0, 0, 0, 0, time.UTC) }
