package orders

import (
	"context"
	"errors"
	"github.com/egor-zakharov/go-musthave-diploma-tpl/internal/models"
)

var (
	ErrDuplicate        = errors.New("duplicate order")
	ErrOrderAnotherUser = errors.New("order created by another user")
	ErrLuhn             = errors.New("luhn error")
)

type Service interface {
	Add(ctx context.Context, orderID string, userID string) error
	GetAllByUser(ctx context.Context, userID string) (*[]models.Order, error)
	Get(ctx context.Context, orderID string, userID string) (*models.Order, error)
}
