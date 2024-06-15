package balance

import (
	"context"
	"github.com/egor-zakharov/go-musthave-diploma-tpl/internal/models"
)

//go:generate mockgen -source=contract.go -destination=contract_mock.go -package=balance

type Service interface {
	GetBalance(ctx context.Context, userID string) (float64, error)
	GetSumWithdraw(ctx context.Context, userID string) (float64, error)
	AddWithdraw(ctx context.Context, withdraw models.Withdrawal, userID string) error
	CanWithdraw(ctx context.Context, sum float64, userID string) (bool, error)
	GetAllWithdrawByUser(ctx context.Context, userID string) (*[]models.Withdrawal, error)
}
