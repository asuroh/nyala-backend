package handler

import (
	"net/http"
	"nyala-backend/usecase"
	"strconv"
)

// LogicHandler ...
type LogicHandler struct {
	Handler
}

// GetFibonacciHandler ...
func (h *LogicHandler) GetFibonacciHandler(w http.ResponseWriter, r *http.Request) {
	n, err := strconv.Atoi(r.URL.Query().Get("n"))
	if err != nil {
		SendBadRequest(w, "Invalid page value")
		return
	}

	logicUC := usecase.LogicUC{ContractUC: h.ContractUC}
	res, err := logicUC.Fibonacci(n)
	if err != nil {
		SendBadRequest(w, err.Error())
		return
	}

	SendSuccess(w, res, nil)
	return
}
