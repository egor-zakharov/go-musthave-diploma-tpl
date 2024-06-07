package accrual

import (
	"context"
	"github.com/egor-zakharov/go-musthave-diploma-tpl/internal/clients/accrual"
	"github.com/egor-zakharov/go-musthave-diploma-tpl/internal/logger"
	"github.com/egor-zakharov/go-musthave-diploma-tpl/internal/models"
	"github.com/egor-zakharov/go-musthave-diploma-tpl/internal/storage/orders"
	"go.uber.org/zap"
	"time"
)

type processor struct {
	log     *zap.Logger
	client  accrual.Client
	storage orders.Storage
}

// TODO горутину тут оставить или вынести вместе с тикером?
func (p processor) Do() {
	go func() {
		ticker := time.NewTicker(500 * time.Millisecond)
		for range ticker.C {
			notTerminated, err := p.storage.GetAllNotTerminated(context.Background())
			if err != nil {
				continue
			}
			for _, order := range *notTerminated {
				accOrder, err := p.client.SendOrder(order.Number)
				if err != nil {
					continue
				}
				updateOrder := models.Order{
					Accrual: accOrder.Accrual,
					Number:  accOrder.Order,
					Status:  accOrder.Status,
				}
				if order.Status == "REGISTERED" {
					logger.Log().Sugar().Infow("Order is just registered", updateOrder)
					updateOrder.Status = "NEW"
				}
				err = p.storage.Set(context.Background(), updateOrder)
				if err != nil {
					logger.Log().Sugar().Errorw("Can not update order", zap.Error(err))
				}
			}
		}
	}()
}

func New(log *zap.Logger, client accrual.Client, storage orders.Storage) Processor {
	return &processor{log, client, storage}
}
