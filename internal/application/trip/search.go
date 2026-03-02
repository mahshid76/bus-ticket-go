package trip
import "bus-ticket-go/internal/domain"

type TripGateway interface {
	FetchAll() ([]domain.Trip, error)
}
type SearchTripsUseCase struct {
	gateway TripGateway
}
func NewSearchTripsUseCase(gateway TripGateway) *SearchTripsUseCase {
	return &SearchTripsUseCase{gateway: gateway}
}

func (uc *SearchTripsUseCase) Execute() ([]domain.Trip, error) {
	return uc.gateway.FetchAll()
}