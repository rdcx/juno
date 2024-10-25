package service

import (
	"juno/pkg/api/extractor/filter"
	"juno/pkg/api/extractor/filter/repo/mem"
	"strings"
	"testing"

	"github.com/google/uuid"
)

func TestCreate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := mem.New()
		service := New(repo)
		userID := uuid.New()
		name := "String equals 'charger'"
		value := "charger"
		fType := filter.FilterTypeStringEquals

		f, err := service.Create(userID, name, fType, value)

		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}

		if f.UserID != userID {
			t.Errorf("Expected %s, got %s", userID, f.UserID)
		}

		if f.Name != name {
			t.Errorf("Expected %s, got %s", name, f.Name)
		}

		if f.Type != fType {
			t.Errorf("Expected %s, got %s", fType, f.Type)
		}

		if f.Value != value {
			t.Errorf("Expected %s, got %s", value, f.Value)
		}

		if f.ID == uuid.Nil {
			t.Errorf("Expected non-nil, got %s", f.ID)
		}
	})

	t.Run("validation error", func(t *testing.T) {
		repo := mem.New()
		service := New(repo)
		userID := uuid.New()
		name := ""
		value := ""
		fType := filter.FilterType("")

		f, err := service.Create(userID, name, fType, value)

		if err == nil {
			t.Errorf("Expected error, got nil")
		}

		if f != nil {
			t.Errorf("Expected nil, got %v", f)
		}

		if !strings.Contains(err.Error(), "name is required") {
			t.Errorf("Expected 'name is required', got %v", err)
		}

		if !strings.Contains(err.Error(), "type is required") {
			t.Errorf("Expected 'type is required', got %v", err)
		}

		if !strings.Contains(err.Error(), "value is required") {
			t.Errorf("Expected 'value is required', got %v", err)
		}
	})
}

func TestGet(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := mem.New()
		service := New(repo)
		userID := uuid.New()
		name := "String equals 'charger'"
		value := "charger"
		fType := filter.FilterTypeStringEquals

		f, err := service.Create(userID, name, fType, value)

		if err != nil {
			t.Fatal(err)
		}

		got, err := service.Get(f.ID)

		if err != nil {
			t.Fatal(err)
		}

		if got.ID != f.ID {
			t.Errorf("Expected %s, got %s", f.ID, got.ID)
		}

		if got.UserID != f.UserID {
			t.Errorf("Expected %s, got %s", f.UserID, got.UserID)
		}

		if got.Name != f.Name {
			t.Errorf("Expected %s, got %s", f.Name, got.Name)
		}

		if got.Type != f.Type {
			t.Errorf("Expected %s, got %s", f.Type, got.Type)
		}

		if got.Value != f.Value {
			t.Errorf("Expected %s, got %s", f.Value, got.Value)
		}
	})

	t.Run("not found", func(t *testing.T) {
		repo := mem.New()
		service := New(repo)
		id := uuid.New()

		_, err := service.Get(id)

		if err != filter.ErrNotFound {
			t.Errorf("Expected %v, got %v", filter.ErrNotFound, err)
		}
	})
}

func TestListByUserID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := mem.New()
		service := New(repo)
		userID := uuid.New()
		name := "String equals 'charger'"
		value := "charger"
		fType := filter.FilterTypeStringEquals

		f1, err := service.Create(userID, name, fType, value)

		if err != nil {
			t.Fatal(err)
		}

		name = "String contains 'charger'"
		fType = filter.FilterTypeStringContains

		f2, err := service.Create(userID, name, fType, value)

		if err != nil {
			t.Fatal(err)
		}

		filters, err := service.ListByUserID(userID)

		if err != nil {
			t.Fatal(err)
		}

		if len(filters) != 2 {
			t.Errorf("Expected 2, got %d", len(filters))
		}

		if filters[0].ID != f1.ID {
			t.Errorf("Expected %s, got %s", f1.ID, filters[0].ID)
		}

		if filters[1].ID != f2.ID {
			t.Errorf("Expected %s, got %s", f2.ID, filters[1].ID)
		}
	})

	t.Run("no filters", func(t *testing.T) {
		repo := mem.New()
		service := New(repo)
		userID := uuid.New()

		filters, err := service.ListByUserID(userID)

		if err != nil {
			t.Fatal(err)
		}

		if len(filters) != 0 {
			t.Errorf("Expected 0, got %d", len(filters))
		}
	})
}
