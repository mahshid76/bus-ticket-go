package booking

type CancelGateway interface {
	Cancel(ticketID int) error
}

type CancelTicketUseCase struct {
	gateway CancelGateway
}

func NewCancelTicketUseCase(g CancelGateway) *CancelTicketUseCase {
	return &CancelTicketUseCase{gateway: g}
}

func (uc *CancelTicketUseCase) Execute(ticketID int) error {
	return uc.gateway.Cancel(ticketID)
}
