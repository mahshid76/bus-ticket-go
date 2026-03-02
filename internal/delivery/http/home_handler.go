package http

import (
	"net/http"

	"bus-ticket-go/internal/application/city"
	"bus-ticket-go/internal/infrastructure/nest"

	"github.com/gin-gonic/gin"
)

type HomeHandler struct {
	listCitiesUC *city.ListCitiesUseCase
}

func NewHomeHandler(client *nest.Client) *HomeHandler {
	gateway := nest.NewCityGateway(client)
	uc := city.NewListCitiesUseCase(gateway)
	return &HomeHandler{listCitiesUC: uc}
}

func (h *HomeHandler) ShowHome(c *gin.Context) {
	cities, err := h.listCitiesUC.Execute()
	if err != nil {
		c.String(500, "cannot load cities")
		return
	}

	c.HTML(http.StatusOK, "home.html", gin.H{
		"Cities": cities,
	})
}