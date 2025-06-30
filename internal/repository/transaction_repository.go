package repository

import (
	"database/sql"
	"fmt"

	"github.com/shopspring/decimal"
)

type TransactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) TransactionExists(transactionID string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM transactions WHERE transaction_id = $1)`
	err := r.db.QueryRow(query, transactionID).Scan(&exists)
	return exists, err
}

func (r *TransactionRepository) CreateTransaction(userID int64, transactionID, sourceType, state string, amount decimal.Decimal) error {
	query := `
		INSERT INTO transactions (user_id, transaction_id, source_type, state, amount, created_at)
		VALUES ($1, $2, $3, $4, $5, NOW())
	`
	_, err := r.db.Exec(query, userID, transactionID, sourceType, state, amount)
	if err != nil {
		return fmt.Errorf("failed to create transaction: %w", err)
	}
	return nil
}

func (r *TransactionRepository) ProcessTransactionWithBalance(userID int64, transactionID, sourceType, state string, amount decimal.Decimal) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	exists, err := r.transactionExistsInTx(tx, transactionID)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("transaction %s already processed", transactionID)
	}

	var currentBalance decimal.Decimal
	query := `SELECT balance FROM users WHERE id = $1 FOR UPDATE`
	err = tx.QueryRow(query, userID).Scan(&currentBalance)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("user with id %d not found", userID)
		}
		return fmt.Errorf("failed to get user balance: %w", err)
	}

	var newBalance decimal.Decimal
	if state == "win" {
		newBalance = currentBalance.Add(amount)
	} else {
		newBalance = currentBalance.Sub(amount)
		if newBalance.IsNegative() {
			return fmt.Errorf("insufficient balance: current=%s, required=%s", currentBalance.String(), amount.String())
		}
	}

	updateQuery := `UPDATE users SET balance = $1, updated_at = NOW() WHERE id = $2`
	_, err = tx.Exec(updateQuery, newBalance, userID)
	if err != nil {
		return fmt.Errorf("failed to update balance: %w", err)
	}

	insertQuery := `
		INSERT INTO transactions (user_id, transaction_id, source_type, state, amount, created_at)
		VALUES ($1, $2, $3, $4, $5, NOW())
	`
	_, err = tx.Exec(insertQuery, userID, transactionID, sourceType, state, amount)
	if err != nil {
		return fmt.Errorf("failed to create transaction: %w", err)
	}

	return tx.Commit()
}

func (r *TransactionRepository) transactionExistsInTx(tx *sql.Tx, transactionID string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM transactions WHERE transaction_id = $1)`
	err := tx.QueryRow(query, transactionID).Scan(&exists)
	return exists, err
}
