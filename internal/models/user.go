package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type User struct {
	ID        int64           `json:"id"`
	Balance   decimal.Decimal `json:"balance"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

type UserBalance struct {
	UserID  int64  `json:"userId"`
	Balance string `json:"balance"`
}
