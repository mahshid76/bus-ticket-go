package search

import "bus-ticket-go/internal/domain"

type SearchGateway interface {
	Search(FromCityID, toCityID int, date string) ([]domain.Trip, error)
}
type SearchTripsUseCase struct {
	gateway SearchGateway
}

func NewSearchTripsUseCase(gateway SearchGateway) *SearchTripsUseCase {
	return &SearchTripsUseCase{gateway: gateway}
}

func (uc *SearchTripsUseCase) Execute(fromCityID, toCityID int, date string) ([]domain.Trip, error) {
	return uc.gateway.Search(fromCityID, toCityID, date)
}
