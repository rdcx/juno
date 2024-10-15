package policy

import (
	"context"
	"juno/pkg/api/assignment"
	"juno/pkg/api/auth"
	"juno/pkg/api/user"
	"testing"

	"github.com/google/uuid"
)

func TestCanCreate(t *testing.T) {
	t.Run("returns allowed", func(t *testing.T) {
		p := Policy{}

		result := p.CanCreate()

		if !result.Allowed {
			t.Errorf("expected allowed, got denied")
		}
	})
}

func TestCanRead(t *testing.T) {
	t.Run("returns allowed if user is owner", func(t *testing.T) {
		p := Policy{}

		u := &user.User{
			ID: uuid.New(),
		}

		result := p.CanRead(auth.WithUser(context.Background(), u), &assignment.Assignment{
			OwnerID: u.ID,
		})

		if !result.Allowed {
			t.Errorf("expected allowed, got denied")
		}
	})

	t.Run("returns denied if user is not owner", func(t *testing.T) {
		p := Policy{}

		u := &user.User{
			ID: uuid.New(),
		}

		result := p.CanRead(auth.WithUser(context.Background(), u), &assignment.Assignment{
			OwnerID: uuid.New(),
		})

		if result.Allowed {
			t.Errorf("expected denied, got allowed")
		}
	})
}

func TestCanUpdate(t *testing.T) {
	t.Run("returns allowed if user is owner", func(t *testing.T) {
		p := Policy{}

		u := &user.User{
			ID: uuid.New(),
		}

		result := p.CanUpdate(auth.WithUser(context.Background(), u), &assignment.Assignment{
			OwnerID: u.ID,
		})

		if !result.Allowed {
			t.Errorf("expected allowed, got denied")
		}
	})

	t.Run("returns denied if user is not owner", func(t *testing.T) {
		p := Policy{}

		u := &user.User{
			ID: uuid.New(),
		}

		result := p.CanUpdate(auth.WithUser(context.Background(), u), &assignment.Assignment{
			OwnerID: uuid.New(),
		})

		if result.Allowed {
			t.Errorf("expected denied, got allowed")
		}
	})
}

func TestCanDelete(t *testing.T) {
	t.Run("returns allowed if user is owner", func(t *testing.T) {
		p := Policy{}

		u := &user.User{
			ID: uuid.New(),
		}

		result := p.CanDelete(auth.WithUser(context.Background(), u), &assignment.Assignment{
			OwnerID: u.ID,
		})

		if !result.Allowed {
			t.Errorf("expected allowed, got denied")
		}
	})

	t.Run("returns denied if user is not owner", func(t *testing.T) {
		p := Policy{}

		u := &user.User{
			ID: uuid.New(),
		}

		result := p.CanDelete(auth.WithUser(context.Background(), u), &assignment.Assignment{
			OwnerID: uuid.New(),
		})

		if result.Allowed {
			t.Errorf("expected denied, got allowed")
		}
	})
}

func TestCanList(t *testing.T) {
	t.Run("returns allowed if user is owner of all assignments", func(t *testing.T) {
		p := Policy{}

		u := &user.User{
			ID: uuid.New(),
		}

		assignments := []*assignment.Assignment{
			{
				OwnerID: u.ID,
			},
			{
				OwnerID: u.ID,
			},
		}

		result := p.CanList(auth.WithUser(context.Background(), u), assignments)

		if !result.Allowed {
			t.Errorf("expected allowed, got denied")
		}

	})

	t.Run("returns denied if user is not owner of all assignments", func(t *testing.T) {
		p := Policy{}

		u := &user.User{
			ID: uuid.New(),
		}

		assignments := []*assignment.Assignment{
			{
				OwnerID: u.ID,
			},
			{
				OwnerID: uuid.New(),
			},
		}

		result := p.CanList(auth.WithUser(context.Background(), u), assignments)

		if result.Allowed {
			t.Errorf("expected denied, got allowed")
		}
	})
}
