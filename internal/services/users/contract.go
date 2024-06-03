package users

import (
	"context"
	"errors"
	"github.com/egor-zakharov/go-musthave-diploma-tpl/internal/models"
)

//go:generate mockgen -source=contract.go -destination=contract_mock.go -package=users

var ErrIncorrectData = errors.New("data incorrect")

type Service interface {
	Register(ctx context.Context, userIn models.User) (*models.User, error)
	Login(ctx context.Context, userIn models.User) (*models.User, error)
}
