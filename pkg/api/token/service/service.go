package service

import (
	"juno/pkg/api/token"
	"juno/pkg/api/transaction"
	"time"

	"github.com/google/uuid"
)

type Service struct {
	transactionService transaction.Service
}

func New(transactionService transaction.Service) *Service {
	return &Service{
		transactionService: transactionService,
	}
}

func (s *Service) Balance(userID uuid.UUID) (int, error) {
	return s.transactionService.Balance(userID)
}

func (s *Service) Deposit(userID uuid.UUID, amount int) error {

	if amount < 0 {
		return token.ErrInvalidAmount
	}

	// simulate a delay
	time.Sleep(1 * time.Second)

	return s.transactionService.CreateTransaction(
		userID,
		amount,
		transaction.DepositKey,
		map[string]string{"provider": "paddle"},
	)
}

func (s *Service) Debit(userID uuid.UUID, tranKey transaction.TransactionKey, amount int, meta map[string]string) error {

	if amount < 0 {
		return token.ErrInvalidAmount
	}

	allowedKeys := []transaction.TransactionKey{
		transaction.QueryExecutionKey,
	}

	if !transaction.ContainsKey(allowedKeys, tranKey) {
		return token.ErrInvalidAmount
	}

	balance, err := s.transactionService.Balance(userID)

	if err != nil {
		return err
	}

	if balance < amount {
		return token.ErrInsufficientBalance
	}

	// simulate a delay
	time.Sleep(1 * time.Second)

	return s.transactionService.CreateTransaction(userID, -amount, transaction.WithdrawalKey, meta)
}
