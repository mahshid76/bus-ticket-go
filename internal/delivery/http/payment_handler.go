package http

import (
	"net/http"
	"strconv"

	"bus-ticket-go/internal/application/booking"
	"bus-ticket-go/internal/infrastructure/nest"

	"github.com/gin-gonic/gin"
)

type PaymentHandler struct {
	confirmUC *booking.ConfirmPaymentUseCase
	cancelUC  *booking.CancelTicketUseCase
}

func NewPaymentHandler(client *nest.Client) *PaymentHandler {
	return &PaymentHandler{
		confirmUC: booking.NewConfirmPaymentUseCase(nest.NewConfirmGateway(client)),
		cancelUC:  booking.NewCancelTicketUseCase(nest.NewCancelGateway(client)),
	}
}

func (h *PaymentHandler) Confirm(c *gin.Context) {
	id, err := strconv.Atoi(c.PostForm("ticketId"))
	if err != nil {
		c.String(400, "invalid ticket id")
		return
	}

	if err := h.confirmUC.Execute(id); err != nil {
		c.String(500, "payment failed")
		return
	}

	// ✅ بعد از پرداخت موفق
	c.Redirect(http.StatusSeeOther, "/success")
}

func (h *PaymentHandler) Cancel(c *gin.Context) {
	id, err := strconv.Atoi(c.PostForm("ticketId"))
	if err != nil {
		c.String(400, "invalid ticket id")
		return
	}

	_ = h.cancelUC.Execute(id)

	// 🔁 بعد از کنسل → برگرد به سرچ
	c.Redirect(http.StatusSeeOther, "/")
}

