package balance

import (
	"context"
	"github.com/egor-zakharov/go-musthave-diploma-tpl/internal/models"
	"github.com/egor-zakharov/go-musthave-diploma-tpl/internal/storage/balance"
	"go.uber.org/zap"
)

type service struct {
	log     *zap.Logger
	storage balance.Storage
}

func New(log *zap.Logger, storage balance.Storage) Service {
	return &service{log: log, storage: storage}
}

func (s *service) GetBalance(ctx context.Context, userID string) (float64, error) {
	return s.storage.GetBalance(ctx, userID)
}

func (s *service) GetSumWithdraw(ctx context.Context, userID string) (float64, error) {
	return s.storage.GetSumWithdrawal(ctx, userID)
}

func (s *service) CanWithdraw(ctx context.Context, sum float64, userID string) (bool, error) {
	getBalance, err := s.storage.GetBalance(ctx, userID)
	if err != nil {
		return false, err
	}
	withdrawal, err := s.storage.GetSumWithdrawal(ctx, userID)
	if err != nil {
		return false, err
	}

	if getBalance+withdrawal < sum {
		return false, nil
	}

	return true, nil
}

func (s *service) AddWithdraw(ctx context.Context, withdraw models.Withdrawal, userID string) error {
	err := s.storage.AddWithdraw(ctx, withdraw, userID)
	if err != nil {
		return err
	}
	newBal := -withdraw.Sum
	err = s.storage.SetBalance(ctx, newBal, userID)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) GetAllWithdrawByUser(ctx context.Context, userID string) (*[]models.Withdrawal, error) {
	return s.storage.GetAllWithdrawByUser(ctx, userID)
}
