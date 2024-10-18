package service

import (
	"juno/pkg/api/transaction"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type Service struct {
	logger          *logrus.Logger
	transactionRepo transaction.Repository
}

func New(logger *logrus.Logger, transactionRepo transaction.Repository) *Service {
	return &Service{
		logger:          logger,
		transactionRepo: transactionRepo,
	}
}

func (s *Service) CreateTransaction(userID uuid.UUID, amount int, key transaction.TransactionKey, meta map[string]string) error {
	return s.transactionRepo.CreateTransaction(
		transaction.NewTransaction(userID, amount, key, meta),
	)
}

func (s *Service) GetTransactionsByUserID(userID uuid.UUID) ([]*transaction.Transaction, error) {
	return s.transactionRepo.GetTransactionsByUserID(userID)
}

func (s *Service) Balance(userID uuid.UUID) (int, error) {
	transactions, err := s.GetTransactionsByUserID(userID)

	if err != nil {
		return 0, err
	}

	balance := 0

	for _, t := range transactions {
		balance += t.Amount
	}

	return balance, nil
}
