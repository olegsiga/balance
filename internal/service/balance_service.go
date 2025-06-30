package service

import (
	"fmt"
	"strconv"

	"balance/internal/models"
	"balance/internal/repository"
	"github.com/shopspring/decimal"
)

type BalanceService struct {
	userRepo        *repository.UserRepository
	transactionRepo *repository.TransactionRepository
}

func NewBalanceService(userRepo *repository.UserRepository, transactionRepo *repository.TransactionRepository) *BalanceService {
	return &BalanceService{
		userRepo:        userRepo,
		transactionRepo: transactionRepo,
	}
}

func (s *BalanceService) GetUserBalance(userID int64) (*models.UserBalance, error) {
	user, err := s.userRepo.GetBalance(userID)
	if err != nil {
		return nil, err
	}

	return &models.UserBalance{
		UserID:  user.ID,
		Balance: user.Balance.StringFixed(2),
	}, nil
}

func (s *BalanceService) ProcessTransaction(userID int64, req *models.TransactionRequest, sourceType string) error {
	amount, err := decimal.NewFromString(req.Amount)
	if err != nil {
		return fmt.Errorf("invalid amount format: %w", err)
	}

	if amount.IsNegative() || amount.IsZero() {
		return fmt.Errorf("amount must be positive")
	}

	decimalPlaces := s.getDecimalPlaces(req.Amount)
	if decimalPlaces > 2 {
		return fmt.Errorf("amount can have at most 2 decimal places")
	}

	return s.transactionRepo.ProcessTransactionWithBalance(userID, req.TransactionID, sourceType, req.State, amount)
}

func (s *BalanceService) getDecimalPlaces(amount string) int {
	decimal, err := decimal.NewFromString(amount)
	if err != nil {
		return 0
	}
	return int(-decimal.Exponent())
}

func (s *BalanceService) ValidateUserID(userIDStr string) (int64, error) {
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid user ID format")
	}
	if userID <= 0 {
		return 0, fmt.Errorf("user ID must be positive")
	}
	return userID, nil
}
