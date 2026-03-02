package main

import (
	"bus-ticket-go/internal/delivery/http"
	"bus-ticket-go/internal/infrastructure/nest"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.LoadHTMLGlob("internal/view/templates/*")

	nestClient := nest.NewClient("http://localhost:3000")
	http.RegisterRoutes(r, nestClient)

	r.Run(":8000")

}
