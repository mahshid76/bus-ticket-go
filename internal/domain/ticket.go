package domain

import "time"

type TicketStatus string

const (
	TicketHeld      TicketStatus = "HELD"
	TicketConfirmed TicketStatus = "CONFIRMED"
	TicketCancelled TicketStatus = "CANCELLED"
)

type Ticket struct {
	ID        int
	Status    TicketStatus
	Trip      Trip
	ExpiresAt time.Time
	CreatedAt time.Time
}