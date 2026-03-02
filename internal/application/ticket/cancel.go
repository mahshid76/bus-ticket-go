package ticket

type CancelTicketGateway interface {
	Cancel(ticketID int) error
}
type CancelTicketUseCase struct {
	gateway CancelTicketGateway
}

func NewCancelTicketUseCase(gateway CancelTicketGateway) *CancelTicketUseCase {
	return &CancelTicketUseCase{gateway: gateway}
}
func (uc *CancelTicketUseCase) Execute(ticketID int) error {
	return uc.gateway.Cancel(ticketID)
}
