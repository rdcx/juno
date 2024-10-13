package handler

import (
	"context"
	"encoding/json"
	"juno/pkg/api/auth"
	"juno/pkg/api/user"
	"juno/pkg/api/user/dto"
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

func TestGet(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		logger := logrus.New()
		repo := mem.New()
		svc := service.New(logger, repo)
		handler := New(logger, svc)

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

		if w.Code != 200 {
			t.Fatalf("Expected status code 200, got %d", w.Code)
		}

		var resp dto.GetUserResponse

		err = json.Unmarshal(w.Body.Bytes(), &resp)

		if err != nil {
			t.Fatalf("Error parsing response: %s", err)
		}

		if resp.Status != "success" {
			t.Fatalf("Expected status success, got %s", resp.Status)
		}

		if resp.Message != "" {
			t.Fatalf("Expected message to be empty, got %s", resp.Message)
		}

		if resp.User.ID != u.ID.String() {
			t.Fatalf("Expected user ID %s, got %s", u.ID, resp.User.ID)
		}

		if resp.User.Email != u.Email {
			t.Fatalf("Expected user email %s, got %s", u.Email, resp.User.Email)
		}
	})

	t.Run("not found", func(t *testing.T) {
		logger := logrus.New()
		repo := mem.New()
		svc := service.New(logger, repo)
		handler := New(logger, svc)

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

		if w.Code != 404 {
			t.Fatalf("Expected status code 404, got %d", w.Code)
		}

		var resp dto.GetUserResponse

		err := json.Unmarshal(w.Body.Bytes(), &resp)

		if err != nil {
			t.Fatalf("Error parsing response: %s", err)
		}

		if resp.Status != "error" {
			t.Fatalf("Expected status error, got %s", resp.Status)
		}

		if resp.Message != user.ErrNotFound.Error() {
			t.Fatalf("Expected message %s, got %s", user.ErrNotFound.Error(), resp.Message)
		}

		if resp.User != nil {
			t.Fatalf("Expected user to be nil, got %v", resp.User)
		}
	})

	t.Run("unauthorized", func(t *testing.T) {
		logger := logrus.New()
		repo := mem.New()
		svc := service.New(logger, repo)
		handler := New(logger, svc)

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

		if w.Code != 404 {
			t.Fatalf("Expected status code 404, got %d", w.Code)
		}

		var resp dto.GetUserResponse

		err = json.Unmarshal(w.Body.Bytes(), &resp)

		if err != nil {
			t.Fatalf("Error parsing response: %s", err)
		}

		if resp.Status != "error" {
			t.Fatalf("Expected status error, got %s", resp.Status)
		}

		if resp.Message != user.ErrNotFound.Error() {
			t.Fatalf("Expected message %s, got %s", user.ErrNotFound.Error(), resp.Message)
		}
	})

	t.Run("invalid id", func(t *testing.T) {
		logger := logrus.New()
		repo := mem.New()
		svc := service.New(logger, repo)
		handler := New(logger, svc)

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

		if w.Code != 400 {
			t.Fatalf("Expected status code 400, got %d", w.Code)
		}

		var resp dto.GetUserResponse

		err := json.Unmarshal(w.Body.Bytes(), &resp)

		if err != nil {
			t.Fatalf("Error parsing response: %s", err)
		}

		if resp.Status != "error" {
			t.Fatalf("Expected status error, got %s", resp.Status)
		}

		if resp.Message != user.ErrInvalidID.Error() {
			t.Fatalf("Expected message %s, got %s", user.ErrInvalidID.Error(), resp.Message)
		}
	})

	t.Run("user context not set", func(t *testing.T) {
		logger := logrus.New()
		repo := mem.New()
		svc := service.New(logger, repo)
		handler := New(logger, svc)

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

		if w.Code != 500 {
			t.Fatalf("Expected status code 500, got %d", w.Code)
		}

		var resp dto.GetUserResponse

		err = json.Unmarshal(w.Body.Bytes(), &resp)

		if err != nil {
			t.Fatalf("Error parsing response: %s", err)
		}

		if resp.Status != "error" {
			t.Fatalf("Expected status error, got %s", resp.Status)
		}

		if resp.Message != user.ErrInternal.Error() {
			t.Fatalf("Expected message %s, got %s", user.ErrInternal.Error(), resp.Message)
		}
	})
}
