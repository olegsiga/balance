package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type Transaction struct {
	ID            int64           `json:"id"`
	UserID        int64           `json:"user_id"`
	TransactionID string          `json:"transaction_id"`
	SourceType    string          `json:"source_type"`
	State         string          `json:"state"`
	Amount        decimal.Decimal `json:"amount"`
	CreatedAt     time.Time       `json:"created_at"`
}

type TransactionRequest struct {
	State         string `json:"state" binding:"required,oneof=win lose"`
	Amount        string `json:"amount" binding:"required"`
	TransactionID string `json:"transactionId" binding:"required"`
}

type TransactionResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}
