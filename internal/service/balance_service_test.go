package service

import (
	"testing"
)

func TestGetDecimalPlaces(t *testing.T) {
	s := &BalanceService{}

	tests := []struct {
		amount   string
		expected int
	}{
		{"10.00", 2},
		{"10.5", 1},
		{"10", 0},
		{"10.123", 3},
	}

	for _, test := range tests {
		result := s.getDecimalPlaces(test.amount)
		if result != test.expected {
			t.Errorf("getDecimalPlaces(%s) = %d; expected %d", test.amount, result, test.expected)
		}
	}
}

func TestValidateUserID(t *testing.T) {
	s := &BalanceService{}

	tests := []struct {
		userIDStr string
		expected  int64
		hasError  bool
	}{
		{"1", 1, false},
		{"123", 123, false},
		{"0", 0, true},
		{"-1", 0, true},
		{"abc", 0, true},
		{"", 0, true},
	}

	for _, test := range tests {
		result, err := s.ValidateUserID(test.userIDStr)
		if test.hasError {
			if err == nil {
				t.Errorf("ValidateUserID(%s) should have returned an error", test.userIDStr)
			}
		} else {
			if err != nil {
				t.Errorf("ValidateUserID(%s) returned unexpected error: %v", test.userIDStr, err)
			}
			if result != test.expected {
				t.Errorf("ValidateUserID(%s) = %d; expected %d", test.userIDStr, result, test.expected)
			}
		}
	}
}
