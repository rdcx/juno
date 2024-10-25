package service

import (
	"juno/pkg/api/extractor/field"

	"github.com/google/uuid"
)

type Service struct {
	repo field.Repository
}

func New(repo field.Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Create(userID, selectorID uuid.UUID, name string, t field.FieldType) (*field.Field, error) {
	sel := &field.Field{
		ID:         uuid.New(),
		UserID:     userID,
		SelectorID: selectorID,
		Name:       name,
		Type:       t,
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

func (s *Service) Get(id uuid.UUID) (*field.Field, error) {
	return s.repo.Get(id)
}

func (s *Service) ListByUserID(userID uuid.UUID) ([]*field.Field, error) {
	return s.repo.ListByUserID(userID)
}

func (s *Service) ListBySelectorID(selectorID uuid.UUID) ([]*field.Field, error) {
	return s.repo.ListBySelectorID(selectorID)
}
