package mem

import (
	"juno/pkg/api/extraction/strategy"
	"testing"

	"github.com/google/uuid"
)

func Test(t *testing.T) {
	repo := New()

	ex := &strategy.Strategy{
		ID:     uuid.New(),
		UserID: uuid.New(),
		Name:   "name",
	}

	err := repo.Create(ex)

	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}

	ex2, err := repo.Get(ex.ID)

	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}

	if ex2.ID != ex.ID {
		t.Errorf("Expected %s, got %s", ex.ID, ex2.ID)
	}

	repo.Create(&strategy.Strategy{
		ID:     uuid.New(),
		UserID: uuid.New(),
	})

	list, err := repo.ListByUserID(ex.UserID)

	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}

	if len(list) != 1 {
		t.Errorf("Expected 1, got %d", len(list))
	}

	if list[0].ID != ex.ID {
		t.Errorf("Expected %s, got %s", ex.ID, list[0].ID)
	}
}
