package mem

import (
	"testing"

	"github.com/google/uuid"
)

func TestAddField(t *testing.T) {
	repo := New()
	strategyID := uuid.New()

	err := repo.AddField(strategyID, uuid.New())

	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}

	err = repo.AddField(strategyID, uuid.New())

	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}

	if len(repo.fields[strategyID]) != 2 {
		t.Errorf("Expected 2, got %d", len(repo.fields[strategyID]))
	}
}

func TestRemoveField(t *testing.T) {
	repo := New()
	strategyID := uuid.New()

	fieldID := uuid.New()

	err := repo.AddField(strategyID, fieldID)

	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}

	err = repo.RemoveField(strategyID, fieldID)

	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}

	if len(repo.fields[strategyID]) != 0 {
		t.Errorf("Expected 0, got %d", len(repo.fields[strategyID]))
	}

}

func TestListFieldIDs(t *testing.T) {
	repo := New()
	strategyID := uuid.New()

	fieldID := uuid.New()

	err := repo.AddField(strategyID, fieldID)

	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}

	list, err := repo.ListFieldIDs(strategyID)

	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}

	if len(list) != 1 {
		t.Errorf("Expected 1, got %d", len(list))
	}

	if list[0] != fieldID {
		t.Errorf("Expected %s, got %s", fieldID, list[0])
	}
}
