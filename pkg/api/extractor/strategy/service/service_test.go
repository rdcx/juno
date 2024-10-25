package service

import (
	"juno/pkg/api/extractor/strategy"
	"juno/pkg/api/extractor/strategy/repo/strategy/mem"
	"testing"

	stratFieldRepo "juno/pkg/api/extractor/strategy/repo/field/mem"
	stratFilterRepo "juno/pkg/api/extractor/strategy/repo/filter/mem"
	stratSelectorRepo "juno/pkg/api/extractor/strategy/repo/selector/mem"

	fieldRepo "juno/pkg/api/extractor/field/repo/mem"
	filterRepo "juno/pkg/api/extractor/filter/repo/mem"
	selectorRepo "juno/pkg/api/extractor/selector/repo/mem"

	fieldService "juno/pkg/api/extractor/field/service"
	filterService "juno/pkg/api/extractor/filter/service"
	selectorService "juno/pkg/api/extractor/selector/service"

	"github.com/google/uuid"
)

func setup() (*Service, *mem.Repository, *stratFieldRepo.Repository, *stratSelectorRepo.Repository, *stratFilterRepo.Repository, *filterService.Service, *fieldService.Service, *selectorService.Service) {
	stratFieldRepo := stratFieldRepo.New()
	stratSelectorRepo := stratSelectorRepo.New()
	stratFilterRepo := stratFilterRepo.New()
	strategyRepo := mem.New()

	filterService := filterService.New(filterRepo.New())
	fieldService := fieldService.New(fieldRepo.New())
	selectorService := selectorService.New(selectorRepo.New())

	service := New(
		strategyRepo,
		stratFilterRepo,
		stratFieldRepo,
		stratSelectorRepo,
		filterService,
		fieldService,
		selectorService)

	return service, strategyRepo, stratFieldRepo, stratSelectorRepo, stratFilterRepo, filterService, fieldService, selectorService
}

func TestCreate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		service, _, _, _, _, _, _, _ := setup()

		userID := uuid.New()
		name := "strategy_name"

		strat, err := service.Create(userID, name)

		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}

		if strat.UserID != userID {
			t.Errorf("Expected %s, got %s", userID, strat.UserID)
		}

		if strat.Name != name {
			t.Errorf("Expected %s, got %s", name, strat.Name)
		}

		if strat.ID == uuid.Nil {
			t.Errorf("Expected non-nil, got %s", strat.ID)
		}
	})
}

func TestGet(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		service, strategyRepo, _, _, _, _, _, _ := setup()

		strat := &strategy.Strategy{
			ID:     uuid.New(),
			UserID: uuid.New(),
			Name:   "strategy_name",
		}

		strategyRepo.Create(strat)

		check, err := service.Get(strat.ID)

		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}

		if check.ID != strat.ID {
			t.Errorf("Expected %s, got %s", strat.ID, check.ID)
		}

		if check.UserID != strat.UserID {
			t.Errorf("Expected %s, got %s", strat.UserID, check.UserID)
		}

		if check.Name != strat.Name {
			t.Errorf("Expected %s, got %s", strat.Name, check.Name)
		}
	})

	t.Run("not found", func(t *testing.T) {
		service, _, _, _, _, _, _, _ := setup()

		id := uuid.New()

		_, err := service.Get(id)

		if err == nil {
			t.Errorf("Expected error, got nil")
		}
	})
}

func TestUpdate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		service, strategyRepo, _, _, _, _, _, _ := setup()

		strat := &strategy.Strategy{
			ID:     uuid.New(),
			UserID: uuid.New(),
			Name:   "strategy_name",
		}

		err := strategyRepo.Create(strat)

		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}

		var copied strategy.Strategy

		copied.ID = strat.ID
		copied.UserID = strat.UserID
		copied.Name = "new_name"

		err = service.Update(&copied)

		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}

		check, _ := service.Get(strat.ID)

		if check.Name != "new_name" {
			t.Errorf("Expected %s, got %s", "new_name", check.Name)
		}
	})

	t.Run("not found", func(t *testing.T) {
		service, _, _, _, _, _, _, _ := setup()

		id := uuid.New()

		strat := &strategy.Strategy{
			ID:     id,
			UserID: uuid.New(),
			Name:   "strategy_name",
		}

		err := service.Update(strat)

		if err == nil {
			t.Errorf("Expected error, got nil")
		}
	})
}

func TestDelete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		service, strategyRepo, _, _, _, _, _, _ := setup()

		strat := &strategy.Strategy{
			ID:     uuid.New(),
			UserID: uuid.New(),
			Name:   "strategy_name",
		}

		strategyRepo.Create(strat)

		err := service.Delete(strat.ID)

		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}

		_, err = service.Get(strat.ID)

		if err == nil {
			t.Errorf("Expected error, got nil")
		}
	})

	t.Run("not found", func(t *testing.T) {
		service, _, _, _, _, _, _, _ := setup()

		id := uuid.New()

		err := service.Delete(id)

		if err == nil {
			t.Errorf("Expected error, got nil")
		}
	})
}

