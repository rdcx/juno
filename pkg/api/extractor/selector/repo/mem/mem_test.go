package mem

import (
	"juno/pkg/api/extractor/selector"
	"testing"

	"github.com/google/uuid"
)

func Test(t *testing.T) {
	repo := New()

	q := &selector.Selector{
		ID:     uuid.New(),
		UserID: uuid.New(),
	}

	err := repo.Create(q)

	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}

	check, err := repo.Get(q.ID)

	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}

	if check.ID != q.ID {
		t.Errorf("Expected %s, got %s", q.ID, check.ID)
	}

	if check.UserID != q.UserID {
		t.Errorf("Expected %s, got %s", q.UserID, check.UserID)
	}

	list, err := repo.ListByUserID(q.UserID)

	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}

	if len(list) != 1 {
		t.Errorf("Expected 1, got %d", len(list))
	}

	if list[0].ID != q.ID {
		t.Errorf("Expected %s, got %s", q.ID, list[0].ID)
	}

	var copied selector.Selector

	copied.ID = q.ID
	copied.UserID = q.UserID

	q.Visibility = selector.VisibilityPrivate

	err = repo.Update(q)

	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}

	check, err = repo.Get(q.ID)

	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}

	if check.Visibility != q.Visibility {
		t.Errorf("Expected %s, got %s", q.Visibility, check.Visibility)
	}
}
