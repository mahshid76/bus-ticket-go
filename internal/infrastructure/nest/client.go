package nest

import (
	"bytes"
	"net/http"
)

type Client struct {
	BaseURL string
	Token   string
}

func NewClient(baseURL string) *Client {
	return &Client{BaseURL: baseURL}
}

func (c *Client) WithToken(token string) {
	c.Token = token
}

func (c *Client) newRequest(method, path string, body []byte) (*http.Request, error) {
	req, err := http.NewRequest(method, c.BaseURL+path, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	if c.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.Token)
	}

	return req, nil
}
