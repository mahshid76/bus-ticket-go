package http

import (
	"net/http"

	"bus-ticket-go/internal/application/trip"
	"bus-ticket-go/internal/infrastructure/nest"

	"github.com/gin-gonic/gin"
)

type TripHandler struct {
	searchUC *trip.SearchTripsUseCase
}

func NewTripHandler(client *nest.Client) *TripHandler {
	gateway := nest.NewTripGateway(client) // ✅ این مهمه
	uc := trip.NewSearchTripsUseCase(gateway)
	return &TripHandler{searchUC: uc}
}

func (h *TripHandler) ShowTrips(c *gin.Context) {
	trips, err := h.searchUC.Execute()
	if err != nil {
		c.String(500, "cannot load trips")
		return
	}

	c.HTML(http.StatusOK, "trips.html", trips)
}
