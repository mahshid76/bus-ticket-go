package nest

import (
	"encoding/json"
	"net/http"
	"strconv"

	"bus-ticket-go/internal/domain"
	"bus-ticket-go/internal/spec"
)

type TripGateway struct {
	client *Client
}

func NewTripGateway(client *Client) *TripGateway {
	return &TripGateway{client: client}
}

func (g *TripGateway) FetchAll() ([]domain.Trip, error) {
	req, err := g.client.newRequest("GET", "/trips", nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var specs []spec.Trip
	if err := json.NewDecoder(resp.Body).Decode(&specs); err != nil {
		return nil, err
	}

	trips := make([]domain.Trip, 0, len(specs))

	for _, t := range specs {
		price, _ := strconv.ParseFloat(t.Price, 64)

		trips = append(trips, domain.Trip{
			ID: t.ID,
			Route: domain.Route{
				Origin: domain.City{
					Name: t.Route.OriginCity.Name,
				},
				Destination: domain.City{
					Name: t.Route.DestinationCity.Name,
				},
			},
			PriceEUR:       price,
			AvailableSeats: t.AvailableSeats,
		})
	}

	return trips, nil
}
