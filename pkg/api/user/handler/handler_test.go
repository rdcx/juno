package handler

import (
	"context"
	"encoding/json"
	"juno/pkg/api/auth"
	"juno/pkg/api/user"
	"juno/pkg/api/user/dto"
	"juno/pkg/api/user/policy"
	"juno/pkg/api/user/repo/mem"
	"juno/pkg/api/user/service"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

func randomEmail() string {
	return uuid.New().String() + "@example.com"
}

func testGetResponse(t *testing.T, w *httptest.ResponseRecorder, expectedCode int, expectedStatus, expectedMessage string, expectedUser *user.User) {
	if w.Code != expectedCode {
		t.Fatalf("Expected status code %d, got %d", expectedCode, w.Code)
	}

	var resp dto.GetUserResponse

	err := json.Unmarshal(w.Body.Bytes(), &resp)

	if err != nil {
		t.Fatalf("Error parsing response: %s", err)
	}

	if resp.Status != expectedStatus {
		t.Fatalf("Expected status %s, got %s", expectedStatus, resp.Status)
	}

	if resp.Message != expectedMessage {
		t.Fatalf("Expected message %s, got %s", expectedMessage, resp.Message)
	}

	if expectedUser == nil && resp.User == nil {
		return
	}

	if resp.User.ID != expectedUser.ID.String() {
		t.Fatalf("Expected user ID %s, got %s", expectedUser.ID, resp.User.ID)
	}

	if resp.User.Email != expectedUser.Email {
		t.Fatalf("Expected user email %s, got %s", expectedUser.Email, resp.User.Email)
	}
}

func TestGet(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		logger := logrus.New()
		repo := mem.New()
		svc := service.New(logger, repo)
		policy := policy.New()
		handler := New(logger, policy, svc)

		u := &user.User{
			ID:    uuid.New(),
			Email: randomEmail(),
		}

		err := repo.Create(u)

		if err != nil {
			t.Errorf("expected err to be nil, got %v", err)
		}

		w := httptest.NewRecorder()

		tc, _ := gin.CreateTestContext(w)

		tc.Params = gin.Params{
			{Key: "id", Value: u.ID.String()},
		}

		req := httptest.NewRequest("GET", "/users/"+u.ID.String(), nil)
		tc.Request = req.WithContext(
			auth.WithUser(context.Background(), u),
		)

		handler.Get(tc)

		testGetResponse(t, w, 200, dto.SUCCESS, "", u)
	})

	t.Run("not found", func(t *testing.T) {
		logger := logrus.New()
		repo := mem.New()
		svc := service.New(logger, repo)
		policy := policy.New()
		handler := New(logger, policy, svc)

		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/users/"+uuid.New().String(), nil)
		tc, _ := gin.CreateTestContext(w)
		tc.Request = req.WithContext(auth.WithUser(context.Background(), &user.User{
			ID: uuid.New(),
		}))

		tc.Params = gin.Params{
			{Key: "id", Value: uuid.New().String()},
		}

		handler.Get(tc)

		testGetResponse(t, w, 404, dto.ERROR, user.ErrNotFound.Error(), nil)
	})

	t.Run("unauthorized", func(t *testing.T) {
		logger := logrus.New()
		repo := mem.New()
		svc := service.New(logger, repo)
		policy := policy.New()
		handler := New(logger, policy, svc)

		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/users/"+uuid.New().String(), nil)
		tc, _ := gin.CreateTestContext(w)
		tc.Request = req.WithContext(auth.WithUser(context.Background(), &user.User{
			ID: uuid.New(),
		}))

		usr := &user.User{
			ID: uuid.New(),
		}

		err := repo.Create(usr)

		if err != nil {
			t.Errorf("expected err to be nil, got %v", err)
		}

		tc.Params = gin.Params{
			{Key: "id", Value: usr.ID.String()},
		}

		handler.Get(tc)

		testGetResponse(t, w, 404, dto.ERROR, user.ErrNotFound.Error(), nil)
	})

	t.Run("invalid id", func(t *testing.T) {
		logger := logrus.New()
		repo := mem.New()
		svc := service.New(logger, repo)
		policy := policy.New()
		handler := New(logger, policy, svc)

		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/users/"+uuid.New().String(), nil)
		tc, _ := gin.CreateTestContext(w)
		tc.Request = req.WithContext(auth.WithUser(context.Background(), &user.User{
			ID: uuid.New(),
		}))

		tc.Params = gin.Params{
			{Key: "id", Value: "invalid"},
		}

		handler.Get(tc)

		testGetResponse(t, w, 400, dto.ERROR, user.ErrInvalidID.Error(), nil)
	})

	t.Run("user context not set", func(t *testing.T) {
		logger := logrus.New()
		repo := mem.New()
		svc := service.New(logger, repo)
		policy := policy.New()
		handler := New(logger, policy, svc)

		u := &user.User{
			ID: uuid.New(),
		}
		err := repo.Create(u)

		if err != nil {
			t.Errorf("expected err to be nil, got %v", err)
		}

		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/users/"+uuid.New().String(), nil)
		tc, _ := gin.CreateTestContext(w)
		tc.Request = req

		tc.Params = gin.Params{
			{Key: "id", Value: u.ID.String()},
		}

		handler.Get(tc)

		testGetResponse(t, w, 500, dto.ERROR, user.ErrInternal.Error(), nil)
	})
}
