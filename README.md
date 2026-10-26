# Busscanner

A bus search engine for Europe, in the spirit of Skyscanner. Pick any two cities and a date,
see every departure, and jump to the operator to book.

Go, server-rendered HTML, no database, no JavaScript framework. It loads in about 40 MB of RAM,
so it runs on a free tier.

## What it does

- Searches real timetables from the public [FlixBus GTFS feed](https://gtfs.gis.flix.tech/gtfs_generic_eu.zip)
  (about 36,000 trips, 1,750 cities, covering the feed's published date range).
- City autocomplete, swap button, sort by earliest or fastest, filter by morning, afternoon or evening.
- Correct durations across time zones and for trips that run past midnight.
- Deep link to the Flix booking page for the chosen stops and date.

## What it does not do (yet)

- **No prices.** GTFS contains schedules, not fares. Flix fares need a partner agreement, so
  results say "Check price on FlixBus". The "Cheapest" sort and price display are already wired
  and appear as soon as a provider returns a `Price`.
- Only one operator. Adding more is one file, see below.
- Only direct trips, no connections.

## Architecture

```
cmd/indexer        downloads the GTFS zip and writes a 1 MB index (shapes.txt is skipped)
cmd/web            loads the index and serves the site
internal/domain    Query, Offer, Provider interface, sorting and time filters
internal/gtfs      GTFS parsing, city grouping, calendar rules
internal/provider  one package per operator (flix)
internal/application/search   runs all providers concurrently, merges, sorts
internal/delivery/http        Gin handlers and templates
```

Dependencies point inward: `domain` knows nothing about GTFS or HTTP.

### Adding a provider

Implement `domain.Provider` (`Name()` and `Search(ctx, Query) ([]Offer, error)`) and pass it to
`search.New(...)` in `cmd/web/main.go`. Providers run concurrently. If one fails, the others
still show. Return `Offer.Price` and the price sort turns on by itself.

### City grouping

GTFS has stops, not cities. `internal/gtfs/city.go` derives the city from the stop name
(`Prague (Bus Terminal Roztyly)` becomes `Prague`), merges stops of the same name within 40 km,
and folds station names of major cities (`Rome Tiburtina`) into the city. It is a heuristic and
is covered by tests, but small towns can still be grouped oddly.

## Run locally

```sh
go run ./cmd/indexer          # downloads the feed, writes data/index.gob (about 3 s)
go run ./cmd/web              # http://localhost:8000
go test ./...
```

Environment: `PORT` (default 8000), `DATA_PATH` (default `data/index.gob`).

## Deploy for free on Render

1. Push this repo to GitHub.
2. On [render.com](https://render.com), choose New, then Blueprint, and select the repo. `render.yaml`
   creates a free Docker web service. No credit card is required.
3. The image downloads the feed at build time, so the timetable is as fresh as the last deploy.
   Trigger a redeploy (or add a weekly scheduled deploy hook) to refresh it.

A free Render service sleeps after 15 minutes without traffic, and the first request after that
takes up to a minute. A free uptime monitor that pings `/healthz` every 10 minutes keeps it awake.
