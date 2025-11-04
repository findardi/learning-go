package service

import (
	"context"
	"go-restapi/internal/model"
	"go-restapi/internal/repository"
	"go-restapi/pkg/common"
	formattime "go-restapi/pkg/common/format-time"
	"go-restapi/pkg/common/token"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo *repository.Queries
}

func NewAuthService(repo *repository.Queries) *AuthService {
	return &AuthService{
		repo: repo,
	}
}

func (a *AuthService) hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	return string(bytes), err
}

func (a *AuthService) comparePassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (a *AuthService) Register(ctx context.Context, request *model.RequestUserRegister) (string, error) {
	hash, err := a.hashPassword(request.Password)
	if err != nil {
		return "", common.ErrGeneratePassword
	}

	result, err := a.repo.UserInsert(ctx, repository.UserInsertParams{
		Username:  request.Username,
		Email:     request.Email,
		Password:  hash,
		CreatedAt: formattime.Now(),
		UpdatedAt: formattime.Now(),
	})

	if err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			if strings.Contains(err.Error(), "username") {
				return "", common.ErrDuplicateUsername
			}
			if strings.Contains(err.Error(), "email") {
				return "", common.ErrDuplicateEmail
			}
		}
		return "", err
	}

	return result, nil
}

func (a *AuthService) Login(ctx context.Context, request *model.RequestUserLogin) (*model.ResponseUserLogin, error) {
	user, err := a.repo.UserGetByEmail(ctx, request.Email)
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			return nil, common.ErrInvalidCredentials
		}
		return nil, err
	}

	match := a.comparePassword(request.Passowrd, user.Password)
	if !match {
		return nil, common.ErrInvalidCredentials
	}

	validToken, err := token.GenerateToken(user.ID, user.Username, "USER")
	if err != nil {
		return nil, err
	}

	return &model.ResponseUserLogin{
		Username: user.Username,
		Token:    validToken,
	}, nil
}
