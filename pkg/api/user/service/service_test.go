package service

import (
	"fmt"
	"juno/pkg/api/user"
	"juno/pkg/api/user/repo/mem"
	"juno/pkg/util"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

func randomEmail() string {
	return fmt.Sprintf("%s@%s.com", uuid.New().String(), uuid.New().String())
}

func TestGet(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := mem.New()

		u := &user.User{
			ID:    uuid.New(),
			Email: randomEmail(),
		}

		err := repo.Create(u)

		if err != nil {
			t.Errorf("expected err to be nil, got %v", err)
		}

		s := New(logrus.New(), repo)

		usr, err := s.Get(u.ID)

		if err != nil {
			t.Errorf("expected err to be nil, got %v", err)
		}

		if usr.ID != u.ID {
			t.Errorf("expected user ID to be %s, got %s", u.ID, usr.ID)
		}

		if usr.Email != u.Email {
			t.Errorf("expected user email to be %s, got %s", u.Email, usr.Email)
		}
	})
}

func TestCreate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := mem.New()

		s := New(logrus.New(), repo)

		email := randomEmail()
		pass := "password"

		u, err := s.Create(email, pass)

		if err != nil {
			t.Errorf("expected err to be nil, got %v", err)
		}

		if u.ID == uuid.Nil {
			t.Errorf("expected user ID to be set, got %s", u.ID)
		}

		u, err = repo.FirstWhereEmail(email)

		if err != nil {
			t.Errorf("expected err to be nil, got %v", err)
		}

		if u.Email != email {
			t.Errorf("expected user email to be %s, got %s", email, u.Email)
		}

		if u.Password == pass {
			t.Errorf("expected user password to be hashed, got %s", u.Password)
		}

		if err := util.CompareBcryptPassword(u.Password, pass); err != nil {
			t.Errorf("expected password to match, got %s", u.Password)
		}
	})

	t.Run("validates", func(t *testing.T) {
		repo := mem.New()

		s := New(logrus.New(), repo)

		email := "invalid"

		_, err := s.Create(email, "short")

		if !strings.Contains(err.Error(), "invalid email") {
			t.Errorf("expected error to contain 'invalid email', got %v", err)
		}

		if !strings.Contains(err.Error(), "invalid password") {
			t.Errorf("expected error to contain 'invalid password', got %v", err)
		}
	})

	t.Run("email exists", func(t *testing.T) {
		repo := mem.New()

		s := New(logrus.New(), repo)

		email := randomEmail()
		pass := "password"

		_, err := s.Create(email, pass)

		if err != nil {
			t.Errorf("expected err to be nil, got %v", err)
		}

		_, err = s.Create(email, pass)

		if err != user.ErrEmailAlreadyExists {
			t.Errorf("expected error to be ErrEmailAlreadyExists, got %v", err)
		}
	})
}
