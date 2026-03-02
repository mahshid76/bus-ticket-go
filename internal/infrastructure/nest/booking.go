package nest

import (
	"encoding/json"
	"net/http"

	"bus-ticket-go/internal/domain"
	"bus-ticket-go/internal/spec"
)

type BookingGateway struct {
	client *Client
}

func NewBookingGateway(client *Client) *BookingGateway {
	return &BookingGateway{client: client}
}

func (g *BookingGateway) Hold(tripID int) (*domain.Ticket, error) {
	payload := spec.HoldTicketRequest{
		TripID: tripID,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	// 👇 body رو به صورت []byte می‌فرستیم
	req, err := g.client.newRequest("POST", "/tickets/hold", body)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var ticket domain.Ticket
	if err := json.NewDecoder(resp.Body).Decode(&ticket); err != nil {
		return nil, err
	}

	return &ticket, nil
}
