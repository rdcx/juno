package handler

import (
	"context"
	"encoding/json"
	"errors"
	"juno/pkg/api/auth"
	"juno/pkg/api/transaction"
	"juno/pkg/api/transaction/dto"
	"juno/pkg/api/user"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type mockTransactionService struct {
	withErr      error
	transactions []*transaction.Transaction
}

func (m *mockTransactionService) GetTransactionsByUserID(userID uuid.UUID) ([]*transaction.Transaction, error) {
	return m.transactions, m.withErr
}

func (m *mockTransactionService) CreateTransaction(userID uuid.UUID, amount int, key transaction.TransactionKey, meta map[string]string) error {
	return m.withErr
}

func (m *mockTransactionService) Balance(userID uuid.UUID) (int, error) {
	return 0, m.withErr
}

func TestList(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		service := &mockTransactionService{
			transactions: []*transaction.Transaction{
				{
					ID:     uuid.New(),
					UserID: uuid.New(),
					Amount: 100,
					Key:    "deposit",
					Meta:   map[string]string{"provider": "paddle"},
				},
			},
		}

		w := httptest.NewRecorder()
		tc, _ := gin.CreateTestContext(w)

		tc.Request = httptest.NewRequestWithContext(
			auth.WithUser(context.Background(), &user.User{ID: uuid.New()}),
			"GET",
			"/transactions",
			nil,
		)

		h := New(service)

		h.List(tc)

		if w.Code != http.StatusOK {
			t.Errorf("Expected 200, got %d", w.Code)
		}

		var res dto.ListResponse

		err := json.Unmarshal(w.Body.Bytes(), &res)

		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}

		if res.Status != dto.SUCCESS {
			t.Errorf("Expected %s, got %s", dto.SUCCESS, res.Status)
		}

		if len(res.Transactions) != 1 {
			t.Errorf("Expected 1, got %d", len(res.Transactions))
		}

		for i, tran := range res.Transactions {
			if tran.ID != service.transactions[i].ID.String() {
				t.Errorf("Expected %s, got %s", service.transactions[i].ID.String(), tran.ID)
			}

			if tran.UserID != service.transactions[i].UserID.String() {
				t.Errorf("Expected %s, got %s", service.transactions[i].UserID.String(), tran.UserID)
			}

			if tran.Amount != service.transactions[i].Amount {
				t.Errorf("Expected %d, got %d", service.transactions[i].Amount, tran.Amount)
			}

			if tran.Key != string(service.transactions[i].Key) {
				t.Errorf("Expected %s, got %s", string(service.transactions[i].Key), tran.Key)
			}

			if tran.Meta["provider"] != service.transactions[i].Meta["provider"] {
				t.Errorf("Expected %s, got %s", service.transactions[i].Meta["provider"], tran.Meta["provider"])
			}
		}
	})

	t.Run("failure", func(t *testing.T) {
		service := &mockTransactionService{
			withErr: errors.New("test error"),
		}

		w := httptest.NewRecorder()
		tc, _ := gin.CreateTestContext(w)

		tc.Request = httptest.NewRequestWithContext(
			auth.WithUser(context.Background(), &user.User{ID: uuid.New()}),
			"GET",
			"/transactions",
			nil,
		)

		h := New(service)

		h.List(tc)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("Expected 500, got %d", w.Code)
		}

		var res dto.ListResponse

		err := json.Unmarshal(w.Body.Bytes(), &res)

		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}

		if res.Status != dto.ERROR {
			t.Errorf("Expected %s, got %s", dto.ERROR, res.Status)
		}

		if res.Message != "test error" {
			t.Errorf("Expected test error, got %s", res.Message)
		}
	})
}
