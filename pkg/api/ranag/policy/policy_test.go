package policy

import (
	"context"
	"juno/pkg/api/auth"
	"juno/pkg/api/ranag"
	"juno/pkg/api/user"
	"testing"

	"github.com/google/uuid"
)

func TestCanCreate(t *testing.T) {
	p := Policy{}

	result := p.CanCreate()

	if result.Allowed == false {
		t.Errorf("expected result to be allowed, got denied")
	}
}

func TestCanRead(t *testing.T) {
	t.Run("allowed", func(t *testing.T) {
		p := Policy{}

		u := &user.User{
			ID: uuid.New(),
		}

		result := p.CanRead(auth.WithUser(context.Background(), u), &ranag.Ranag{
			OwnerID: u.ID,
		})

		if result.Allowed == false {
			t.Errorf("expected result to be allowed, got denied")
		}
	})

	t.Run("denied", func(t *testing.T) {
		p := Policy{}

		u := &user.User{
			ID: uuid.New(),
		}

		result := p.CanRead(auth.WithUser(context.Background(), u), &ranag.Ranag{
			OwnerID: uuid.New(),
		})

		if result.Allowed == true {
			t.Errorf("expected result to be denied, got allowed")
		}

		if result.Reason != "user not allowed to read ranag" {
			t.Errorf("expected reason to be 'user not allowed to read ranag', got %s", result.Reason)
		}
	})

	t.Run("user not found in context", func(t *testing.T) {
		p := Policy{}

		result := p.CanRead(context.Background(), &ranag.Ranag{})

		if result.Allowed == true {
			t.Errorf("expected result to be denied, got allowed")
		}

		if result.Error.Error() != user.ErrUserNotFoundInContext.Error() {
			t.Errorf("expected error to be 'user not found in context', got %s", result.Reason)
		}
	})
}

func TestCanUpdate(t *testing.T) {
	t.Run("allowed", func(t *testing.T) {
		p := Policy{}

		u := &user.User{
			ID: uuid.New(),
		}

		result := p.CanUpdate(auth.WithUser(context.Background(), u), &ranag.Ranag{
			OwnerID: u.ID,
		})

		if result.Allowed == false {
			t.Errorf("expected result to be allowed, got denied")
		}
	})

	t.Run("denied", func(t *testing.T) {
		p := Policy{}

		u := &user.User{
			ID: uuid.New(),
		}

		result := p.CanUpdate(auth.WithUser(context.Background(), u), &ranag.Ranag{
			OwnerID: uuid.New(),
		})

		if result.Allowed == true {
			t.Errorf("expected result to be denied, got allowed")
		}

		if result.Reason != "user not allowed to update ranag" {
			t.Errorf("expected reason to be 'user not allowed to update ranag', got %s", result.Reason)
		}
	})

	t.Run("user not found in context", func(t *testing.T) {
		p := Policy{}

		result := p.CanUpdate(context.Background(), &ranag.Ranag{})

		if result.Allowed == true {
			t.Errorf("expected result to be denied, got allowed")
		}

		if result.Error.Error() != user.ErrUserNotFoundInContext.Error() {
			t.Errorf("expected error to be 'user not found in context', got %s", result.Reason)
		}
	})
}

func TestCanDelete(t *testing.T) {
	t.Run("allowed", func(t *testing.T) {
		p := Policy{}

		u := &user.User{
			ID: uuid.New(),
		}

		result := p.CanDelete(auth.WithUser(context.Background(), u), &ranag.Ranag{
			OwnerID: u.ID,
		})

		if result.Allowed == false {
			t.Errorf("expected result to be allowed, got denied")
		}
	})

	t.Run("denied", func(t *testing.T) {
		p := Policy{}

		u := &user.User{
			ID: uuid.New(),
		}

		result := p.CanDelete(auth.WithUser(context.Background(), u), &ranag.Ranag{
			OwnerID: uuid.New(),
		})

		if result.Allowed == true {
			t.Errorf("expected result to be denied, got allowed")
		}

		if result.Reason != "user not allowed to delete ranag" {
			t.Errorf("expected reason to be 'user not allowed to delete ranag', got %s", result.Reason)
		}
	})

	t.Run("user not found in context", func(t *testing.T) {
		p := Policy{}

		result := p.CanDelete(context.Background(), &ranag.Ranag{})

		if result.Allowed == true {
			t.Errorf("expected result to be denied, got allowed")
		}

		if result.Error.Error() != user.ErrUserNotFoundInContext.Error() {
			t.Errorf("expected error to be 'user not found in context', got %s", result.Reason)
		}
	})
}
