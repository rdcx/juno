package dto

import (
	"fmt"
	"juno/pkg/api/user"
	"testing"

	"github.com/google/uuid"
)

func randomEmail() string {
	return fmt.Sprintf("%s@%s.com", uuid.New().String(), uuid.New().String())
}

func TestNewSuccessGetUserResponse(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		u := &user.User{
			ID:    uuid.New(),
			Email: randomEmail(),
		}

		resp := NewSuccessGetUserResponse(u)

		if resp.Status != "success" {
			t.Errorf("expected status to be success, got %s", resp.Status)
		}

		if resp.User == nil {
			t.Errorf("expected user to be set, got nil")
		}

		if resp.User.ID != u.ID.String() {
			t.Errorf("expected user ID to be %s, got %s", u.ID.String(), resp.User.ID)
		}

		if resp.User.Email != u.Email {
			t.Errorf("expected user email to be %s, got %s", u.Email, resp.User.Email)
		}

	})
}

func TestNewErrorGetUserResponse(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		message := "error message"
		resp := NewErrorGetUserResponse(message)

		if resp.Status != "error" {
			t.Errorf("expected status to be error, got %s", resp.Status)
		}

		if resp.Message != message {
			t.Errorf("expected message to be %s, got %s", message, resp.Message)
		}
	})
}
