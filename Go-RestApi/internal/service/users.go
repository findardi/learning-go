package service

import (
	"context"
	"go-restapi/internal/model"
	"go-restapi/internal/repository"
	"go-restapi/pkg/common/appmiddleware"
	formattime "go-restapi/pkg/common/format-time"
)

type UserService struct {
	repo *repository.Queries
}

func NewUserService(repo *repository.Queries) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (u *UserService) Profile(ctx context.Context) (*model.ResponseUser, error) {
	payload, err := appmiddleware.GetUserFromContext(ctx)
	if err != nil {
		return nil, err
	}

	result, err := u.repo.UserGetByID(ctx, payload.UserID)
	if err != nil {
		return nil, err
	}

	return &model.ResponseUser{
		Id:        result.ID,
		Username:  result.Username,
		Email:     result.Email,
		CreatedAt: formattime.FormatRFC3339(result.CreatedAt),
		UpdatedAt: formattime.FormatRFC3339(result.UpdatedAt),
	}, nil
}
