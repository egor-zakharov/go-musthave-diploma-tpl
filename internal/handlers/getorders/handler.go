package getorders

import (
	"encoding/json"
	"errors"
	"github.com/egor-zakharov/go-musthave-diploma-tpl/internal/dto"
	"github.com/egor-zakharov/go-musthave-diploma-tpl/internal/middlewares"
	"github.com/egor-zakharov/go-musthave-diploma-tpl/internal/services/orders"
	orders2 "github.com/egor-zakharov/go-musthave-diploma-tpl/internal/storage/orders"
	"net/http"
	"time"
)

type Handler struct {
	order orders.Service
}

func New(order orders.Service) *Handler {
	return &Handler{order: order}
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middlewares.ContextUserIDKey).(string)
	ords, err := h.order.GetAllByUser(r.Context(), userID)
	if errors.Is(err, orders2.ErrNotFound) {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	if err != nil {
		http.Error(w, "Cannot get orders", http.StatusInternalServerError)
		return
	}

	// заполняем модель ответа
	var resp []dto.GetOrdersResponse

	for _, order := range *ords {
		stringDate := order.UploadedAt.Format(time.RFC3339)
		date, _ := time.Parse(time.RFC3339, stringDate)
		resp = append(resp, dto.GetOrdersResponse{
			Number:     order.Number,
			Status:     order.Status,
			Accrual:    order.Accrual,
			UploadedAt: date,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	enc := json.NewEncoder(w)
	if err := enc.Encode(resp); err != nil {
		return
	}
}
