package service

import (
	"juno/pkg/api/extractor/strategy"
	"juno/pkg/api/extractor/strategy/repo/mem"
	"strings"
	"testing"

	"github.com/google/uuid"
)

func TestCreate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := mem.New()
		service := New(repo)
		userID := uuid.New()
		selectorID := uuid.New()
		name := "product_title"
		fType := strategy.StrategyTypeInteger

		f, err := service.Create(userID, selectorID, name, fType)

		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}

		if f.UserID != userID {
			t.Errorf("Expected %s, got %s", userID, f.UserID)
		}

		if f.SelectorID != selectorID {
			t.Errorf("Expected %s, got %s", selectorID, f.SelectorID)
		}

		if f.Name != name {
			t.Errorf("Expected %s, got %s", name, f.Name)
		}

		if f.Type != fType {
			t.Errorf("Expected %s, got %s", fType, f.Type)
		}

		if f.ID == uuid.Nil {
			t.Errorf("Expected non-nil, got %s", f.ID)
		}

	})

	t.Run("validation error", func(t *testing.T) {
		repo := mem.New()
		service := New(repo)
		userID := uuid.Nil
		selectorID := uuid.Nil
		name := ""
		fType := strategy.StrategyType("")

		f, err := service.Create(userID, selectorID, name, fType)

		if err == nil {
			t.Errorf("Expected error, got nil")
		}

		if f != nil {
			t.Errorf("Expected nil, got %v", f)
		}

		if !strings.Contains(err.Error(), "name is required") {
			t.Errorf("Expected 'name is required', got %v", err)
		}

		if !strings.Contains(err.Error(), "selector_id is required") {
			t.Errorf("Expected 'selector_id is required', got %v", err)
		}

		if !strings.Contains(err.Error(), "type is required") {
			t.Errorf("Expected 'type is required', got %v", err)
		}
	})
}

func TestGet(t *testing.T) {
	repo := mem.New()
	service := New(repo)
	userID := uuid.New()
	selectorID := uuid.New()
	name := "product_title"
	fType := strategy.StrategyTypeInteger

	f, err := service.Create(userID, selectorID, name, fType)

	if err != nil {
		t.Fatal(err)
	}

	check, err := service.Get(f.ID)

	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}

	if check.ID != f.ID {
		t.Errorf("Expected %s, got %s", f.ID, check.ID)
	}

	if check.UserID != f.UserID {
		t.Errorf("Expected %s, got %s", f.UserID, check.UserID)
	}

	if check.SelectorID != f.SelectorID {
		t.Errorf("Expected %s, got %s", f.SelectorID, check.SelectorID)
	}

	if check.Name != f.Name {
		t.Errorf("Expected %s, got %s", f.Name, check.Name)
	}

	if check.Type != f.Type {
		t.Errorf("Expected %s, got %s", f.Type, check.Type)
	}
}

func TestListByUserID(t *testing.T) {
	repo := mem.New()
	service := New(repo)
	userID := uuid.New()
	selectorID := uuid.New()
	name := "product_title"
	fType := strategy.StrategyTypeInteger

	f, err := service.Create(userID, selectorID, name, fType)

	if err != nil {
		t.Fatal(err)
	}

	list, err := service.ListByUserID(userID)

	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}

	if len(list) != 1 {
		t.Errorf("Expected 1, got %d", len(list))
	}

	if list[0].ID != f.ID {
		t.Errorf("Expected %s, got %s", f.ID, list[0].ID)
	}
}

func TestListBySelectorID(t *testing.T) {
	repo := mem.New()
	service := New(repo)
	userID := uuid.New()
	selectorID := uuid.New()
	name := "product_title"
	fType := strategy.StrategyTypeInteger

	f, err := service.Create(userID, selectorID, name, fType)

	if err != nil {
		t.Fatal(err)
	}

	list, err := service.ListBySelectorID(selectorID)

	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}

	if len(list) != 1 {
		t.Errorf("Expected 1, got %d", len(list))
	}

	if list[0].ID != f.ID {
		t.Errorf("Expected %s, got %s", f.ID, list[0].ID)
	}
}
