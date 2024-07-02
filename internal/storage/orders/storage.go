package orders

import (
	"context"
	"database/sql"
	"errors"
	"github.com/egor-zakharov/go-musthave-diploma-tpl/internal/models"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"time"
)

const timeOut = 500 * time.Second

type storage struct {
	db *sql.DB
}

func New(db *sql.DB) Storage {
	return &storage{db: db}
}

func (s *storage) Add(ctx context.Context, orderID string, userID string) (*models.Order, error) {
	ctx, cancel := context.WithTimeout(ctx, timeOut)
	defer cancel()

	var number, usrID string
	var uploadedAt time.Time

	row := s.db.QueryRowContext(ctx, `INSERT INTO orders(number, user_id) VALUES ($1, $2) returning number, user_id, uploaded_at`, orderID, userID).Scan(&number, &usrID, &uploadedAt)
	var pgErr *pgconn.PgError
	if errors.As(row, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
		return nil, ErrConflict
	}

	order := &models.Order{Number: number, UserID: usrID, UploadedAt: uploadedAt}
	return order, nil
}

func (s *storage) GetAllByUser(ctx context.Context, userID string) (*[]models.Order, error) {
	ctx, cancel := context.WithTimeout(ctx, timeOut)
	defer cancel()

	var orders []models.Order

	rows, err := s.db.QueryContext(ctx, `SELECT number, status, accrual, user_id, uploaded_at FROM orders WHERE user_id=$1 order by uploaded_at`, userID)
	if err != nil {
		return nil, err
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {

		var number, status, userID string
		var accrual sql.NullFloat64
		var uploadedAt time.Time

		err = rows.Scan(&number, &status, &accrual, &userID, &uploadedAt)
		if err != nil {
			return nil, err
		}
		order := models.Order{
			Number:     number,
			Status:     status,
			Accrual:    accrual.Float64,
			UserID:     userID,
			UploadedAt: uploadedAt,
		}
		orders = append(orders, order)
	}
	if len(orders) == 0 {
		return nil, ErrNotFound
	}
	return &orders, err
}

func (s *storage) GetAllNotTerminated(ctx context.Context) (*[]models.Order, error) {
	ctx, cancel := context.WithTimeout(ctx, timeOut)
	defer cancel()

	var orders []models.Order

	rows, err := s.db.QueryContext(ctx, `SELECT number, status, accrual, user_id, uploaded_at FROM orders WHERE status not in ('INVALID','PROCESSED')`)
	if err != nil {
		return nil, err
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {

		var number, status, userID string
		var accrual sql.NullFloat64
		var uploadedAt time.Time

		err = rows.Scan(&number, &status, &accrual, &userID, &uploadedAt)
		if err != nil {
			return nil, err
		}
		order := models.Order{
			Number:     number,
			Status:     status,
			UserID:     userID,
			Accrual:    accrual.Float64,
			UploadedAt: uploadedAt,
		}
		orders = append(orders, order)
	}
	return &orders, err
}

func (s *storage) Set(ctx context.Context, order models.Order) error {
	ctx, cancel := context.WithTimeout(ctx, timeOut)
	defer cancel()

	_, err := s.db.ExecContext(ctx,
		`UPDATE orders SET accrual=$1, status=$2 WHERE number=$3`, order.Accrual, order.Status, order.Number)
	return err
}

func (s *storage) Get(ctx context.Context, orderID string) (*models.Order, error) {
	ctx, cancel := context.WithTimeout(ctx, timeOut)
	defer cancel()

	var order *models.Order

	rows := s.db.QueryRowContext(ctx, `SELECT number, status, accrual, user_id, uploaded_at FROM orders WHERE number=$1`, orderID)

	var number, status, userID string
	var accrual sql.NullFloat64
	var uploadedAt time.Time

	err := rows.Scan(&number, &status, &accrual, &userID, &uploadedAt)
	if err != nil {
		return nil, err
	}
	order = &models.Order{
		Number:     number,
		Status:     status,
		Accrual:    accrual.Float64,
		UserID:     userID,
		UploadedAt: uploadedAt,
	}

	return order, err
}
