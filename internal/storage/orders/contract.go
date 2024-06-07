package orders

import (
	"context"
	"errors"
	"github.com/egor-zakharov/go-musthave-diploma-tpl/internal/models"
)

var (
	ErrConflict = errors.New("data conflict")
	ErrNotFound = errors.New("not found")
)

type Storage interface {
	Add(ctx context.Context, orderID string, userID string) (models.Order, error)
	GetAllByUser(ctx context.Context, userID string) (*[]models.Order, error)
	GetAllNotTerminated(ctx context.Context) (*[]models.Order, error)
	Set(ctx context.Context, order models.Order) error
	Get(ctx context.Context, orderID string, userID string) (*models.Order, error)
}
