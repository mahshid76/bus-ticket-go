package http

import (
	"fmt"
	"strconv"
	"time"

	"bus-ticket-go/internal/application/booking"
	"bus-ticket-go/internal/infrastructure/nest"

	"github.com/gin-gonic/gin"
)

type BookingHandler struct {
	holdUC *booking.HoldSeatUseCase
}

func NewBookingHandler(client *nest.Client) *BookingHandler {
	gateway := nest.NewBookingGateway(client)
	uc := booking.NewHoldSeatUseCase(gateway)
	return &BookingHandler{holdUC: uc}
}

func (h *BookingHandler) ShowBooking(c *gin.Context) {
	tripIDStr := c.Query("tripId")
	if tripIDStr == "" {
		c.String(400, "tripId is required")
		return
	}

	tripID, err := strconv.Atoi(tripIDStr)
	if err != nil {
		c.String(400, "invalid tripId")
		return
	}

	ticket, err := h.holdUC.Execute(tripID)
	if err != nil {
		c.String(500, "cannot hold seat")
		return
	}

	fmt.Println("TICKET ID:", ticket.ID)
	fmt.Println("EXPIRES AT:", ticket.ExpiresAt)
	fmt.Println("EXPIRES UNIX:", ticket.ExpiresAt.UnixMilli())
	fmt.Println("NOW:", time.Now().UnixMilli())

	c.HTML(200, "booking.html", gin.H{
		"Ticket":    ticket,
		"ExpiresAt": ticket.ExpiresAt,
	})
}
