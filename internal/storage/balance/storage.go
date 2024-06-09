package balance

import (
	"context"
	"database/sql"
	"github.com/egor-zakharov/go-musthave-diploma-tpl/internal/models"
	"time"
)

const timeOut = 500 * time.Millisecond

type storage struct {
	db *sql.DB
}

func New(db *sql.DB) Storage {
	return &storage{db: db}
}

func (s *storage) GetBalance(ctx context.Context, userID string) (float64, error) {
	ctx, cancel := context.WithTimeout(ctx, timeOut)
	defer cancel()
	var balance float64

	row := s.db.QueryRowContext(ctx, `SELECT sum(sum) FROM balances WHERE user_id=$1`, userID)
	var nullBalance sql.NullFloat64

	err := row.Scan(&nullBalance)
	if err != nil {
		return balance, err
	}

	balance = nullBalance.Float64
	return balance, nil
}

func (s *storage) SetBalance(ctx context.Context, sum float64, userID string) error {
	ctx, cancel := context.WithTimeout(ctx, timeOut)
	defer cancel()

	_, err := s.db.ExecContext(ctx,
		`INSERT into balances (sum, user_id) values ($1,$2) on conflict (user_id) DO UPDATE set sum = balances.sum + $1 where balances.user_id=$2`, sum, userID)

	return err
}

func (s *storage) GetWithdrawal(ctx context.Context, userID string) (float64, error) {
	ctx, cancel := context.WithTimeout(ctx, timeOut)
	defer cancel()
	var withdraw float64

	row := s.db.QueryRowContext(ctx, `SELECT SUM(sum) FROM withdrawals WHERE user_id=$1`, userID)
	var nullWithdraw sql.NullFloat64

	err := row.Scan(&nullWithdraw)
	if err != nil {
		return withdraw, err
	}

	withdraw = nullWithdraw.Float64
	return withdraw, nil
}

func (s *storage) AddWithdraw(ctx context.Context, withdraw models.Withdrawal, userID string) error {
	ctx, cancel := context.WithTimeout(ctx, timeOut)
	defer cancel()

	_, err := s.db.ExecContext(ctx,
		`INSERT INTO withdrawals(number, sum, user_id) VALUES ($1, $2, $3)`, withdraw.OrderNumber, withdraw.Sum, userID)
	return err
}
