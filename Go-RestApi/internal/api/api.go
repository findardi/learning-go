package api

import (
	"go-restapi/internal/service"
	"go-restapi/pkg/common"
	"net/http"
)

type API struct {
	Auth *AuthAPI
	User *UserApi
}

func NewAPI(service *service.Service) *API {
	return &API{
		Auth: NewAuthAPI(service),
		User: NewUserApi(service),
	}
}

func (a *API) HealthCheck(w http.ResponseWriter, r *http.Request) {
	common.RespondWithSuccess(w, "Server Running", map[string]string{
		"status": "healthy",
	})
}
