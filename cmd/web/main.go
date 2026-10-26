package main

import (
	"log"
	"os"
	"time"

	"bus-ticket-go/internal/application/search"
	httpd "bus-ticket-go/internal/delivery/http"
	"bus-ticket-go/internal/gtfs"
	"bus-ticket-go/internal/provider/flix"
)

func main() {
	dataPath := env("DATA_PATH", "data/index.gob")
	ix, err := gtfs.Load(dataPath)
	if err != nil {
		log.Fatalf("%v\nrun `go run ./cmd/indexer` first to build %s", err, dataPath)
	}
	log.Printf("loaded %d cities, %d trips (feed ends %d)", len(ix.Cities), len(ix.Trips), ix.FeedEnd)

	f := flix.New(ix)
	uc := search.New(5*time.Second, f) // add more providers here
	r, err := httpd.New(f, uc)
	if err != nil {
		log.Fatal(err)
	}
	addr := ":" + env("PORT", "8000")
	log.Printf("listening on %s", addr)
	log.Fatal(r.Run(addr))
}

func env(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
