package api

import (
	"go-restapi/internal/model"
	"go-restapi/internal/service"
	"go-restapi/pkg/common"
	"net/http"
)

type AuthAPI struct {
	service *service.Service
}

func NewAuthAPI(service *service.Service) *AuthAPI {
	return &AuthAPI{
		service: service,
	}
}

func (a *AuthAPI) Register(w http.ResponseWriter, r *http.Request) {
	var request model.RequestUserRegister

	if err := common.ReadJSON(r, &request); err != nil {
		common.ErrorResponseJSON(w, http.StatusBadRequest, "Invalid JSON Format", err.Error())
		return
	}

	if validationErr := common.ValidateRequest(request); validationErr != nil {
		common.ErrorResponseJSON(w, validationErr.Status, validationErr.Message, validationErr.Errors)
		return
	}

	result, err := a.service.Auth.Register(r.Context(), &request)
	if err != nil {
		common.RespondWithError(w, err)
		return
	}

	common.RespondWithSuccess(w, "Success Created User", result)
}

func (a *AuthAPI) Login(w http.ResponseWriter, r *http.Request) {
	var request model.RequestUserLogin

	if err := common.ReadJSON(r, &request); err != nil {
		common.ErrorResponseJSON(w, http.StatusBadRequest, "Invalid JSON Format", err.Error())
		return
	}

	if validationErr := common.ValidateRequest(request); validationErr != nil {
		common.ErrorResponseJSON(w, validationErr.Status, validationErr.Message, validationErr.Errors)
		return
	}

	result, err := a.service.Auth.Login(r.Context(), &request)
	if err != nil {
		common.RespondWithError(w, err)
		return
	}

	common.RespondWithSuccess(w, "Success Login User", result)
}
