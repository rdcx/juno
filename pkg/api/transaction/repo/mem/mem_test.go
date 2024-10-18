package mem

import (
	"juno/pkg/api/transaction"
	"testing"

	"github.com/google/uuid"
)

func TestCreateTransaction(t *testing.T) {
	repo := New()

	tran := transaction.NewTransaction(
		uuid.New(),
		1,
		transaction.QueryExecutionKey,
		map[string]string{"query_id": uuid.New().String()},
	)

	if err := repo.CreateTransaction(tran); err != nil {
		t.Errorf("Expected nil, got %v", err)
	}

	if len(repo.tokens) != 1 {
		t.Errorf("Expected 1, got %d", len(repo.tokens))
	}

	if repo.tokens[0].ID != tran.ID {
		t.Errorf("Expected %s, got %s", tran.ID, repo.tokens[0].ID)
	}
}

func TestGetTransactionsByUserID(t *testing.T) {
	repo := New()

	userID := uuid.New()
	tran := token.NewTransaction(
		userID,
		1,
		token.TransactionQueryExecution,
		map[string]string{"query_id": uuid.New().String()},
	)

	if err := repo.CreateTransaction(tran); err != nil {
		t.Errorf("Expected nil, got %v", err)
	}

	transactions, err := repo.GetTransactionsByUserID(userID)
	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}

	if len(transactions) != 1 {
		t.Errorf("Expected 1, got %d", len(transactions))
	}

	if transactions[0].ID != tran.ID {
		t.Errorf("Expected %s, got %s", tran.ID, transactions[0].ID)
	}

	if transactions[0].UserID != tran.UserID {
		t.Errorf("Expected %s, got %s", tran.UserID, transactions[0].UserID)
	}
}
