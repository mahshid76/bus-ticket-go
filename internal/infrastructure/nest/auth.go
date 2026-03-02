package nest

import (
	"encoding/json"
	"net/http"

	"bus-ticket-go/internal/domain"
	"bus-ticket-go/internal/spec"
)

type AuthGateway struct {
	client *Client
}

func NewAuthGateway(client *Client) *AuthGateway {
	return &AuthGateway{client: client}
}

func (g *AuthGateway) Login(email, password string) (domain.UserSession, error) {
	reqBody := spec.LoginRequest{
		Email:    email,
		Password: password,
	}

	data, err := json.Marshal(reqBody)
	if err != nil {
		return domain.UserSession{}, err
	}

	req, err := g.client.newRequest("POST", "/auth/login", data)
	if err != nil {
		return domain.UserSession{}, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return domain.UserSession{}, err
	}
	defer resp.Body.Close()

	var res spec.LoginResponse
	json.NewDecoder(resp.Body).Decode(&res)

	return domain.UserSession{
		AccessToken: res.AccessToken,
	}, nil
}
