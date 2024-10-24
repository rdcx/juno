package service

import (
	"juno/pkg/api/extractor/selector"

	"github.com/google/uuid"
)

type Service struct {
	repo selector.Repository
}

func New(repo selector.Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Create(userID uuid.UUID, name, value string, visibility selector.Visibility) (*selector.Selector, error) {
	sel := &selector.Selector{
		ID:         uuid.New(),
		UserID:     userID,
		Name:       name,
		Value:      value,
		Visibility: visibility,
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

func (s *Service) Get(id uuid.UUID) (*selector.Selector, error) {
	return s.repo.Get(id)
}

func (s *Service) ListByUserID(userID uuid.UUID) ([]*selector.Selector, error) {
	return s.repo.ListByUserID(userID)
}
