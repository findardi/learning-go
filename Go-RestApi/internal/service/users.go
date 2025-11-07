package service

import (
	"context"
	"go-restapi/internal/model"
	"go-restapi/internal/repository"
	"go-restapi/pkg/common/appmiddleware"
	formattime "go-restapi/pkg/common/format-time"
	"go-restapi/pkg/common/logger"

	"go.uber.org/zap"
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
		logger.Error("failed to get user from context", zap.String("path", "users/profile"), zap.Error(err))
		return nil, err
	}

	result, err := u.repo.UserGetByID(ctx, payload.UserID)
	if err != nil {
		logger.Error("failed to get user by id", zap.String("path", "users/profile"), zap.Error(err))
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
