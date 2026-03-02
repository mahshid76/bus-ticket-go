package nest

import (
	"encoding/json"
	"net/http"

	"bus-ticket-go/internal/domain"
	"bus-ticket-go/internal/spec"
)

type CityGateway struct {
	client *Client
}

func NewCityGateway(client *Client) *CityGateway {
	return &CityGateway{client: client}
}

func (g *CityGateway) FetchAll() ([]domain.City, error) {
	req, err := g.client.newRequest("GET", "/cities", nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var specs []spec.City
	if err := json.NewDecoder(resp.Body).Decode(&specs); err != nil {
		return nil, err
	}

	cities := make([]domain.City, 0, len(specs))
	for _, c := range specs {
		cities = append(cities, domain.City{
			ID:   c.ID,
			Name: c.Name,
		})
	}

	return cities, nil
}
