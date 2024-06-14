package users

import (
	"context"
	"database/sql"
	"errors"
	"github.com/egor-zakharov/go-musthave-diploma-tpl/internal/models"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"time"
)

const timeOut = 500 * time.Millisecond

type storage struct {
	db *sql.DB
}

func New(db *sql.DB) Storage {
	return &storage{db: db}
}

func (s *storage) Register(ctx context.Context, userIn models.User) (*models.User, error) {
	ctx, cancel := context.WithTimeout(ctx, timeOut)
	defer cancel()

	var userID, login, password string
	row := s.db.QueryRowContext(ctx, `INSERT INTO users(login,password) VALUES ($1, $2) returning id, login, password`, userIn.Login, userIn.Password).Scan(&userID, &login, &password)
	var pgErr *pgconn.PgError
	if errors.As(row, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
		return nil, ErrConflict
	}
	user := &models.User{UserID: userID, Login: login, Password: password}
	return user, nil
}

func (s *storage) Login(ctx context.Context, login string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(ctx, timeOut)
	defer cancel()

	row := s.db.QueryRowContext(ctx, `SELECT id, password FROM users WHERE login=$1`, login)
	var id, password string
	err := row.Scan(&id, &password)
	if err != nil {
		return nil, err
	}
	user := &models.User{UserID: id, Login: login, Password: password}
	return user, nil
}
