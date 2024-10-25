package mem

import (
	"testing"

	"github.com/google/uuid"
)

func TestAddFilter(t *testing.T) {
	repo := New()
	strategyID := uuid.New()

	err := repo.AddFilter(strategyID, uuid.New())

	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}

	err = repo.AddFilter(strategyID, uuid.New())

	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}

	if len(repo.filters[strategyID]) != 2 {
		t.Errorf("Expected 2, got %d", len(repo.filters[strategyID]))
	}
}

func TestRemoveFilter(t *testing.T) {
	repo := New()
	strategyID := uuid.New()

	filterID := uuid.New()

	err := repo.AddFilter(strategyID, filterID)

	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}

	err = repo.RemoveFilter(strategyID, filterID)

	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}

	if len(repo.filters[strategyID]) != 0 {
		t.Errorf("Expected 0, got %d", len(repo.filters[strategyID]))
	}

}

func TestListFilterIDs(t *testing.T) {
	repo := New()
	strategyID := uuid.New()

	filterID := uuid.New()

	err := repo.AddFilter(strategyID, filterID)

	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}

	list, err := repo.ListFilterIDs(strategyID)

	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}

	if len(list) != 1 {
		t.Errorf("Expected 1, got %d", len(list))
	}

	if list[0] != filterID {
		t.Errorf("Expected %s, got %s", filterID, list[0])
	}
}
