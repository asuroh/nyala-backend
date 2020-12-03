package handler

import (
	"net/http"
	"nyala-backend/server/request"
	"nyala-backend/usecase"

	validator "gopkg.in/go-playground/validator.v9"
)

// OrderHandler ...
type OrderHandler struct {
	Handler
}

// CreateHandler ...
func (h *OrderHandler) CreateHandler(w http.ResponseWriter, r *http.Request) {
	user := requestIDFromContextInterface(r.Context(), "user")
	customerID := user["id"].(string)

	req := request.OrderRequest{}
	if err := h.Handler.Bind(r, &req); err != nil {
		SendBadRequest(w, err.Error())
		return
	}
	if err := h.Handler.Validate.Struct(req); err != nil {
		h.SendRequestValidationError(w, err.(validator.ValidationErrors))
		return
	}

	orderUC := usecase.OrderUC{ContractUC: h.ContractUC}
	res, err := orderUC.Create(&req, customerID)
	if err != nil {
		SendBadRequest(w, err.Error())
		return
	}

	SendSuccess(w, res, nil)
	return
}
