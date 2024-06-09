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
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	//TODO по-хорошему переписать на insert returning
	_, err = tx.ExecContext(ctx, `INSERT INTO users(login,password) VALUES ($1, $2)`, userIn.Login, userIn.Password)
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
		return nil, ErrConflict
	}

	user, err := s.Login(ctx, userIn.Login)
	if err != nil {
		return nil, err
	}
	return user, err
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
