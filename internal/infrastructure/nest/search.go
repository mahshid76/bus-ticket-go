package nest

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"

	"bus-ticket-go/internal/domain"
	"bus-ticket-go/internal/spec"
)

type SearchGateway struct {
	client *Client
}

func NewSearchGateway(client *Client) *SearchGateway {
	return &SearchGateway{client: client}
}

func (g *SearchGateway) Search(fromCityID, toCityID int, date string) ([]domain.Trip, error) {
	// build query params
	q := url.Values{}
	q.Set("fromCity", strconv.Itoa(fromCityID))
	q.Set("toCity", strconv.Itoa(toCityID))
	q.Set("date", date)

	// create request to Nest backend
	req, err := g.client.newRequest("GET", "/trips?"+q.Encode(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// decode response
	var specs []spec.Trip
	if err := json.NewDecoder(resp.Body).Decode(&specs); err != nil {
		return nil, err
	}

	// map Spec → Domain
	trips := make([]domain.Trip, 0, len(specs))
	for _, t := range specs {
		price, _ := strconv.ParseFloat(t.Price, 64)

		trips = append(trips, domain.Trip{
			ID: t.ID,
			Route: domain.Route{
				Origin: domain.City{
					ID:   t.Route.OriginCity.ID,
					Name: t.Route.OriginCity.Name,
				},
				Destination: domain.City{
					ID:   t.Route.DestinationCity.ID,
					Name: t.Route.DestinationCity.Name,
				},
			},
			DepartureTime:  t.DepartureTime,
			ArrivalTime:    t.ArrivalTime,
			PriceEUR:       price,
			AvailableSeats: t.AvailableSeats,
		})
	}

	return trips, nil
}
