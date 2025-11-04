package service

import "go-restapi/internal/repository"

type Service struct {
	Users *UserService
	Auth  *AuthService
}

func New(repo *repository.Queries) *Service {
	return &Service{
		Users: NewUserService(repo),
		Auth:  NewAuthService(repo),
	}
}
