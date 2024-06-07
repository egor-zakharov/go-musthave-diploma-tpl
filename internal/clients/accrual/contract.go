package accrual

import (
	"errors"
	"github.com/egor-zakharov/go-musthave-diploma-tpl/internal/dto"
)

var (
	ErrAccrualServerError     = errors.New("accrual server error")
	ErrAccrualTooManyRequests = errors.New("too many requests to accrual")
	ErrAccrualNoData          = errors.New("order is not registered")
)

type Client interface {
	SendOrder(orderID string) (*dto.AccrualOrderResponse, error)
}
