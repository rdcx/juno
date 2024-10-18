package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"juno/pkg/api/auth"
	"juno/pkg/api/token/dto"
	"juno/pkg/api/user"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type mockTokenService struct {
	balance int
	withErr error

	depositAmount int
}

func (m *mockTokenService) Balance(userID uuid.UUID) (int, error) {
	return m.balance, m.withErr
}

func (m *mockTokenService) Deposit(userID uuid.UUID, amount int) error {
	if m.withErr != nil {
		return m.withErr
	}
	m.depositAmount = amount
	return nil
}

func TestBalance(t *testing.T) {
	t.Run("success", func(t *testing.T) {

		handler := New(logrus.New(), &mockTokenService{balance: 100, withErr: nil})

		userID := uuid.New()

		w := httptest.NewRecorder()
		tc, _ := gin.CreateTestContext(w)

		tc.Request = httptest.NewRequestWithContext(
			auth.WithUser(context.Background(), &user.User{ID: userID}),
			"GET",
			"/tokens/balance",
			nil,
		)

		handler.Balance(tc)

		var res dto.BalanceResponse

		err := json.Unmarshal(w.Body.Bytes(), &res)

		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}

		if res.Status != dto.SUCCESS {
			t.Errorf("Expected %s, got %s", dto.SUCCESS, res.Status)
		}

		if res.Balance != 100 {
			t.Errorf("Expected 100, got %d", res.Balance)
		}
	})

	t.Run("failure", func(t *testing.T) {

		errTest := errors.New("test error")

		handler := New(logrus.New(), &mockTokenService{balance: 0, withErr: errTest})

		userID := uuid.New()

		w := httptest.NewRecorder()
		tc, _ := gin.CreateTestContext(w)

		tc.Request = httptest.NewRequestWithContext(
			auth.WithUser(context.Background(), &user.User{ID: userID}),
			"GET",
			"/tokens/balance",
			nil,
		)

		handler.Balance(tc)

		var res dto.BalanceResponse

		err := json.Unmarshal(w.Body.Bytes(), &res)

		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}

		if res.Status != dto.ERROR {
			t.Errorf("Expected %s, got %s", dto.ERROR, res.Status)
		}

		if res.Message != errTest.Error() {
			t.Errorf("Expected %s, got %s", errTest.Error(), res.Message)
		}
	})
}

func TestDeposit(t *testing.T) {
	t.Run("success", func(t *testing.T) {

		mockTokenService := &mockTokenService{balance: 100, withErr: nil}
		handler := New(logrus.New(), mockTokenService)

		userID := uuid.New()

		w := httptest.NewRecorder()
		tc, _ := gin.CreateTestContext(w)

		req := dto.DepositRequest{Amount: 100}

		body, err := json.Marshal(req)

		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}

		tc.Request = httptest.NewRequestWithContext(
			auth.WithUser(context.Background(), &user.User{ID: userID}),
			"POST",
			"/tokens/deposit",
			bytes.NewBuffer(body),
		)

		handler.Deposit(tc)

		var res dto.DepositResponse

		err = json.Unmarshal(w.Body.Bytes(), &res)

		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}

		if res.Status != dto.SUCCESS {
			t.Errorf("Expected %s, got %s", dto.SUCCESS, res.Status)
		}

		if mockTokenService.depositAmount != 100 {
			t.Errorf("Expected 100, got %d", mockTokenService.depositAmount)
		}
	})

	t.Run("failure", func(t *testing.T) {

		errTest := errors.New("test error")

		mockTokenService := &mockTokenService{balance: 100, withErr: errTest}
		handler := New(logrus.New(), mockTokenService)

		userID := uuid.New()

		w := httptest.NewRecorder()
		tc, _ := gin.CreateTestContext(w)

		req := dto.DepositRequest{Amount: 100}

		body, err := json.Marshal(req)

		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}

		tc.Request = httptest.NewRequestWithContext(
			auth.WithUser(context.Background(), &user.User{ID: userID}),
			"POST",
			"/tokens/deposit",
			bytes.NewBuffer(body),
		)

		handler.Deposit(tc)

		var res dto.DepositResponse

		err = json.Unmarshal(w.Body.Bytes(), &res)

		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}

		if res.Status != dto.ERROR {
			t.Errorf("Expected %s, got %s", dto.ERROR, res.Status)
		}

		if res.Message != errTest.Error() {
			t.Errorf("Expected %s, got %s", errTest.Error(), res.Message)
		}

		if mockTokenService.depositAmount != 0 {
			t.Errorf("Expected 0, got %d", mockTokenService.depositAmount)
		}
	})

}
