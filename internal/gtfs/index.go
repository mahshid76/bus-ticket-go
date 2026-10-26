// Package gtfs turns a GTFS feed into a compact in-memory index that can
// answer "which trips run from city A to city B on date D" quickly.
package gtfs

import (
	"compress/gzip"
	"encoding/gob"
	"fmt"
	"os"
	"time"
)

type Stop struct {
	ID   string
	Name string
	Lat  float64
	Lon  float64
	TZ   string
	City int32
}

type City struct {
	Name  string
	Stops []int32
	Trips []int32 // trips that serve at least one stop of this city
}

type Route struct {
	Name   string
	Agency string
}

type StopTime struct {
	Stop    int32
	Arr     int32 // seconds after the service day's noon minus 12h; may exceed 86400
	Dep     int32
	Pickup  bool
	DropOff bool
}

type Trip struct {
	Route   int32
	Service int32
	Stops   []StopTime
}

type Service struct {
	Days    [7]bool // indexed by time.Weekday
	Start   int32   // yyyymmdd
	End     int32
	Added   []int32
	Removed []int32
}

type Index struct {
	Generated time.Time
	FeedStart int32
	FeedEnd   int32
	Stops     []Stop
	Cities    []City
	Routes    []Route
	Services  []Service
	Trips     []Trip
}

func (ix *Index) Save(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	zw := gzip.NewWriter(f)
	if err := gob.NewEncoder(zw).Encode(ix); err != nil {
		return err
	}
	return zw.Close()
}

func Load(path string) (*Index, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open index: %w", err)
	}
	defer f.Close()
	zr, err := gzip.NewReader(f)
	if err != nil {
		return nil, err
	}
	var ix Index
	if err := gob.NewDecoder(zr).Decode(&ix); err != nil {
		return nil, err
	}
	return &ix, nil
}

// DateInt converts a date to the yyyymmdd form used by GTFS.
func DateInt(t time.Time) int32 {
	return int32(t.Year()*10000 + int(t.Month())*100 + t.Day())
}

// Active reports whether service s runs on the given calendar date.
func (s *Service) Active(t time.Time) bool {
	d := DateInt(t)
	for _, x := range s.Removed {
		if x == d {
			return false
		}
	}
	for _, x := range s.Added {
		if x == d {
			return true
		}
	}
	return d >= s.Start && d <= s.End && s.Days[t.Weekday()]
}
