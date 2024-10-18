package transaction

import (
	"github.com/google/uuid"
)

type TransactionKey string

const (
	DepositKey        TransactionKey = "deposit"
	QueryExecutionKey TransactionKey = "query_execution"
	WithdrawalKey     TransactionKey = "withdrawal"
	NodeEarningsKey   TransactionKey = "node_earnings"
)

type Transaction struct {
	ID     uuid.UUID
	UserID uuid.UUID
	Amount int
	Key    TransactionKey
	Meta   map[string]string
}

func NewTransaction(userID uuid.UUID, amount int, key TransactionKey, meta map[string]string) *Transaction {
	return &Transaction{
		ID:     uuid.New(),
		UserID: userID,
		Amount: amount,
		Key:    key,
		Meta:   meta,
	}
}

type Repository interface {
	CreateTransaction(t *Transaction) error
	GetTransactionsByUserID(userID uuid.UUID) ([]*Transaction, error)
}

type Service interface {
	CreateTransaction(
		userID uuid.UUID,
		amount int,
		key TransactionKey,
		meta map[string]string,
	) error
	Balance(userID uuid.UUID) (int, error)
	GetTransactionsByUserID(userID uuid.UUID) ([]*Transaction, error)
}
