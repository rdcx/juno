package mysql

import (
	"database/sql"
	"encoding/json"
	"juno/pkg/api/transaction"

	"github.com/google/uuid"
)

type Repository struct {
	db *sql.DB
}

func New(db *sql.DB) *Repository {
	return &Repository{db}
}

func (r *Repository) CreateTransaction(t *transaction.Transaction) error {

	meta, err := json.Marshal(t.Meta)

	if err != nil {
		return err
	}

	_, err = r.db.Exec(`
		INSERT INTO transactions (id, user_id, type, amount, meta)
		VALUES (?, ?, ?, ?, ?)
	`, t.ID, t.UserID, t.Key, t.Amount, meta)
	return err
}

func (r *Repository) GetTransactionsByUserID(userID uuid.UUID) ([]*transaction.Transaction, error) {
	rows, err := r.db.Query(`
		SELECT id, user_id, type, amount, meta
		FROM transactions
		WHERE user_id = ?
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []*transaction.Transaction

	var meta string
	for rows.Next() {
		var t transaction.Transaction
		if err := rows.Scan(&t.ID, &t.UserID, &t.Key, &t.Amount, &meta); err != nil {
			return nil, err
		}

		err := json.Unmarshal([]byte(meta), &t.Meta)

		if err != nil {
			return nil, err
		}

		transactions = append(transactions, &t)
	}

	return transactions, nil
}
