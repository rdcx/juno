package service

import (
	"juno/pkg/api/extractor/field"
	"juno/pkg/api/extractor/filter"
	"juno/pkg/api/extractor/selector"
	"juno/pkg/api/extractor/strategy"

	"github.com/google/uuid"
)

type Service struct {
	strategyRepo      strategy.Repository
	stratFilterRepo   strategy.StrategyFilterRepository
	stratFieldRepo    strategy.StrategyFieldRepository
	stratSelectorRepo strategy.StrategySelectorRepository

	filterService   filter.Service
	fieldService    field.Service
	selectorService selector.Service
}

func New(
	strategyRepo strategy.Repository,
	stratFilterRepo strategy.StrategyFilterRepository,
	stratFieldRepo strategy.StrategyFieldRepository,
	stratSelectorRepo strategy.StrategySelectorRepository,
	filterService filter.Service,
	fieldService field.Service,
	selectorService selector.Service,
) *Service {
	return &Service{
		strategyRepo:      strategyRepo,
		stratFilterRepo:   stratFilterRepo,
		stratFieldRepo:    stratFieldRepo,
		stratSelectorRepo: stratSelectorRepo,
		filterService:     filterService,
		fieldService:      fieldService,
		selectorService:   selectorService,
	}
}

func (s *Service) Create(userID uuid.UUID, name string) (*strategy.Strategy, error) {
	sel := &strategy.Strategy{
		ID:     uuid.New(),
		UserID: userID,
		Name:   name,
	}

	err := sel.Validate()

	if err != nil {
		return nil, err
	}

	err = s.strategyRepo.Create(sel)

	if err != nil {
		return nil, err
	}

	return sel, nil
}

func (s *Service) Get(id uuid.UUID) (*strategy.Strategy, error) {
	strat, err := s.strategyRepo.Get(id)

	if err != nil {
		return nil, err
	}

	selectorIDs, err := s.stratSelectorRepo.ListSelectorIDs(id)

	if err != nil {
		return nil, err
	}

	for _, selectorID := range selectorIDs {
		selector, err := s.selectorService.Get(selectorID)

		if err != nil {
			return nil, err
		}

		strat.Selectors = append(strat.Selectors, selector)
	}

	filterIDs, err := s.stratFilterRepo.ListFilterIDs(id)

	if err != nil {
		return nil, err
	}

	for _, filterID := range filterIDs {
		filter, err := s.filterService.Get(filterID)

		if err != nil {
			return nil, err
		}

		strat.Filters = append(strat.Filters, filter)
	}

	fieldIDs, err := s.stratFieldRepo.ListFieldIDs(id)

	if err != nil {
		return nil, err
	}

	for _, fieldID := range fieldIDs {
		field, err := s.fieldService.Get(fieldID)

		if err != nil {
			return nil, err
		}

		strat.Fields = append(strat.Fields, field)
	}

	return strat, nil
}

func (s *Service) Update(strat *strategy.Strategy) error {

	if _, err := s.strategyRepo.Get(strat.ID); err != nil {
		return err
	}

	return s.strategyRepo.Update(strat)
}

func (s *Service) ListByUserID(userID uuid.UUID) ([]*strategy.Strategy, error) {
	strats, err := s.strategyRepo.ListByUserID(userID)

	if err != nil {
		return nil, err
	}

	for _, strat := range strats {
		selectorIDs, err := s.stratSelectorRepo.ListSelectorIDs(strat.ID)

		if err != nil {
			return nil, err
		}

		for _, selectorID := range selectorIDs {
			selector, err := s.selectorService.Get(selectorID)

			if err != nil {
				return nil, err
			}

			strat.Selectors = append(strat.Selectors, selector)
		}

		filterIDs, err := s.stratFilterRepo.ListFilterIDs(strat.ID)

		if err != nil {
			return nil, err
		}

		for _, filterID := range filterIDs {
			filter, err := s.filterService.Get(filterID)

			if err != nil {
				return nil, err
			}

			strat.Filters = append(strat.Filters, filter)
		}

		fieldIDs, err := s.stratFieldRepo.ListFieldIDs(strat.ID)

		if err != nil {
			return nil, err
		}

		for _, fieldID := range fieldIDs {
			field, err := s.fieldService.Get(fieldID)

			if err != nil {
				return nil, err
			}

			strat.Fields = append(strat.Fields, field)
		}
	}

	return strats, nil
}

func (s *Service) Delete(id uuid.UUID) error {
	if _, err := s.strategyRepo.Get(id); err != nil {
		return err
	}
	return s.strategyRepo.Delete(id)
}

func (s *Service) AddSelector(id, selectorID uuid.UUID) error {
	if _, err := s.strategyRepo.Get(id); err != nil {
		return err
	}

	if _, err := s.selectorService.Get(selectorID); err != nil {
		return err
	}

	return s.stratSelectorRepo.AddSelector(id, selectorID)
}

func (s *Service) RemoveSelector(id, selectorID uuid.UUID) error {
	if _, err := s.strategyRepo.Get(id); err != nil {
		return err
	}

	if _, err := s.selectorService.Get(selectorID); err != nil {
		return err
	}

	return s.stratSelectorRepo.RemoveSelector(id, selectorID)
}

func (s *Service) AddFilter(id, filterID uuid.UUID) error {
	if _, err := s.strategyRepo.Get(id); err != nil {
		return err
	}

	if _, err := s.filterService.Get(filterID); err != nil {
		return err
	}

	return s.stratFilterRepo.AddFilter(id, filterID)
}

func (s *Service) RemoveFilter(id, filterID uuid.UUID) error {
	if _, err := s.strategyRepo.Get(id); err != nil {
		return err
	}

	if _, err := s.filterService.Get(filterID); err != nil {
		return err
	}

	return s.stratFilterRepo.RemoveFilter(id, filterID)
}

func (s *Service) AddField(id, fieldID uuid.UUID) error {
	if _, err := s.strategyRepo.Get(id); err != nil {
		return err
	}

	if _, err := s.fieldService.Get(fieldID); err != nil {
		return err
	}

	return s.stratFieldRepo.AddField(id, fieldID)
}

func (s *Service) RemoveField(id, fieldID uuid.UUID) error {
	if _, err := s.strategyRepo.Get(id); err != nil {
		return err
	}

	if _, err := s.fieldService.Get(fieldID); err != nil {
		return err
	}

	return s.stratFieldRepo.RemoveField(id, fieldID)
}
