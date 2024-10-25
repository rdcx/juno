package service

import (
	"juno/pkg/api/extractor/filter"

	"github.com/google/uuid"
)

type Service struct {
	repo filter.Repository
}

func New(repo filter.Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Create(userID, fieldID uuid.UUID, name string, t filter.FilterType, value string) (*filter.Filter, error) {
	sel := &filter.Filter{
		ID:      uuid.New(),
		UserID:  userID,
		FieldID: fieldID,
		Name:    name,
		Value:   value,
		Type:    t,
	}

	err := sel.Validate()

	if err != nil {
		return nil, err
	}

	err = s.repo.Create(sel)

	if err != nil {
		return nil, err
	}

	return sel, nil
}

func (s *Service) Get(id uuid.UUID) (*filter.Filter, error) {
	return s.repo.Get(id)
}

func (s *Service) ListByUserID(userID uuid.UUID) ([]*filter.Filter, error) {
	return s.repo.ListByUserID(userID)
}
