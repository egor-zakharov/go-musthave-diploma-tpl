package balance

import (
	"context"
	"github.com/egor-zakharov/go-musthave-diploma-tpl/internal/models"
)

// TODO Для баланса отдельный сторейдж или в orders? потому что работает с таблицей orders и withdrawals
type Storage interface {
	GetBalance(ctx context.Context, userID string) (float64, error)
	GetWithdrawal(ctx context.Context, userID string) (float64, error)
	AddWithdraw(ctx context.Context, withdraw models.Withdrawal, userID string) error
	SetBalance(ctx context.Context, sum float64, userID string) error
}
