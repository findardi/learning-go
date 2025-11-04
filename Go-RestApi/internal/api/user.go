package api

import (
	"go-restapi/internal/service"
	"go-restapi/pkg/common"
	"net/http"
)

type UserApi struct {
	service *service.Service
}

func NewUserApi(service *service.Service) *UserApi {
	return &UserApi{
		service: service,
	}
}

func (a *UserApi) Profile(w http.ResponseWriter, r *http.Request) {
	result, err := a.service.Users.Profile(r.Context())
	if err != nil {
		common.RespondWithError(w, err)
	}

	common.RespondWithSuccess(w, "Profile retrieved success", result)
}
