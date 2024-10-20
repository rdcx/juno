package mysql

import (
	"database/sql"
	"juno/pkg/api/transaction"
	"juno/pkg/api/transaction/migration/mysql"
	"testing"

	_ "github.com/go-sql-driver/mysql"

	"github.com/google/uuid"
)

func newTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("mysql", "root:juno@tcp(localhost:3306)/transaction_test")

	if err != nil {
		t.Fatalf("could not connect to mysql: %v", err)
	}

	if err := mysql.ExecuteMigrations(db); err != nil {
		t.Fatalf("could not execute migrations: %v", err)
	}

	return db
}

func TestCreateTransaction(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		db := newTestDB(t)

		repo := New(db)

		tran := transaction.NewTransaction(
			uuid.New(),
			1,
			transaction.QueryExecutionKey,
			map[string]string{"query_id": uuid.New().String()},
		)

		if err := repo.CreateTransaction(tran); err != nil {
			t.Errorf("Expected nil, got %v", err)
		}

		transactions, err := repo.GetTransactionsByUserID(tran.UserID)

		if err != nil {
			t.Fatalf("Expected nil, got %v", err)
		}

		if len(transactions) != 1 {
			t.Errorf("Expected 1, got %d", len(transactions))
		}

		if transactions[0].ID != tran.ID {
			t.Errorf("Expected %s, got %s", tran.ID, transactions[0].ID)
		}
	})

	t.Run("failure", func(t *testing.T) {
		db := newTestDB(t)

		repo := New(db)

		tran := transaction.NewTransaction(
			uuid.New(),
			1,
			transaction.QueryExecutionKey,
			map[string]string{"query_id": uuid.New().String()},
		)

		if err := repo.CreateTransaction(tran); err != nil {
			t.Errorf("Expected nil, got %v", err)
		}

		if err := repo.CreateTransaction(tran); err == nil {
			t.Errorf("Expected error, got nil")
		}

	})
}

func TestGetTransactionsByUserID(t *testing.T) {
	db := newTestDB(t)

	repo := New(db)

	userID := uuid.New()
	tran := transaction.NewTransaction(
		userID,
		1,
		transaction.QueryExecutionKey,
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

	if transactions[0].Amount != tran.Amount {
		t.Errorf("Expected %f, got %f", tran.Amount, transactions[0].Amount)
	}

	if transactions[0].Key != tran.Key {
		t.Errorf("Expected %s, got %s", tran.Key, transactions[0].Key)
	}

	if transactions[0].Meta["query_id"] != tran.Meta["query_id"] {
		t.Errorf("Expected %s, got %s", tran.Meta["query_id"], transactions[0].Meta["query_id"])
	}
}
