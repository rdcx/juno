package policy

import (
	"context"
	"juno/pkg/api/auth"
	"juno/pkg/api/extractor/filter"
	"juno/pkg/api/user"
	"testing"

	"github.com/google/uuid"
)

func TestCreate(t *testing.T) {
	t.Run("allowed", func(t *testing.T) {
		pol := New()

		if pol.CanCreate().Allowed != true {
			t.Errorf("Expected true, got false")
		}
	})
}

func TestRead(t *testing.T) {
	t.Run("allowed", func(t *testing.T) {
		pol := New()

		u := &user.User{
			ID: uuid.New(),
		}
		ctx := auth.WithUser(context.Background(), u)
		if pol.CanRead(ctx, &filter.Filter{UserID: u.ID}).Allowed != true {
			t.Errorf("Expected true, got false")
		}
	})

	t.Run("denied", func(t *testing.T) {
		pol := New()

		u := &user.User{
			ID: uuid.New(),
		}
		ctx := auth.WithUser(context.Background(), u)
		if pol.CanRead(ctx, &filter.Filter{UserID: uuid.New()}).Allowed != false {
			t.Errorf("Expected false, got true")
		}
	})
}

func TestUpdate(t *testing.T) {
	t.Run("allowed", func(t *testing.T) {
		pol := New()

		u := &user.User{
			ID: uuid.New(),
		}
		ctx := auth.WithUser(context.Background(), u)
		if pol.CanUpdate(ctx, &filter.Filter{UserID: u.ID}).Allowed != true {
			t.Errorf("Expected true, got false")
		}
	})

	t.Run("denied", func(t *testing.T) {
		pol := New()

		u := &user.User{
			ID: uuid.New(),
		}
		ctx := auth.WithUser(context.Background(), u)
		if pol.CanUpdate(ctx, &filter.Filter{UserID: uuid.New()}).Allowed != false {
			t.Errorf("Expected false, got true")
		}
	})
}

func TestDelete(t *testing.T) {
	t.Run("allowed", func(t *testing.T) {
		pol := New()

		u := &user.User{
			ID: uuid.New(),
		}
		ctx := auth.WithUser(context.Background(), u)
		if pol.CanDelete(ctx, &filter.Filter{UserID: u.ID}).Allowed != true {
			t.Errorf("Expected true, got false")
		}
	})

	t.Run("denied", func(t *testing.T) {
		pol := New()

		u := &user.User{
			ID: uuid.New(),
		}
		ctx := auth.WithUser(context.Background(), u)
		if pol.CanDelete(ctx, &filter.Filter{UserID: uuid.New()}).Allowed != false {
			t.Errorf("Expected false, got true")
		}
	})
}

func TestList(t *testing.T) {
	t.Run("allowed", func(t *testing.T) {
		pol := New()

		u := &user.User{
			ID: uuid.New(),
		}
		ctx := auth.WithUser(context.Background(), u)
		if pol.CanList(ctx, []*filter.Filter{{UserID: u.ID}}).Allowed != true {
			t.Errorf("Expected true, got false")
		}
	})

	t.Run("denied", func(t *testing.T) {
		pol := New()

		u := &user.User{
			ID: uuid.New(),
		}
		ctx := auth.WithUser(context.Background(), u)
		if pol.CanList(ctx, []*filter.Filter{{UserID: uuid.New()}}).Allowed != false {
			t.Errorf("Expected false, got true")
		}
	})
}
