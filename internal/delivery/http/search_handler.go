package http

import (
	"net/http"
	"strconv"

	"bus-ticket-go/internal/application/search"
	"bus-ticket-go/internal/domain"
	"bus-ticket-go/internal/infrastructure/nest"

	"github.com/gin-gonic/gin"
)

type SearchHandler struct {
	searchUC *search.SearchTripsUseCase
}

func NewSearchHandler(client *nest.Client) *SearchHandler {
	gateway := nest.NewSearchGateway(client)
	uc := search.NewSearchTripsUseCase(gateway)
	return &SearchHandler{searchUC: uc}
}

func (h *SearchHandler) DoSearch(c *gin.Context) {
	fromID, _ := strconv.Atoi(c.Query("fromCity"))
	toID, _ := strconv.Atoi(c.Query("toCity"))
	date := c.Query("date")

	sortBy := c.Query("sort")
	timeGrp := c.Query("time")

	trips, err := h.searchUC.Execute(fromID, toID, date)
	if err != nil {
		c.String(500, "search failed")
		return
	}

	if sortBy == "price" {
		domain.SortByPrice(trips)
	}

	if timeGrp != "" {
		trips = domain.FilterByTimeGroup(trips, domain.TimeGroup(timeGrp))
	}

	c.HTML(http.StatusOK, "search.html", gin.H{
		"Trips": trips,
		"Date":  date,
	})
}
