package spec

import "time"

type Ticket struct {
	ID        int       `json:"id"`
	Status    string    `json:"status"`
	Trip      Trip      `json:"trip"`
	ExpiresAt time.Time `json:"expires_at"`
}
type HoldTicketRequest struct {
	TripID int `json:"tripId"`
}
