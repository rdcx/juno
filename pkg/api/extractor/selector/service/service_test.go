package service

import (
	"juno/pkg/api/extractor/selector"
	"juno/pkg/api/extractor/selector/repo/mem"
	"strings"
	"testing"

	"github.com/google/uuid"
)

func TestCreate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := mem.New()
		service := New(repo)
		userID := uuid.New()
		name := "name"
		value := "#productTitle"
		visibility := selector.VisibilityPrivate

		j, err := service.Create(userID, name, value, visibility)

		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}

		if j.UserID != userID {
			t.Errorf("Expected %s, got %s", userID, j.UserID)
		}

		if j.Name != "name" {
			t.Errorf("Expected name, got %s", j.Name)
		}

		if j.Value != "#productTitle" {
			t.Errorf("Expected #productTitle, got %s", j.Value)
		}

		if j.Visibility != selector.VisibilityPrivate {
			t.Errorf("Expected %s, got %s", selector.VisibilityPrivate, j.Visibility)
		}

		if j.ID == uuid.Nil {
			t.Errorf("Expected not empty, got %s", j.ID)
		}
	})

	t.Run("validates", func(t *testing.T) {
		repo := mem.New()
		service := New(repo)
		userID := uuid.New()
		name := ""
		value := ""
		visibility := selector.VisibilityPrivate

		_, err := service.Create(userID, name, value, visibility)

		if err == nil {
			t.Errorf("Expected error, got nil")
		}

		if !strings.Contains(err.Error(), "name is required") {
			t.Errorf("Expected name is required, got %s", err.Error())
		}

		if !strings.Contains(err.Error(), "value is required") {
			t.Errorf("Expected value is required, got %s", err.Error())
		}
	})
}

func TestGet(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := mem.New()
		service := New(repo)
		userID := uuid.New()
		name := "name"
		value := "#productTitle"
		visibility := selector.VisibilityPrivate

		j, err := service.Create(userID, name, value, visibility)

		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}

		check, err := service.Get(j.ID)

		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}

		if check.ID != j.ID {
			t.Errorf("Expected %s, got %s", j.ID, check.ID)
		}

		if check.UserID != j.UserID {
			t.Errorf("Expected %s, got %s", j.UserID, check.UserID)
		}

		if check.Name != j.Name {
			t.Errorf("Expected %s, got %s", j.Name, check.Name)
		}

		if check.Value != j.Value {
			t.Errorf("Expected %s, got %s", j.Value, check.Value)
		}

		if check.Visibility != j.Visibility {
			t.Errorf("Expected %s, got %s", j.Visibility, check.Visibility)
		}
	})

	t.Run("not found", func(t *testing.T) {
		repo := mem.New()
		service := New(repo)
		_, err := service.Get(uuid.New())

		if err != selector.ErrNotFound {
			t.Errorf("Expected %v, got %v", selector.ErrNotFound, err)
		}
	})
}

func TestListByUserID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := mem.New()
		service := New(repo)
		userID := uuid.New()
		name := "name"
		value := "#productTitle"
		visibility := selector.VisibilityPrivate

		j, err := service.Create(userID, name, value, visibility)
		service.Create(uuid.New(), name, value, visibility)

		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}

		list, err := service.ListByUserID(userID)

		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}

		if len(list) != 1 {
			t.Errorf("Expected 1, got %d", len(list))
		}

		if list[0].ID != j.ID {
			t.Errorf("Expected %s, got %s", j.ID, list[0].ID)
		}
	})
}
