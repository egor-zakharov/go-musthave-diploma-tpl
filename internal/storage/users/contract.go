package users

import (
	"context"
	"errors"
	"github.com/egor-zakharov/go-musthave-diploma-tpl/internal/models"
)

//go:generate mockgen -source=contract.go -destination=contract_mock.go -package=users

var (
	ErrConflict = errors.New("data conflict")
)

type Storage interface {
	Register(ctx context.Context, userIn models.User) (*models.User, error)
	Login(ctx context.Context, login string) (*models.User, error)
}
