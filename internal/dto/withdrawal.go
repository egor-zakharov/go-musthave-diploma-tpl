package dto

import "time"

type WithdrawalRequest struct {
	Number string  `json:"order"`
	Sum    float64 `json:"sum"`
}

type WithdrawalsResponse struct {
	OrderNumber string    `json:"order"`
	Sum         float64   `json:"sum"`
	ProcessedAt time.Time `json:"processed_at"`
}
