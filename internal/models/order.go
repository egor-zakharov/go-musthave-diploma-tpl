package models

import "time"

type Order struct {
	Number     string    `json:"number"`
	UserID     string    `json:"-"`
	Status     string    `json:"status"`
	Accrual    float64   `json:"accrual"`
	UploadedAt time.Time `json:"uploaded_at"`
}
