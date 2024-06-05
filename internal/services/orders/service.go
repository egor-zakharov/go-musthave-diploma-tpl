package orders

import (
	"context"
	"errors"
	"github.com/EClaesson/go-luhn"
	"github.com/egor-zakharov/go-musthave-diploma-tpl/internal/models"
	"github.com/egor-zakharov/go-musthave-diploma-tpl/internal/storage/orders"
	"go.uber.org/zap"
)

type service struct {
	log     *zap.Logger
	storage orders.Storage
}

func New(log *zap.Logger, storage orders.Storage) Service {
	return &service{log: log, storage: storage}
}

func (s *service) Add(ctx context.Context, orderID string, userID string) error {
	ok := s.isValid(orderID)
	if !ok {
		return ErrLuhn
	}

	order, err := s.storage.Add(ctx, orderID, userID)
	if errors.Is(err, orders.ErrConflict) {
		err = ErrDuplicate
	}
	if order.UserID != userID {
		err = ErrOrderAnotherUser
	}
	return err
}

func (s *service) isValid(orderID string) bool {
	ok, _ := luhn.IsValid(orderID)
	return ok
}

func (s *service) Get(ctx context.Context, userID string) (*[]models.Order, error) {
	return s.storage.Get(ctx, userID)
}
