package createorder

import (
	"errors"
	"github.com/egor-zakharov/go-musthave-diploma-tpl/internal/middlewares"
	"github.com/egor-zakharov/go-musthave-diploma-tpl/internal/services/orders"
	"io"
	"net/http"
)

type Handler struct {
	order orders.Service
}

func New(order orders.Service) *Handler {
	return &Handler{order: order}
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	if contentType != "text/plain" {
		http.Error(w, "Invalid request content type", http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Incorrect input json", http.StatusInternalServerError)
		return
	}

	orderID := string(body)
	userID := r.Context().Value(middlewares.ContextUserIDKey).(string)
	err = h.order.Add(r.Context(), orderID, userID)
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
	w.WriteHeader(http.StatusAccepted)
}
