package createwithdraw

import (
	"encoding/json"
	"errors"
	"github.com/egor-zakharov/go-musthave-diploma-tpl/internal/dto"
	"github.com/egor-zakharov/go-musthave-diploma-tpl/internal/logger"
	"github.com/egor-zakharov/go-musthave-diploma-tpl/internal/middlewares"
	"github.com/egor-zakharov/go-musthave-diploma-tpl/internal/models"
	"github.com/egor-zakharov/go-musthave-diploma-tpl/internal/services/balance"
	"github.com/egor-zakharov/go-musthave-diploma-tpl/internal/services/orders"
	"io"
	"net/http"
)

type Handler struct {
	orders  orders.Service
	balance balance.Service
}

func New(orders orders.Service, balance balance.Service) *Handler {
	return &Handler{orders: orders, balance: balance}
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middlewares.ContextUserIDKey).(string)
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		http.Error(w, "Invalid request content type", http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Incorrect input json", http.StatusInternalServerError)
		return
	}

	requestData := &dto.WithdrawalRequest{}
	err = json.Unmarshal(body, requestData)
	if err != nil {
		http.Error(w, "Incorrect input json", http.StatusBadRequest)
		return
	}

	err = h.orders.Add(r.Context(), requestData.Number, userID)
	if errors.Is(err, orders.ErrLuhn) {
		http.Error(w, "Invalid order number", http.StatusUnprocessableEntity)
		return
	}
	if errors.Is(err, orders.ErrOrderAnotherUser) {
		http.Error(w, "Order created by another user", http.StatusConflict)
		return
	}
	if errors.Is(err, orders.ErrDuplicate) {
		http.Error(w, "Duplicate order", http.StatusOK)
		return
	}

	ok, err := h.balance.CanWithdraw(r.Context(), requestData.Sum, userID)
	if err != nil {
		logger.Log().Sugar().Infow("Can withdraw", err)
		http.Error(w, "Cannot can withdraw", http.StatusInternalServerError)
		return
	}

	if !ok {
		http.Error(w, "Not enough money", http.StatusPaymentRequired)
		return
	}

	withdrawal := models.Withdrawal{
		OrderNumber: requestData.Number,
		Sum:         requestData.Sum,
	}
	err = h.balance.AddWithdraw(r.Context(), withdrawal, userID)
	if err != nil {
		logger.Log().Sugar().Infow("add withdraw", err)
		http.Error(w, "Cannot add withdraw", http.StatusInternalServerError)
		return
	}
}
