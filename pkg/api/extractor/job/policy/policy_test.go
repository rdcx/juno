package policy

import (
	"context"
	"juno/pkg/api/auth"
	"juno/pkg/api/extractor/job"
	"juno/pkg/api/user"
	"testing"

	"github.com/google/uuid"
)

func TestCreate(t *testing.T) {
	t.Run("anyone can create a job", func(t *testing.T) {
		p := New()

		result := p.CanCreate()

		if !result.Allowed {
			t.Errorf("Expected allowed, got denied")
		}
	})
}

func TestUpdate(t *testing.T) {
	t.Run("only job owner can update job", func(t *testing.T) {
		p := New()
		userID := uuid.New()
		ctx := auth.WithUser(context.Background(), &user.User{ID: userID})
		result := p.CanUpdate(ctx, &job.Job{UserID: userID})

		if !result.Allowed {
			t.Errorf("Expected allowed, got denied")
		}

		result = p.CanUpdate(ctx, &job.Job{UserID: uuid.New()})
		if result.Allowed {
			t.Errorf("Expected denied, got allowed")
		}
	})
}

func TestGet(t *testing.T) {
	t.Run("only job owner can get job", func(t *testing.T) {
		p := New()
		userID := uuid.New()
		ctx := auth.WithUser(context.Background(), &user.User{ID: userID})
		result := p.CanGet(ctx, &job.Job{UserID: userID})

		if !result.Allowed {
			t.Errorf("Expected allowed, got denied")
		}

		result = p.CanGet(ctx, &job.Job{UserID: uuid.New()})
		if result.Allowed {
			t.Errorf("Expected denied, got allowed")
		}
	})
}

func TestList(t *testing.T) {
	t.Run("only job owner can list jobs", func(t *testing.T) {
		p := New()
		userID := uuid.New()
		ctx := auth.WithUser(context.Background(), &user.User{ID: userID})
		jobs := []*job.Job{{UserID: userID}, {UserID: userID}}
		result := p.CanList(ctx, jobs)

		if !result.Allowed {
			t.Errorf("Expected allowed, got denied")
		}

		jobs = []*job.Job{{UserID: userID}, {UserID: uuid.New()}}
		result = p.CanList(ctx, jobs)
		if result.Allowed {
			t.Errorf("Expected denied, got allowed")
		}
	})
}
