package service

import (
	"errors"
	"juno/pkg/api/auth"
	"juno/pkg/api/user"
	"juno/pkg/util"

	usrRepo "juno/pkg/api/user/repo/mem"
	usrService "juno/pkg/api/user/service"

	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// Helper function to generate a JWT token for testing
func generateTestJWT(userID uuid.UUID, email string, secret string, exp time.Time) string {
	claims := jwt.MapClaims{
		"id":    userID.String(),
		"email": email,
		"exp":   exp.Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte(secret))
	return tokenString
}

// Test TokenToUser
func TestTokenToUser(t *testing.T) {
	os.Setenv("SECRET", "mysecretkey")
	secret := os.Getenv("SECRET")

	userID := uuid.New()
	email := "test@example.com"

	t.Run("Valid token", func(t *testing.T) {
		exp := time.Now().Add(time.Hour)
		token := generateTestJWT(userID, email, secret, exp)

		u, err := TokenToUser(token)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if u.ID != userID {
			t.Errorf("Expected user ID %v, got %v", userID, u.ID)
		}
		if u.Email != email {
			t.Errorf("Expected email %v, got %v", email, u.Email)
		}
	})

	t.Run("Expired token", func(t *testing.T) {
		exp := time.Now().Add(-time.Hour)
		token := generateTestJWT(userID, email, secret, exp)

		_, err := TokenToUser(token)
		if !errors.Is(err, auth.ErrExpiredToken) {
			t.Errorf("Expected error %v, got %v", auth.ErrExpiredToken, err)
		}
	})

	t.Run("Invalid token", func(t *testing.T) {
		invalidToken := "invalid.token.here"

		_, err := TokenToUser(invalidToken)
		if err == nil {
			t.Errorf("Expected error, got nil")
		}
	})
}

// Test Token function
func TestToken(t *testing.T) {
	os.Setenv("SECRET", "mysecretkey")

	userID := uuid.New()
	email := "test@example.com"
	u := &user.User{
		ID:    userID,
		Email: email,
	}

	t.Run("Generate valid token", func(t *testing.T) {
		token, err := Token(u)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		// Validate the token using TokenToUser
		parsedUser, err := TokenToUser(token)
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
}

// Test Authenticate method
func TestAuthenticate(t *testing.T) {

	pass, _ := util.BcryptPassword("validpassword")

	usrRepo := usrRepo.New()
	usrRepo.Create(&user.User{
		ID:       uuid.New(),
		Email:    "test@example.com",
		Password: pass,
	})

	usrService := usrService.New(logrus.New(), usrRepo)

	logger := logrus.New()
	service := New(logger, usrService)

	t.Run("Valid credentials", func(t *testing.T) {
		token, err := service.Authenticate("test@example.com", "validpassword")
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		// Verify the token is valid
		_, err = TokenToUser(token)
		if err != nil {
			t.Errorf("Expected valid token, got error: %v", err)
		}
	})

	t.Run("Invalid email", func(t *testing.T) {
		_, err := service.Authenticate("invalid@example.com", "validpassword")
		if !errors.Is(err, user.ErrNotFound) {
			t.Errorf("Expected error %v, got %v", user.ErrNotFound, err)
		}
	})

	t.Run("Invalid password", func(t *testing.T) {
		_, err := service.Authenticate("test@example.com", "invalidpassword")
		if !errors.Is(err, auth.ErrInvalidEmailOrPassword) {
			t.Errorf("Expected error %v, got %v", user.ErrInvalidPassword, err)
		}
	})
}
