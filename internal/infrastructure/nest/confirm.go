package nest

import (
	"net/http"
	"strconv"
)

type ConfirmGateway struct {
	client *Client
}

func NewConfirmGateway(client *Client) *ConfirmGateway {
	return &ConfirmGateway{client: client}
}

func (g *ConfirmGateway) Confirm(ticketID int) error {
	payload := []byte(`{"ticketId":` + strconv.Itoa(ticketID) + `}`)

	req, err := g.client.newRequest("POST", "/payments/confirm", payload)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
