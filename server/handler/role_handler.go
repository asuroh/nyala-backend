package handler

import (
	"kriyapeople/usecase"
	"net/http"
)

// RoleHandler ...
type RoleHandler struct {
	Handler
}

// SelectAllHandler ...
func (h *RoleHandler) SelectAllHandler(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")
	by := r.URL.Query().Get("by")
	sort := r.URL.Query().Get("sort")

	uc := usecase.RoleUC{ContractUC: h.ContractUC}
	res, err := uc.SelectAll(search, by, sort)
	if err != nil {
		SendBadRequest(w, err.Error())
		return
	}

	SendSuccess(w, res, nil)
	return
}
