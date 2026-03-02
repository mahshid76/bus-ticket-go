package auth

import "bus-ticket-go/internal/domain"

type AuthGetway interface {
	Login(email, password string) (domain.UserSession, error)
}

type LoginUseCase struct {
	gateway AuthGetway
}

func NewLoginUseCase(gateway AuthGetway) *LoginUseCase {
	return &LoginUseCase{gateway: gateway}
}
func (uc *LoginUseCase) Execute(email, password string) (domain.UserSession, error) {
	return uc.gateway.Login(email, password)
}
