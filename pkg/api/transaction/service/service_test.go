package service

import (
	"juno/pkg/api/transaction"
	"juno/pkg/api/transaction/repo/mem"
	"testing"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

func TestCreateTransaction(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := mem.New()

		service := New(logrus.New(), repo)

		tran := transaction.NewTransaction(
			uuid.New(),
			138492,
			transaction.QueryExecutionKey,
			map[string]string{"query_id": uuid.New().String()},
		)

		if err := service.CreateTransaction(
			tran.UserID,
			tran.Amount,
			tran.Key,
			tran.Meta,
		); err != nil {
			t.Errorf("Expected nil, got %v", err)
		}

		transactions, err := service.GetTransactionsByUserID(tran.UserID)

		if err != nil {
			t.Fatalf("Expected nil, got %v", err)
		}

		if len(transactions) != 1 {
			t.Errorf("Expected 1, got %d", len(transactions))
		}

		if transactions[0].Amount != tran.Amount {
			t.Errorf("Expected %d, got %d", 138492, transactions[0].Amount)
		}
	})
}

func TestBalance(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := mem.New()

		service := New(logrus.New(), repo)

		userID := uuid.New()

		tran := transaction.NewTransaction(
			userID,
			138492,
			transaction.DepositKey,
			map[string]string{"query_id": uuid.New().String()},
		)

		if err := service.CreateTransaction(
			tran.UserID,
			tran.Amount,
			tran.Key,
			tran.Meta,
		); err != nil {
			t.Errorf("Expected nil, got %v", err)
		}

		// debit
		tran = transaction.NewTransaction(
			userID,
			-1000,
			transaction.WithdrawalKey,
			map[string]string{"query_id": uuid.New().String()},
		)

		if err := service.CreateTransaction(
			tran.UserID,
			tran.Amount,
			tran.Key,
			tran.Meta,
		); err != nil {
			t.Errorf("Expected nil, got %v", err)
		}

		balance, err := service.Balance(userID)

		if err != nil {
			t.Fatalf("Expected nil, got %v", err)
		}

		if balance != 137492 {
			t.Errorf("Expected 137492, got %d", balance)
		}
	})
}
