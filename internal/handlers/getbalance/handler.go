package getbalance

import (
	"encoding/json"
	"github.com/egor-zakharov/go-musthave-diploma-tpl/internal/dto"
	"github.com/egor-zakharov/go-musthave-diploma-tpl/internal/middlewares"
	"github.com/egor-zakharov/go-musthave-diploma-tpl/internal/services/balance"
	"net/http"
)

type Handler struct {
	balance balance.Service
}

func New(balance balance.Service) *Handler {
	return &Handler{balance: balance}
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middlewares.ContextUserIDKey).(string)
	bal, err := h.balance.GetBalance(r.Context(), userID)
	if err != nil {
		http.Error(w, "Cannot get balance", http.StatusInternalServerError)
		return
	}

	withdrawal, err := h.balance.GetSumWithdraw(r.Context(), userID)
	if err != nil {
		http.Error(w, "Cannot get withdraw", http.StatusInternalServerError)
		return
	}

	// заполняем модель ответа
	resp := dto.GetBalanceResponse{
		Current:   bal,
		Withdrawn: withdrawal,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	enc := json.NewEncoder(w)
	if err := enc.Encode(resp); err != nil {
		return
	}
}
