package http

import (
	"bus-ticket-go/internal/infrastructure/nest"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, client *nest.Client) {
	payment := NewPaymentHandler(client)
	home := NewHomeHandler(client)
	auth := NewAuthHandler(client)
	trips := NewTripHandler(client)
	search := NewSearchHandler(client)
	booking := NewBookingHandler(client)

	r.GET("/", home.ShowHome)
	r.GET("/search", search.DoSearch)

	r.GET("/login", auth.ShowLogin)
	r.POST("/login", auth.DoLogin)

	r.GET("/trips", trips.ShowTrips)

	r.POST("/pay", payment.Confirm)
	r.POST("/cancel", payment.Cancel)

	r.GET("/booking", booking.ShowBooking)

	r.GET("/success", func(c *gin.Context) {
		c.HTML(200, "success.html", nil)

	})
}
