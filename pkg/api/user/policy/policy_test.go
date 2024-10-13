package policy

import (
	"context"
	"juno/pkg/api/auth"
	"juno/pkg/api/user"
	"juno/pkg/can"
	"testing"

	"github.com/google/uuid"
)

func testResult(t *testing.T, result can.Result, allowed bool, expectedReason string, err error) {
	if result.Allowed != allowed {
		t.Fatalf("Expected result to be %t, got %t", allowed, result.Allowed)
	}

	if result.Reason != expectedReason {
		t.Fatalf("Expected reason to be %s, got %s", expectedReason, result.Reason)
	}

	if result.Error != err {
		t.Fatalf("Expected error to be %v, got %v", err, result.Error)
	}
}

func TestCreate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		p := Policy{}
		result := p.CanCreate()

		testResult(t, result, true, "", nil)
	})
}

func TestUpdate(t *testing.T) {
	t.Run("allowed", func(t *testing.T) {
		p := Policy{}
		u := &user.User{ID: uuid.New()}

		ctx := auth.WithUser(context.Background(), u)
		result := p.CanUpdate(ctx, u)

		testResult(t, result, true, "", nil)
	})

	t.Run("denied", func(t *testing.T) {
		p := Policy{}
		u := &user.User{ID: uuid.New()}

		ctx := auth.WithUser(context.Background(), &user.User{ID: uuid.New()})
		result := p.CanUpdate(ctx, u)

		testResult(t, result, false, "user not allowed to update user", nil)
	})

	t.Run("error", func(t *testing.T) {
		p := Policy{}
		u := &user.User{ID: uuid.New()}

		result := p.CanUpdate(context.Background(), u)

		testResult(t, result, false, "", user.ErrUserNotFoundInContext)
	})
}

func TestRead(t *testing.T) {
	t.Run("allowed", func(t *testing.T) {
		p := Policy{}
		u := &user.User{ID: uuid.New()}

		ctx := auth.WithUser(context.Background(), u)
		result := p.CanRead(ctx, u)

		testResult(t, result, true, "", nil)
	})

	t.Run("denied", func(t *testing.T) {
		p := Policy{}
		u := &user.User{ID: uuid.New()}

		ctx := auth.WithUser(context.Background(), &user.User{ID: uuid.New()})
		result := p.CanRead(ctx, u)

		testResult(t, result, false, "user not allowed to read user", nil)
	})

	t.Run("error", func(t *testing.T) {
		p := Policy{}
		u := &user.User{ID: uuid.New()}

		result := p.CanRead(context.Background(), u)

		testResult(t, result, false, "", user.ErrUserNotFoundInContext)
	})
}

func TestDelete(t *testing.T) {
	t.Run("allowed", func(t *testing.T) {
		p := Policy{}
		u := &user.User{ID: uuid.New()}

		ctx := auth.WithUser(context.Background(), u)
		result := p.CanDelete(ctx, u)

		testResult(t, result, true, "", nil)
	})

	t.Run("denied", func(t *testing.T) {
		p := Policy{}
		u := &user.User{ID: uuid.New()}

		ctx := auth.WithUser(context.Background(), &user.User{ID: uuid.New()})
		result := p.CanDelete(ctx, u)

		testResult(t, result, false, "user not allowed to delete user", nil)
	})

	t.Run("error", func(t *testing.T) {
		p := Policy{}
		u := &user.User{ID: uuid.New()}

		result := p.CanDelete(context.Background(), u)

		testResult(t, result, false, "", user.ErrUserNotFoundInContext)
	})
}
