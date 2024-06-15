package balance

import (
	"context"
	"errors"
	"github.com/egor-zakharov/go-musthave-diploma-tpl/internal/models"
)

//go:generate mockgen -source=contract.go -destination=contract_mock.go -package=balance

var ErrNotFound = errors.New("not found")

// TODO Для баланса и списаний отдельные сторейдж? потому что работает с таблицей balances и withdrawals
type Storage interface {
	GetBalance(ctx context.Context, userID string) (float64, error)
	GetSumWithdrawal(ctx context.Context, userID string) (float64, error)
	AddWithdraw(ctx context.Context, withdraw models.Withdrawal, userID string) error
	SetBalance(ctx context.Context, sum float64, userID string) error
	GetAllWithdrawByUser(ctx context.Context, userID string) (*[]models.Withdrawal, error)
}
