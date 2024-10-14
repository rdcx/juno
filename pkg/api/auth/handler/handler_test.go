package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"juno/pkg/api/auth/dto"
	authService "juno/pkg/api/auth/service"
	"juno/pkg/api/user"
	usrRepo "juno/pkg/api/user/repo/mem"
	usrService "juno/pkg/api/user/service"
	"juno/pkg/util"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

func randomEmail() string {
	return fmt.Sprintf("%s@%s.com", uuid.New().String(), uuid.New().String())
}

func TestToken(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		os.Setenv("SECRET", "secret")

		userID := uuid.New()

		hash, _ := util.BcryptPassword("password")

		u := &user.User{
			ID:       userID,
			Email:    randomEmail(),
			Password: hash,
		}

		usrRepo := usrRepo.New()
		usrSvc := usrService.New(logrus.New(), usrRepo)

		err := usrRepo.Create(u)

		if err != nil {
			t.Errorf("expected err to be nil, got %v", err)
		}

		authSvc := authService.New(logrus.New(), usrSvc)
		authHandler := New(logrus.New(), authSvc)

		var req dto.TokenRequest

		req.Email = u.Email
		req.Password = "password"

		encoded, err := json.Marshal(req)

		if err != nil {
			t.Fatalf("Error encoding request: %s", err)
		}

		w := httptest.NewRecorder()
		tc, _ := gin.CreateTestContext(w)
		tc.Request = httptest.NewRequest(http.MethodPost, "/auth/token", bytes.NewBuffer(encoded))

		authHandler.Token(tc)

		if w.Code != 200 {
			t.Errorf("expected status code 200, got %v", w.Code)
		}

		// Check if the response body contains the user ID
		var resp dto.TokenResponse
		err = json.Unmarshal(w.Body.Bytes(), &resp)

		if err != nil {
			t.Fatalf("Error parsing response: %s", err)
		}

		if resp.Token == "" {
			t.Fatalf("Expected token, got empty string")
		}

		// Validate the token using TokenToUser
		parsedUser, err := authService.TokenToUser(resp.Token)

		if err != nil {
			t.Fatalf("Expected no error when parsing token, got %v", err)
		}

		if parsedUser.ID != u.ID {
			t.Errorf("Expected user ID %v, got %v", u.ID, parsedUser.ID)
		}

		if parsedUser.Email != u.Email {
			t.Errorf("Expected email %v, got %v", u.Email, parsedUser.Email)
		}

	})

	t.Run("invalid password", func(t *testing.T) {
		os.Setenv("SECRET", "secret")

		userID := uuid.New()

		hash, _ := util.BcryptPassword("password")

		u := &user.User{
			ID:       userID,
			Email:    randomEmail(),
			Password: hash,
		}

		usrRepo := usrRepo.New()
		usrSvc := usrService.New(logrus.New(), usrRepo)

		err := usrRepo.Create(u)

		if err != nil {
			t.Errorf("expected err to be nil, got %v", err)
		}

		authSvc := authService.New(logrus.New(), usrSvc)
		authHandler := New(logrus.New(), authSvc)

		var req dto.TokenRequest

		req.Email = u.Email
		req.Password = "invalid"

		encoded, err := json.Marshal(req)

		if err != nil {
			t.Fatalf("Error encoding request: %s", err)
		}

		w := httptest.NewRecorder()

		tc, _ := gin.CreateTestContext(w)

		tc.Request = httptest.NewRequest(http.MethodPost, "/auth/token", bytes.NewBuffer(encoded))

		authHandler.Token(tc)

		if w.Code != 400 {
			t.Errorf("expected status code 400, got %v", w.Code)
		}

		// Check if the response body contains the user ID
		var resp dto.TokenResponse
		err = json.Unmarshal(w.Body.Bytes(), &resp)

		if err != nil {
			t.Fatalf("Error parsing response: %s", err)
		}

		if resp.Message == "" {
			t.Fatalf("Expected error, got empty string")
		}

		if resp.Message != "invalid email or password" {
			t.Fatalf("Expected error 'invalid email or password', got %s", resp.Message)
		}

		if resp.Token != "" {
			t.Fatalf("Expected empty token, got %s", resp.Token)
		}

	})
}