func TestAddFilter(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		service, strategyRepo, _, _, _, filterService, _, _ := setup()

		strat := &strategy.Strategy{
			ID:     uuid.New(),
			UserID: uuid.New(),
			Name:   "strategy_name",
		}

		strategyRepo.Create(strat)

		filter, err := filterService.Create(strat.UserID, "filter_name", "filter_type", "filter_value")

		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}

		err = service.AddFilter(strat.ID, filter.ID)

		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}

		check, _ := service.Get(strat.ID)

		if len(check.Filters) != 1 {
			t.Errorf("Expected 1, got %d", len(check.Filters))
		}

		if check.Filters[0].ID != filter.ID {
			t.Errorf("Expected %s, got %s", filter.ID, check.Filters[0])
		}
	})

	t.Run("not found", func(t *testing.T) {
		service, _, _, _, _, _, _, _ := setup()

		id := uuid.New()
		filterID := uuid.New()

		err := service.AddFilter(id, filterID)

		if err == nil {
			t.Errorf("Expected error, got nil")
		}
	})
}

func TestRemoveFilter(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		service, strategyRepo, _, _, _, filterService, _, _ := setup()

		strat := &strategy.Strategy{
			ID:     uuid.New(),
			UserID: uuid.New(),
			Name:   "strategy_name",
		}

		strategyRepo.Create(strat)

		filter, _ := filterService.Create(strat.UserID, "filter_name", "filter_type", "filter_value")

		service.AddFilter(strat.ID, filter.ID)

		err := service.RemoveFilter(strat.ID, filter.ID)

		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}

		check, _ := service.Get(strat.ID)

		if len(check.Filters) != 0 {
			t.Errorf("Expected 0, got %d", len(check.Filters))
		}
	})

	t.Run("not found", func(t *testing.T) {
		service, _, _, _, _, _, _, _ := setup()

		id := uuid.New()
		filterID := uuid.New()

		err := service.RemoveFilter(id, filterID)

		if err == nil {
			t.Errorf("Expected error, got nil")
		}
	})
}

func TestAddField(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		service, strategyRepo, _, _, _, _, fieldService, _ := setup()

		strat := &strategy.Strategy{
			ID:     uuid.New(),
			UserID: uuid.New(),
			Name:   "strategy_name",
		}

		strategyRepo.Create(strat)

		field, err := fieldService.Create(strat.UserID, uuid.New(), "field_name", "field_type")

		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}

		err = service.AddField(strat.ID, field.ID)

		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}

		check, _ := service.Get(strat.ID)

		if len(check.Fields) != 1 {
			t.Errorf("Expected 1, got %d", len(check.Fields))
		}

		if check.Fields[0].ID != field.ID {
			t.Errorf("Expected %s, got %s", field.ID, check.Fields[0])
		}
	})

	t.Run("not found", func(t *testing.T) {
		service, _, _, _, _, _, _, _ := setup()

		id := uuid.New()
		fieldID := uuid.New()

		err := service.AddField(id, fieldID)

		if err == nil {
			t.Errorf("Expected error, got nil")
		}
	})
}

func TestRemoveField(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		service, strategyRepo, _, _, _, _, fieldService, _ := setup()

		strat := &strategy.Strategy{
			ID:     uuid.New(),
			UserID: uuid.New(),
			Name:   "strategy_name",
		}

		strategyRepo.Create(strat)

		field, _ := fieldService.Create(strat.UserID, uuid.New(), "field_name", "field_type")

		service.AddField(strat.ID, field.ID)

		err := service.RemoveField(strat.ID, field.ID)

		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}

		check, _ := service.Get(strat.ID)

		if len(check.Fields) != 0 {
			t.Errorf("Expected 0, got %d", len(check.Fields))
		}
	})

	t.Run("not found", func(t *testing.T) {
		service, _, _, _, _, _, _, _ := setup()

		id := uuid.New()
		fieldID := uuid.New()

		err := service.RemoveField(id, fieldID)

		if err == nil {
			t.Errorf("Expected error, got nil")
		}
	})
}

func TestAddSelector(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		service, strategyRepo, _, _, _, _, _, selectorService := setup()

		strat := &strategy.Strategy{
			ID:     uuid.New(),
			UserID: uuid.New(),
			Name:   "strategy_name",
		}

		strategyRepo.Create(strat)

		selector, err := selectorService.Create(strat.UserID, "selector_name", "selector_type", "selector_value")

		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}

		err = service.AddSelector(strat.ID, selector.ID)

		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}

		check, _ := service.Get(strat.ID)

		if len(check.Selectors) != 1 {
			t.Errorf("Expected 1, got %d", len(check.Selectors))
		}

		if check.Selectors[0].ID != selector.ID {
			t.Errorf("Expected %s, got %s", selector.ID, check.Selectors[0])
		}
	})

	t.Run("not found", func(t *testing.T) {
		service, _, _, _, _, _, _, _ := setup()

		id := uuid.New()
		selectorID := uuid.New()

		err := service.AddSelector(id, selectorID)

		if err == nil {
			t.Errorf("Expected error, got nil")
		}
	})
}

func TestRemoveSelector(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		service, strategyRepo, _, _, _, _, _, selectorService := setup()

		strat := &strategy.Strategy{
			ID:     uuid.New(),
			UserID: uuid.New(),
			Name:   "strategy_name",
		}

		strategyRepo.Create(strat)

		selector, _ := selectorService.Create(strat.UserID, "selector_name", "selector_type", "selector_value")

		service.AddSelector(strat.ID, selector.ID)

		err := service.RemoveSelector(strat.ID, selector.ID)

		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}

		check, _ := service.Get(strat.ID)

		if len(check.Selectors) != 0 {
			t.Errorf("Expected 0, got %d", len(check.Selectors))
		}
	})

	t.Run("not found", func(t *testing.T) {
		service, _, _, _, _, _, _, _ := setup()

		id := uuid.New()
		selectorID := uuid.New()

		err := service.RemoveSelector(id, selectorID)

		if err == nil {
			t.Errorf("Expected error, got nil")
		}
	})
}
