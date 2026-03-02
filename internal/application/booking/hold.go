package booking

import "bus-ticket-go/internal/domain"

type HoldGateway interface {
	Hold(tripID int) (*domain.Ticket, error)
}

type HoldSeatUseCase struct {
	gateway HoldGateway
}

func NewHoldSeatUseCase(g HoldGateway) *HoldSeatUseCase {
	return &HoldSeatUseCase{gateway: g}
}

func (uc *HoldSeatUseCase) Execute(tripID int) (*domain.Ticket, error) {
	return uc.gateway.Hold(tripID)
}