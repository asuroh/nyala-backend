package handler

import (
	"kriyapeople/server/request"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"

	"kriyapeople/usecase"

	validator "gopkg.in/go-playground/validator.v9"
)

// AdminHandler ...
type AdminHandler struct {
	Handler
}

// LoginHandler ...
func (h *AdminHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	req := request.UserLoginRequest{}
	if err := h.Handler.Bind(r, &req); err != nil {
		SendBadRequest(w, err.Error())
		return
	}
	if err := h.Handler.Validate.Struct(req); err != nil {
		h.SendRequestValidationError(w, err.(validator.ValidationErrors))
		return
	}

	adminUc := usecase.AdminUC{ContractUC: h.ContractUC}
	res, err := adminUc.Login(req)
	if err != nil {
		SendBadRequest(w, err.Error())
		return
	}

	SendSuccess(w, res, nil)
	return
}

// GetAllHandler ...
func (h *AdminHandler) GetAllHandler(w http.ResponseWriter, r *http.Request) {
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		SendBadRequest(w, "Invalid page value")
		return
	}
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		SendBadRequest(w, "Invalid limit value")
		return
	}
	search := r.URL.Query().Get("search")
	by := r.URL.Query().Get("by")
	sort := r.URL.Query().Get("sort")

	adminUc := usecase.AdminUC{ContractUC: h.ContractUC}
	res, p, err := adminUc.FindAll(search, page, limit, by, sort)
	if err != nil {
		SendBadRequest(w, err.Error())
		return
	}

	SendSuccess(w, res, p)
	return
}

// GetByIDHandler ...
func (h *AdminHandler) GetByIDHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		SendBadRequest(w, "Parameter must be filled")
		return
	}

	adminUc := usecase.AdminUC{ContractUC: h.ContractUC}
	res, err := adminUc.FindByID(id, false)
	if err != nil {
		SendBadRequest(w, err.Error())
		return
	}

	SendSuccess(w, res, nil)
	return
}

// CreateHandler ...
func (h *AdminHandler) CreateHandler(w http.ResponseWriter, r *http.Request) {
	req := request.UserRequest{}
	if err := h.Handler.Bind(r, &req); err != nil {
		SendBadRequest(w, err.Error())
		return
	}
	if err := h.Handler.Validate.Struct(req); err != nil {
		h.SendRequestValidationError(w, err.(validator.ValidationErrors))
		return
	}

	adminUc := usecase.AdminUC{ContractUC: h.ContractUC}
	res, err := adminUc.Create(&req)
	if err != nil {
		SendBadRequest(w, err.Error())
		return
	}

	SendSuccess(w, res, nil)
	return
}

// UpdateHandler ...
func (h *AdminHandler) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		SendBadRequest(w, "Parameter must be filled")
		return
	}

	req := request.UserRequest{}
	if err := h.Handler.Bind(r, &req); err != nil {
		SendBadRequest(w, err.Error())
		return
	}
	if err := h.Handler.Validate.Struct(req); err != nil {
		h.SendRequestValidationError(w, err.(validator.ValidationErrors))
		return
	}

	adminUc := usecase.AdminUC{ContractUC: h.ContractUC}
	res, err := adminUc.Update(id, &req)
	if err != nil {
		SendBadRequest(w, err.Error())
		return
	}

	SendSuccess(w, res, nil)
	return
}

// DeleteHandler ...
func (h *AdminHandler) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		SendBadRequest(w, "Parameter must be filled")
		return
	}

	adminUc := usecase.AdminUC{ContractUC: h.ContractUC}
	res, err := adminUc.Delete(id)
	if err != nil {
		SendBadRequest(w, err.Error())
		return
	}

	SendSuccess(w, res, nil)
	return
}
