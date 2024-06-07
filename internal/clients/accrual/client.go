package accrual

import (
	"encoding/json"
	"fmt"
	"github.com/egor-zakharov/go-musthave-diploma-tpl/internal/dto"
	"github.com/egor-zakharov/go-musthave-diploma-tpl/internal/logger"
	"go.uber.org/zap"
	"io"
	"net/http"
)

type accrual struct {
	log  *zap.Logger
	addr string
}

func New(log *zap.Logger, addr string) Client {
	return &accrual{log, addr}
}

func (a accrual) SendOrder(orderID string) (*dto.AccrualOrderResponse, error) {
	order := &dto.AccrualOrderResponse{}
	requestURL := fmt.Sprintf("%s/api/orders/%s", a.addr, orderID)
	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		logger.Log().Sugar().Errorw("New request error", zap.Error(err))
		return order, err
	}

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.Log().Sugar().Errorw("Do request error", zap.Error(err))
		return order, err
	}

	if response.StatusCode == http.StatusOK {
		resBody, err := io.ReadAll(response.Body)
		defer response.Body.Close()
		if err != nil {
			logger.Log().Sugar().Errorw("Read body error", zap.Error(err))
			return order, err
		}
		err = json.Unmarshal(resBody, &order)
		if err != nil {
			logger.Log().Sugar().Errorw("Unmarshal body error", zap.Error(err))
			return order, err
		}
		return order, nil
	}

	if response.StatusCode == http.StatusInternalServerError {
		logger.Log().Sugar().Errorw("Internal Server Error", zap.Error(ErrAccrualServerError))
		return order, ErrAccrualServerError
	}

	if response.StatusCode == http.StatusTooManyRequests {
		logger.Log().Sugar().Errorw("Too Many Requests", zap.Error(ErrAccrualTooManyRequests))
		return order, ErrAccrualTooManyRequests
	}

	if response.StatusCode == http.StatusNoContent {
		logger.Log().Sugar().Errorw("No contend", zap.Error(ErrAccrualNoData))
		return order, ErrAccrualNoData
	}

	return order, nil
}
