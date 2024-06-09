package orders

import (
	"context"
	"errors"
	"github.com/egor-zakharov/go-musthave-diploma-tpl/internal/models"
)

//go:generate mockgen -source=contract.go -destination=contract_mock.go -package=orders

var (
	ErrConflict = errors.New("data conflict")
	ErrNotFound = errors.New("not found")
)

type Storage interface {
	Add(ctx context.Context, orderID string, userID string) (*models.Order, error)
	GetAllByUser(ctx context.Context, userID string) (*[]models.Order, error)
	GetAllNotTerminated(ctx context.Context) (*[]models.Order, error)
	Set(ctx context.Context, order models.Order) error
	Get(ctx context.Context, orderID string) (*models.Order, error)
}
