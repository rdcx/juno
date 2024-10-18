package service

import (
	"juno/pkg/api/transaction"
	tranRepo "juno/pkg/api/transaction/repo/mem"
	tranService "juno/pkg/api/transaction/service"
	"testing"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

func TestBalance(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		logger := logrus.New()
		tranRepo := tranRepo.New()
		tranService := tranService.New(logger, tranRepo)

		tokenService := New(tranService)

		userID := uuid.New()

		tranRepo.CreateTransaction(&transaction.Transaction{
			ID:     uuid.New(),
			UserID: userID,
			Amount: 100,
			Key:    transaction.DepositKey,
			Meta:   map[string]string{"provider": "paddle"},
		})

		balance, err := tokenService.Balance(userID)

		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}

		if balance != 100 {
			t.Errorf("Expected 100, got %d", balance)
		}
	})
}

func TestDeposit(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		logger := logrus.New()
		tranRepo := tranRepo.New()
		tranService := tranService.New(logger, tranRepo)

		tokenService := New(tranService)

		userID := uuid.New()

		tranRepo.CreateTransaction(&transaction.Transaction{
			ID:     uuid.New(),
			UserID: userID,
			Amount: 100,
			Key:    transaction.DepositKey,
			Meta:   map[string]string{"provider": "paddle"},
		})

		err := tokenService.Deposit(userID, 100)

		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}

		balance, err := tokenService.Balance(userID)

		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}

		if balance != 100 {
			t.Errorf("Expected 100, got %d", balance)
		}
	})

	t.Run("invalid amount", func(t *testing.T) {
		logger := logrus.New()
		tranRepo := tranRepo.New()
		tranService := tranService.New(logger, tranRepo)

		tokenService := New(tranService)

		userID := uuid.New()

		err := tokenService.Deposit(userID, 0)
		if err == nil {
			t.Errorf("Expected error, got nil")
		}
	})
}

func TestDebit(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		logger := logrus.New()
		tranRepo := tranRepo.New()
		tranService := tranService.New(logger, tranRepo)

		tokenService := New(tranService)

		userID := uuid.New()

		tranRepo.CreateTransaction(&transaction.Transaction{
			ID:     uuid.New(),
			UserID: userID,
			Amount: 100,
			Key:    transaction.DepositKey,
			Meta:   map[string]string{"provider": "paddle"},
		})

		err := tokenService.Debit(userID, transaction.QueryExecutionKey, 100, map[string]string{"provider": "paddle"})

		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}

		balance, err := tokenService.Balance(userID)

		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}

		if balance != 0 {
			t.Errorf("Expected -100, got %d", balance)
		}
	})

	t.Run("invalid amount", func(t *testing.T) {
		logger := logrus.New()
		tranRepo := tranRepo.New()
		tranService := tranService.New(logger, tranRepo)

		tokenService := New(tranService)

		userID := uuid.New()

		err := tokenService.Debit(userID, transaction.QueryExecutionKey, -100, map[string]string{"provider": "paddle"})
		if err == nil {
			t.Errorf("Expected error, got nil")
		}
	})

	t.Run("insufficient balance", func(t *testing.T) {
		logger := logrus.New()
		tranRepo := tranRepo.New()
		tranService := tranService.New(logger, tranRepo)

		tokenService := New(tranService)

		userID := uuid.New()

		err := tokenService.Debit(userID, transaction.QueryExecutionKey, 100, map[string]string{"provider": "paddle"})
		if err == nil {
			t.Errorf("Expected error, got nil")
		}
	})

	t.Run("invalid transaction key", func(t *testing.T) {
		logger := logrus.New()
		tranRepo := tranRepo.New()
		tranService := tranService.New(logger, tranRepo)

		tokenService := New(tranService)

		userID := uuid.New()

		err := tokenService.Debit(userID, transaction.DepositKey, 100, map[string]string{"provider": "paddle"})
		if err == nil {
			t.Errorf("Expected error, got nil")
		}
	})
}
