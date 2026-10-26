// Command indexer downloads the FlixBus GTFS feed and writes the compact
// index the web server loads at startup.
package main

import (
	"archive/zip"
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"bus-ticket-go/internal/gtfs"
)

func main() {
	url := flag.String("url", "https://gtfs.gis.flix.tech/gtfs_generic_eu.zip", "GTFS zip URL")
	zipPath := flag.String("zip", "", "use a local zip instead of downloading")
	out := flag.String("out", "data/index.gob", "output index file")
	flag.Parse()

	var data []byte
	var err error
	if *zipPath != "" {
		data, err = os.ReadFile(*zipPath)
	} else {
		data, err = download(*url)
	}
	if err != nil {
		log.Fatal(err)
	}
	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		log.Fatal(err)
	}
	ix, err := gtfs.Build(zr)
	if err != nil {
		log.Fatal(err)
	}
	if err := os.MkdirAll(dirOf(*out), 0o755); err != nil {
		log.Fatal(err)
	}
	if err := ix.Save(*out); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("indexed %d stops, %d cities, %d trips (feed %d to %d)\n",
		len(ix.Stops), len(ix.Cities), len(ix.Trips), ix.FeedStart, ix.FeedEnd)
}

func download(url string) ([]byte, error) {
	c := &http.Client{Timeout: 5 * time.Minute}
	resp, err := c.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("download %s: %s", url, resp.Status)
	}
	return io.ReadAll(resp.Body)
}

func dirOf(p string) string {
	for i := len(p) - 1; i >= 0; i-- {
		if p[i] == '/' {
			return p[:i]
		}
	}
	return "."
}
