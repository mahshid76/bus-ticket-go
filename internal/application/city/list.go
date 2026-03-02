package city

import "bus-ticket-go/internal/domain"

type CityGateway interface {
	FetchAll() ([]domain.City, error)
}

type ListCitiesUseCase struct {
	gateway CityGateway
}

func NewListCitiesUseCase(gateway CityGateway) *ListCitiesUseCase {
	return &ListCitiesUseCase{gateway: gateway}
}

func (uc *ListCitiesUseCase) Execute() ([]domain.City, error) {
	return uc.gateway.FetchAll()
}
