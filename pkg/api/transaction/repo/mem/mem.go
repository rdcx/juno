package mem

import (
	"juno/pkg/api/transaction"

	"github.com/google/uuid"
)

type Repository struct {
	tokens []*transaction.Transaction
}

func New() *Repository {
	return &Repository{}
}

func (r *Repository) CreateTransaction(t *transaction.Transaction) error {
	r.tokens = append(r.tokens, t)
	return nil
}

func (r *Repository) GetTransactionsByUserID(userID uuid.UUID) ([]*transaction.Transaction, error) {
	var transactions []*transaction.Transaction
	for _, t := range r.tokens {
		if t.UserID == userID {
			transactions = append(transactions, t)
		}
	}
	return transactions, nil
}
