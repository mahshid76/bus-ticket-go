package nest

import (
	"encoding/json"
	"net/http"

	"bus-ticket-go/internal/spec"
)

type PaymentGateway struct {
	client *Client
}

func NewPaymentGateway(client *Client) *PaymentGateway {
	return &PaymentGateway{client: client}
}

func (g *PaymentGateway) CreateIntent(ticketID int) (string, error) {
	payload := spec.CreatePaymentIntentRequest{TicketID: ticketID}
	body, _ := json.Marshal(payload)

	req, err := g.client.newRequest("POST", "/payments/intent", body)
	if err != nil {
		return "", err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var res struct {
		ClientSecret string `json:"clientSecret"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return "", err
	}
	return res.ClientSecret, nil
}
