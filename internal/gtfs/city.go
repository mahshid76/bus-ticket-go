package gtfs

import (
	"math"
	"regexp"
	"strings"
)

// GTFS has stops, not cities. Flix stop names look like
// "Prague (Bus Station Holešovice)", "Brodnica, Sądowa street" or
// "Berlin central bus station", so the city is the part before the first
// "(" or ",", minus generic words such as "airport" or "central station".
var genericTail = regexp.MustCompile(`(?i)\s+(central\s+)?(bus|train|railway|coach|main|central|south|north|east|west)?\s*(station|terminal|airport|stop|zob|hbf|centrum|centre|center|busstasjon|busstation|busbahnhof|bahnhof|dworzec|stacja|stanica|stajalište|stajaliste|autostanica|aut\.)(\b.*)?$`)
var trailingNum = regexp.MustCompile(`\s+\d+$`)

// Large cities whose stations are named "Rome Tiburtina", "Berlin-Spandau",
// "London Victoria" and so on. Travellers search for the city, not the station.
var majorCities = strings.Fields("Rome London Paris Brussels Berlin Madrid Barcelona Milan Munich Vienna Prague Warsaw Budapest Amsterdam Hamburg Frankfurt Cologne Zurich Lisbon Stockholm Copenhagen Oslo Helsinki Dublin Athens Bucharest Sofia Belgrade Zagreb Naples Turin Venice Lyon Marseille Stuttgart Dresden Leipzig Hannover Nuremberg Krakow Gdansk Wroclaw Poznan Bratislava Ljubljana Seville Bilbao Antwerp Rotterdam Edinburgh Glasgow Manchester Birmingham")

func foldMajor(n string) string {
	for _, m := range majorCities {
		if len(n) > len(m) && strings.HasPrefix(n, m) && (n[len(m)] == ' ' || n[len(m)] == '-') {
			return m
		}
	}
	return n
}

func cityName(stopName string) string {
	return foldMajor(rawCityName(stopName))
}

func rawCityName(stopName string) string {
	n := stopName
	if i := strings.IndexAny(n, "(,"); i > 0 {
		n = n[:i]
	}
	n = trailingNum.ReplaceAllString(strings.TrimSpace(n), "")
	for i := 0; i < 3; i++ {
		m := genericTail.ReplaceAllString(n, "")
		m = strings.TrimSpace(m)
		if m == "" || m == n {
			break
		}
		n = m
	}
	n = trailingNum.ReplaceAllString(n, "")
	if n == "" {
		n = strings.TrimSpace(stopName)
	}
	return n
}

func haversineKm(lat1, lon1, lat2, lon2 float64) float64 {
	const r = 6371.0
	p := math.Pi / 180
	a := math.Sin((lat2-lat1)*p/2)*math.Sin((lat2-lat1)*p/2) +
		math.Cos(lat1*p)*math.Cos(lat2*p)*math.Sin((lon2-lon1)*p/2)*math.Sin((lon2-lon1)*p/2)
	return 2 * r * math.Asin(math.Sqrt(a))
}

// groupCities assigns every stop to a city. Stops with the same derived name
// are merged when they are within 40 km of each other, so two unrelated towns
// that share a name stay separate. A multi-word name such as "Berlin Wannsee"
// folds into "Berlin" when that city exists within 30 km.
func groupCities(stops []Stop) []City {
	type cluster struct {
		name     string
		stops    []int32
		lat, lon float64
	}
	var clusters []*cluster
	byName := map[string][]*cluster{}

	place := func(i int32, name string, maxKm float64) *cluster {
		s := stops[i]
		for _, c := range byName[name] {
			if haversineKm(c.lat, c.lon, s.Lat, s.Lon) <= maxKm {
				n := float64(len(c.stops))
				c.lat = (c.lat*n + s.Lat) / (n + 1)
				c.lon = (c.lon*n + s.Lon) / (n + 1)
				c.stops = append(c.stops, i)
				return c
			}
		}
		c := &cluster{name: name, stops: []int32{i}, lat: s.Lat, lon: s.Lon}
		clusters = append(clusters, c)
		byName[name] = append(byName[name], c)
		return c
	}

	names := make([]string, len(stops))
	for i := range stops {
		names[i] = cityName(stops[i].Name)
	}
	// first pass: exact names
	for i := range stops {
		place(int32(i), names[i], 40)
	}
	// second pass: fold "Berlin Wannsee" into "Berlin" when "Berlin" is nearby
	for _, c := range clusters {
		w := strings.Fields(c.name)
		if len(w) < 2 {
			continue
		}
		for _, cand := range byName[w[0]] {
			if cand != c && haversineKm(cand.lat, cand.lon, c.lat, c.lon) <= 30 {
				cand.stops = append(cand.stops, c.stops...)
				c.stops = nil
				break
			}
		}
	}

	var out []City
	for _, c := range clusters {
		if len(c.stops) == 0 {
			continue
		}
		out = append(out, City{Name: c.name, Stops: c.stops})
	}
	for ci := range out {
		for _, s := range out[ci].Stops {
			stops[s].City = int32(ci)
		}
	}
	return out
}
