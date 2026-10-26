package gtfs

import (
	"archive/zip"
	"encoding/csv"
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Build parses the GTFS zip. shapes.txt (about 150 MB) is never read.
func Build(zr *zip.Reader) (*Index, error) {
	files := map[string]*zip.File{}
	for _, f := range zr.File {
		files[f.Name] = f
	}
	read := func(name string, fn func(col func(string) string) error) error {
		f, ok := files[name]
		if !ok {
			return fmt.Errorf("%s missing from feed", name)
		}
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer rc.Close()
		r := csv.NewReader(rc)
		r.FieldsPerRecord = -1
		r.ReuseRecord = true
		head, err := r.Read()
		if err != nil {
			return err
		}
		pos := map[string]int{}
		for i, h := range head {
			pos[strings.TrimPrefix(h, "\xef\xbb\xbf")] = i
		}
		var rec []string
		col := func(k string) string {
			i, ok := pos[k]
			if !ok || i >= len(rec) {
				return ""
			}
			return rec[i]
		}
		for {
			rec, err = r.Read()
			if err == io.EOF {
				return nil
			}
			if err != nil {
				return err
			}
			if err := fn(col); err != nil {
				return err
			}
		}
	}

	ix := &Index{Generated: time.Now().UTC()}

	agencies := map[string]string{}
	if err := read("agency.txt", func(c func(string) string) error {
		agencies[c("agency_id")] = strings.TrimSuffix(c("agency_name"), "-eu")
		return nil
	}); err != nil {
		return nil, err
	}

	if _, ok := files["feed_info.txt"]; ok {
		_ = read("feed_info.txt", func(c func(string) string) error {
			ix.FeedStart = atoi32(c("feed_start_date"))
			ix.FeedEnd = atoi32(c("feed_end_date"))
			return nil
		})
	}

	stopIdx := map[string]int32{}
	if err := read("stops.txt", func(c func(string) string) error {
		lat, _ := strconv.ParseFloat(c("stop_lat"), 64)
		lon, _ := strconv.ParseFloat(c("stop_lon"), 64)
		stopIdx[c("stop_id")] = int32(len(ix.Stops))
		ix.Stops = append(ix.Stops, Stop{ID: c("stop_id"), Name: c("stop_name"), Lat: lat, Lon: lon, TZ: c("stop_timezone")})
		return nil
	}); err != nil {
		return nil, err
	}

	routeIdx := map[string]int32{}
	if err := read("routes.txt", func(c func(string) string) error {
		name := c("route_short_name")
		if name == "" {
			name = c("route_long_name")
		}
		routeIdx[c("route_id")] = int32(len(ix.Routes))
		ix.Routes = append(ix.Routes, Route{Name: name, Agency: agencies[c("agency_id")]})
		return nil
	}); err != nil {
		return nil, err
	}

	svcIdx := map[string]int32{}
	if err := read("calendar.txt", func(c func(string) string) error {
		var s Service
		for i, d := range []string{"sunday", "monday", "tuesday", "wednesday", "thursday", "friday", "saturday"} {
			s.Days[i] = c(d) == "1"
		}
		s.Start, s.End = atoi32(c("start_date")), atoi32(c("end_date"))
		svcIdx[c("service_id")] = int32(len(ix.Services))
		ix.Services = append(ix.Services, s)
		return nil
	}); err != nil {
		return nil, err
	}
	if _, ok := files["calendar_dates.txt"]; ok {
		if err := read("calendar_dates.txt", func(c func(string) string) error {
			i, ok := svcIdx[c("service_id")]
			if !ok {
				return nil
			}
			d := atoi32(c("date"))
			if c("exception_type") == "1" {
				ix.Services[i].Added = append(ix.Services[i].Added, d)
			} else {
				ix.Services[i].Removed = append(ix.Services[i].Removed, d)
			}
			return nil
		}); err != nil {
			return nil, err
		}
	}

	tripIdx := map[string]int32{}
	if err := read("trips.txt", func(c func(string) string) error {
		r, ok1 := routeIdx[c("route_id")]
		s, ok2 := svcIdx[c("service_id")]
		if !ok1 || !ok2 {
			return nil
		}
		tripIdx[c("trip_id")] = int32(len(ix.Trips))
		ix.Trips = append(ix.Trips, Trip{Route: r, Service: s})
		return nil
	}); err != nil {
		return nil, err
	}

	type row struct {
		seq int
		st  StopTime
	}
	rows := make([][]row, len(ix.Trips))
	if err := read("stop_times.txt", func(c func(string) string) error {
		t, ok1 := tripIdx[c("trip_id")]
		s, ok2 := stopIdx[c("stop_id")]
		if !ok1 || !ok2 {
			return nil
		}
		seq, _ := strconv.Atoi(c("stop_sequence"))
		arr, dep := parseHMS(c("arrival_time")), parseHMS(c("departure_time"))
		rows[t] = append(rows[t], row{seq, StopTime{
			Stop: s, Arr: arr, Dep: dep,
			Pickup: c("pickup_type") != "1", DropOff: c("drop_off_type") != "1",
		}})
		return nil
	}); err != nil {
		return nil, err
	}
	for i := range ix.Trips {
		sort.Slice(rows[i], func(a, b int) bool { return rows[i][a].seq < rows[i][b].seq })
		for _, r := range rows[i] {
			ix.Trips[i].Stops = append(ix.Trips[i].Stops, r.st)
		}
	}

	ix.Cities = groupCities(ix.Stops)
	for ti, t := range ix.Trips {
		seen := map[int32]bool{}
		for _, st := range t.Stops {
			c := ix.Stops[st.Stop].City
			if !seen[c] {
				seen[c] = true
				ix.Cities[c].Trips = append(ix.Cities[c].Trips, int32(ti))
			}
		}
	}
	return ix, nil
}

func atoi32(s string) int32 {
	n, _ := strconv.Atoi(strings.TrimSpace(s))
	return int32(n)
}

func parseHMS(s string) int32 {
	p := strings.Split(strings.TrimSpace(s), ":")
	if len(p) < 2 {
		return 0
	}
	h, _ := strconv.Atoi(p[0])
	m, _ := strconv.Atoi(p[1])
	sec := 0
	if len(p) > 2 {
		sec, _ = strconv.Atoi(p[2])
	}
	return int32(h*3600 + m*60 + sec)
}
