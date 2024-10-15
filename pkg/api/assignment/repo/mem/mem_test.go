package mem

import (
	"juno/pkg/api/assignment"
	"testing"

	"github.com/google/uuid"
)

func TestGet(t *testing.T) {
	t.Run("returns assignment by ID", func(t *testing.T) {
		repo := New()
		a := &assignment.Assignment{
			ID:       uuid.New(),
			OwnerID:  uuid.New(),
			EntityID: uuid.New(),
			Offset:   0,
			Length:   10,
		}

		repo.Store(a)

		result, err := repo.Get(a.ID)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if result != a {
			t.Errorf("expected assignment %v, got %v", a, result)
		}
	})

	t.Run("returns error if assignment not found", func(t *testing.T) {
		repo := New()
		id := uuid.New()

		_, err := repo.Get(id)

		if err != assignment.ErrNotFound {
			t.Errorf("expected error %v, got %v", assignment.ErrNotFound, err)
		}
	})
}

func TestListByEntityID(t *testing.T) {
	t.Run("returns assignments for entity ID", func(t *testing.T) {
		repo := New()
		entityID := uuid.New()
		assignments := []*assignment.Assignment{
			{
				ID:       uuid.New(),
				OwnerID:  uuid.New(),
				EntityID: entityID,
				Offset:   0,
				Length:   10,
			},
			{
				ID:       uuid.New(),
				OwnerID:  uuid.New(),
				EntityID: entityID,
				Offset:   10,
				Length:   10,
			},
		}

		for _, a := range assignments {
			repo.Store(a)
		}

		result, err := repo.ListByEntityID(entityID)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(result) != 2 {
			t.Fatalf("expected 2 assignments, got %d", len(result))
		}

		for i, a := range result {
			if a.ID != assignments[i].ID {
				t.Errorf("expected assignment ID %v, got %v", assignments[i].ID, a.ID)
			}
		}

	})
}

func TestStore(t *testing.T) {
	t.Run("stores assignment", func(t *testing.T) {
		repo := New()
		a := &assignment.Assignment{
			ID:       uuid.New(),
			OwnerID:  uuid.New(),
			EntityID: uuid.New(),
			Offset:   0,
			Length:   10,
		}

		err := repo.Store(a)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(repo.assignments) != 1 {
			t.Fatalf("expected 1 assignment, got %d", len(repo.assignments))
		}

		if repo.assignments[a.ID] != a {
			t.Errorf("expected assignment %v, got %v", a, repo.assignments[a.ID])
		}
	})
}

func TestDelete(t *testing.T) {
	t.Run("deletes assignment", func(t *testing.T) {
		repo := New()
		a := &assignment.Assignment{
			ID:       uuid.New(),
			OwnerID:  uuid.New(),
			EntityID: uuid.New(),
			Offset:   0,
			Length:   10,
		}

		repo.Store(a)

		err := repo.Delete(a.ID)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(repo.assignments) != 0 {
			t.Fatalf("expected 0 assignments, got %d", len(repo.assignments))
		}
	})
}
