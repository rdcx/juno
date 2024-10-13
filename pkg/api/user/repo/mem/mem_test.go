package mem

import (
	"fmt"
	"juno/pkg/api/user"
	"juno/pkg/util"
	"testing"

	"github.com/google/uuid"
)

func randomEmail() string {
	return fmt.Sprintf("%s@%s.com", uuid.New().String(), uuid.New().String())
}

func TestCreate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		pass := "password"
		hash, err := util.BcryptPassword(pass)

		if err != nil {
			t.Errorf("expected err to be nil, got %v", err)
		}

		u := &user.User{
			ID:       uuid.New(),
			Email:    randomEmail(),
			Password: hash,
		}

		repo := New()

		err = repo.Create(u)

		if err != nil {
			t.Errorf("expected err to be nil, got %v", err)
		}

		if u.ID == uuid.Nil {
			t.Errorf("expected user ID to be set, got %v", u.ID)
		}

		usr := repo.users[u.ID]

		if usr == nil {
			t.Errorf("expected user to be in the repo, got %v", usr)
		}

		if usr.Email != u.Email {
			t.Errorf("expected user email to be %s, got %s", u.Email, usr.Email)
		}

		if err := util.CompareBcryptPassword(usr.Password, pass); err != nil {
			t.Errorf("expected password to be hashed, got %v", err)
		}
	})
}

func TestGet(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		pass := "password"
		hash, err := util.BcryptPassword(pass)

		if err != nil {
			t.Errorf("expected err to be nil, got %v", err)
		}

		u := &user.User{
			ID:       uuid.New(),
			Email:    randomEmail(),
			Password: hash,
		}

		repo := New()

		repo.users[u.ID] = u

		usr, err := repo.Get(u.ID)

		if err != nil {
			t.Errorf("expected err to be nil, got %v", err)
		}

		if usr == nil {
			t.Errorf("expected user to be in the repo, got %v", usr)
		}

		if usr.Email != u.Email {
			t.Errorf("expected user email to be %s, got %s", u.Email, usr.Email)
		}

		if u.Password != hash {
			t.Errorf("expected user password to be %s, got %s", hash, usr.Password)
		}
	})
}

func TestFirstWhereEmail(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		pass := "password"
		hash, err := util.BcryptPassword(pass)

		if err != nil {
			t.Errorf("expected err to be nil, got %v", err)
		}

		u := &user.User{
			ID:       uuid.New(),
			Email:    randomEmail(),
			Password: hash,
		}

		repo := New()

		repo.users[u.ID] = u

		usr, err := repo.FirstWhereEmail(u.Email)

		if err != nil {
			t.Errorf("expected err to be nil, got %v", err)
		}

		if usr == nil {
			t.Errorf("expected user to be in the repo, got %v", usr)
		}

		if usr.Email != u.Email {
			t.Errorf("expected user email to be %s, got %s", u.Email, usr.Email)
		}

		if u.Password != hash {
			t.Errorf("expected user password to be %s, got %s", hash, usr.Password)
		}

	})

	t.Run("failure", func(t *testing.T) {
		repo := New()

		email := randomEmail()
		usr, err := repo.FirstWhereEmail(email)

		if err != user.ErrNotFound {
			t.Errorf("expected err to be ErrNotFound, got %v", err)
		}

		if usr != nil {
			t.Errorf("expected user to be nil, got %v", usr)
		}
	})
}

func TestUpdate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		pass := "password"
		hash, err := util.BcryptPassword(pass)

		if err != nil {
			t.Errorf("expected err to be nil, got %v", err)
		}

		u := &user.User{
			ID:       uuid.New(),
			Email:    "ross@example.com",
			Password: hash,
		}

		repo := New()

		copyUser := *u

		repo.users[u.ID] = &copyUser

		newEmail := "steve@example.com"

		u.Email = newEmail

		err = repo.Update(u)

		if err != nil {
			t.Errorf("expected err to be nil, got %v", err)
		}

		usr := repo.users[u.ID]

		if usr == nil {
			t.Errorf("expected user to be in the repo, got %v", usr)
		}

		if usr.Email != newEmail {
			t.Errorf("expected user email to be %s, got %s", newEmail, usr.Email)
		}

		if u.Password != hash {
			t.Errorf("expected user password to be %s, got %s", hash, usr.Password)
		}

	})
}

func TestDelete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		pass := "password"
		hash, err := util.BcryptPassword(pass)

		if err != nil {
			t.Errorf("expected err to be nil, got %v", err)
		}

		u := &user.User{
			ID:       uuid.New(),
			Email:    "alan@example.com",
			Password: hash,
		}

		repo := New()

		repo.users[u.ID] = u

		err = repo.Delete(u.ID)

		if err != nil {
			t.Errorf("expected err to be nil, got %v", err)
		}

		usr, ok := repo.users[u.ID]

		if ok {
			t.Errorf("expected user to be deleted, got %v", usr)
		}

		if usr != nil {
			t.Errorf("expected user to be deleted, got %v", usr)
		}
	})
}
