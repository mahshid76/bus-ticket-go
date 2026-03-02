package ticket

type HoldSeatGateway interface {
	HoldSeat(tripID int) error
}

type HoldSeatUseCase struct {
	gateway HoldSeatGateway
}

func NewHoldSeatUseCase(gateway HoldSeatGateway) *HoldSeatUseCase {
	return &HoldSeatUseCase{gateway: gateway}
}

func (uc *HoldSeatUseCase) Execute(tripID int) error {
	return uc.gateway.HoldSeat(tripID)
}