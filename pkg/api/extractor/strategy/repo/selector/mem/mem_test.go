package mem

import (
	"testing"

	"github.com/google/uuid"
)

func TestAddSelector(t *testing.T) {
	repo := New()
	strategyID := uuid.New()

	err := repo.AddSelector(strategyID, uuid.New())

	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}

	err = repo.AddSelector(strategyID, uuid.New())

	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}

	if len(repo.selectors[strategyID]) != 2 {
		t.Errorf("Expected 2, got %d", len(repo.selectors[strategyID]))
	}
}

func TestRemoveSelector(t *testing.T) {
	repo := New()
	strategyID := uuid.New()

	selectorID := uuid.New()

	err := repo.AddSelector(strategyID, selectorID)

	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}

	err = repo.RemoveSelector(strategyID, selectorID)

	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}

	if len(repo.selectors[strategyID]) != 0 {
		t.Errorf("Expected 0, got %d", len(repo.selectors[strategyID]))
	}

}

func TestListSelectorIDs(t *testing.T) {
	repo := New()
	strategyID := uuid.New()

	selectorID := uuid.New()

	err := repo.AddSelector(strategyID, selectorID)

	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}

	list, err := repo.ListSelectorIDs(strategyID)

	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}

	if len(list) != 1 {
		t.Errorf("Expected 1, got %d", len(list))
	}

	if list[0] != selectorID {
		t.Errorf("Expected %s, got %s", selectorID, list[0])
	}
}
