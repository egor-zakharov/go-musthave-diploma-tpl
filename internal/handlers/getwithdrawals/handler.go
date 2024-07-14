package getwithdrawals

import (
	"encoding/json"
	"errors"
	"github.com/egor-zakharov/go-musthave-diploma-tpl/internal/dto"
	"github.com/egor-zakharov/go-musthave-diploma-tpl/internal/middlewares"
	"github.com/egor-zakharov/go-musthave-diploma-tpl/internal/services/balance"
	balanceStorage "github.com/egor-zakharov/go-musthave-diploma-tpl/internal/storage/balance"
	"net/http"
	"time"
)

type Handler struct {
	balance balance.Service
}

func New(balance balance.Service) *Handler {
	return &Handler{balance: balance}
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {

	userID := r.Context().Value(middlewares.ContextUserIDKey).(string)
	withdrawals, err := h.balance.GetAllWithdrawByUser(r.Context(), userID)
	if errors.Is(err, balanceStorage.ErrNotFound) {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	if err != nil {
		http.Error(w, "Cannot get withdrawals", http.StatusInternalServerError)
		return
	}

	// заполняем модель ответа
	var resp []dto.WithdrawalsResponse

	for _, withdrawal := range *withdrawals {
		stringDate := withdrawal.ProcessedAt.Format(time.RFC3339)
		date, _ := time.Parse(time.RFC3339, stringDate)
		resp = append(resp, dto.WithdrawalsResponse{
			OrderNumber: withdrawal.OrderNumber,
			Sum:         withdrawal.Sum,
			ProcessedAt: date,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	enc := json.NewEncoder(w)
	if err := enc.Encode(resp); err != nil {
		return
	}
}
