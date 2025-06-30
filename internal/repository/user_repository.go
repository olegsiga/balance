package repository

import (
	"database/sql"
	"fmt"

	"balance/internal/models"
	"github.com/shopspring/decimal"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetBalance(userID int64) (*models.User, error) {
	var user models.User
	query := `SELECT id, balance, created_at, updated_at FROM users WHERE id = $1`

	err := r.db.QueryRow(query, userID).Scan(
		&user.ID,
		&user.Balance,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user with id %d not found", userID)
		}
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) UpdateBalance(userID int64, newBalance decimal.Decimal) error {
	query := `UPDATE users SET balance = $1, updated_at = NOW() WHERE id = $2`
	result, err := r.db.Exec(query, newBalance, userID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user with id %d not found", userID)
	}

	return nil
}
