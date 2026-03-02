package nest

import (
	"encoding/json"
	"net/http"
	"strconv"

	"bus-ticket-go/internal/domain"
	"bus-ticket-go/internal/spec"
)

type TicketGateway struct {
	client *Client
}

func NewTicketGateway(client *Client) *TicketGateway {
	return &TicketGateway{client: client}
}

func (g *TicketGateway) HoldSeat(tripID int) error {
	req, err := g.client.newRequest("POST", "/trips/"+strconv.Itoa(tripID)+"/hold", nil)
	if err != nil {
		return err
	}

	_, err = http.DefaultClient.Do(req)
	return err
}

func (g *TicketGateway) Buy(tripID int) (domain.Ticket, error) {
	req, err := g.client.newRequest("POST", "/tickets/"+strconv.Itoa(tripID) , nil)
	if err != nil {
		return domain.Ticket{}, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return domain.Ticket{}, err
	}
	defer resp.Body.Close()

	var specTicket spec.Ticket
	json.NewDecoder(resp.Body).Decode(&specTicket)

	return domain.Ticket{
		ID:     specTicket.ID,
		Status: domain.TicketStatus(specTicket.Status),
	}, nil
}

func (g *TicketGateway) Cancel(ticketID int) error {
	req, err := g.client.newRequest("POST", "/tickets/"+strconv.Itoa(ticketID)+"/cancel", nil)
	if err != nil {
		return err
	}

	_, err = http.DefaultClient.Do(req)
	return err
}
