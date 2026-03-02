package http

import (
	"net/http"

	"bus-ticket-go/internal/application/auth"
	"bus-ticket-go/internal/infrastructure/nest"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct{
	loginUC * auth.LoginUseCase
}

func NewAuthHandler (client * nest.Client)* AuthHandler{
	gateway := nest.NewAuthGateway(client)
	uc := auth.NewLoginUseCase(gateway)
	return &AuthHandler{loginUC :uc}
}

func (h *AuthHandler) ShowLogin(c *gin.Context){
	c.HTML(http.StatusOK, "login.html", nil)
}

func (h *AuthHandler) DoLogin(c * gin.Context){
	email := c.PostForm("email")
	password := c.PostForm("password")

	session, err := h.loginUC.Execute(email, password)
	if err != nil  {
		c.String(401, "login failed")
		return
	}

	c.SetCookie("token", session.AccessToken, 3600, "/", "", false, true)

	c.Redirect(http.StatusFound, "/trips")
}