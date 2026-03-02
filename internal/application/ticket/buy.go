package ticket

import "bus-ticket-go/internal/domain"

type TicketGateway interface {
	BuyTicket(tripID int) (domain.Ticket, error)
}
type BuyTicketUseCase struct {
	gateway TicketGateway
}

func NewBuyTicketUseCase(gateway TicketGateway) *BuyTicketUseCase {
	return &BuyTicketUseCase{gateway: gateway}
}

func (uc *BuyTicketUseCase) Execute(tripID int) (domain.Ticket, error) {
	return uc.gateway.BuyTicket(tripID)
}
