package nest

import (
	"net/http"
	"strconv"
)

type CancelGateway struct {
	client *Client
}

func NewCancelGateway(client *Client) *CancelGateway {
	return &CancelGateway{client: client}
}

func (g *CancelGateway) Cancel(ticketID int) error {
	payload := []byte(`{"ticketId":` + strconv.Itoa(ticketID) + `}`)

	req, err := g.client.newRequest("POST", "/tickets/cancel", payload)
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
