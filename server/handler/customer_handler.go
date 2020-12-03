package handler

import (
	"net/http"
	"nyala-backend/server/request"
	"nyala-backend/usecase"

	validator "gopkg.in/go-playground/validator.v9"
)

// CustomerHandler ...
type CustomerHandler struct {
	Handler
}

// LoginHandler ...
func (h *CustomerHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	req := request.CustomerLoginRequest{}
	if err := h.Handler.Bind(r, &req); err != nil {
		SendBadRequest(w, err.Error())
		return
	}
	if err := h.Handler.Validate.Struct(req); err != nil {
		h.SendRequestValidationError(w, err.(validator.ValidationErrors))
		return
	}

	customerUC := usecase.CustomerUC{ContractUC: h.ContractUC}
	res, err := customerUC.Login(req)
	if err != nil {
		SendBadRequest(w, err.Error())
		return
	}

	SendSuccess(w, res, nil)
	return
}

// RefreshTokenHandler ...
func (h *CustomerHandler) RefreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	user := requestIDFromContextInterface(r.Context(), "user")
	customerID := user["id"].(string)

	customerUC := usecase.CustomerUC{ContractUC: h.ContractUC}
	res, err := customerUC.RefreshToken(customerID)
	if err != nil {
		SendBadRequest(w, err.Error())
		return
	}

	SendSuccess(w, res, nil)
	return
}

// RegisterHandler ...
func (h *CustomerHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	req := request.CustomerRequest{}
	if err := h.Handler.Bind(r, &req); err != nil {
		SendBadRequest(w, err.Error())
		return
	}
	if err := h.Handler.Validate.Struct(req); err != nil {
		h.SendRequestValidationError(w, err.(validator.ValidationErrors))
		return
	}

	customerUC := usecase.CustomerUC{ContractUC: h.ContractUC}
	res, err := customerUC.Register(&req)
	if err != nil {
		SendBadRequest(w, err.Error())
		return
	}

	SendSuccess(w, res, nil)
	return
}
