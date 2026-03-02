package booking

type ConfirmGateway interface {
	Confirm(ticketID int) error
}

type ConfirmPaymentUseCase struct {
	gateway ConfirmGateway
}

func NewConfirmPaymentUseCase(g ConfirmGateway) *ConfirmPaymentUseCase {
	return &ConfirmPaymentUseCase{gateway: g}
}

func (uc *ConfirmPaymentUseCase) Execute(ticketID int) error {
	return uc.gateway.Confirm(ticketID)
}
