package mem

import (
	"juno/pkg/api/extractor/selector"

	"github.com/google/uuid"
)

type Repository struct {
	selectors map[uuid.UUID]selector.Selector
}

func New() *Repository {
	return &Repository{
		selectors: make(map[uuid.UUID]selector.Selector),
	}
}

func (r *Repository) Create(q *selector.Selector) error {
	r.selectors[q.ID] = *q
	return nil
}

func (r *Repository) Get(id uuid.UUID) (*selector.Selector, error) {
	q, ok := r.selectors[id]
	if !ok {
		return nil, selector.ErrNotFound
	}

	return &q, nil
}

func (r *Repository) ListByUserID(userID uuid.UUID) ([]*selector.Selector, error) {
	var selectors []*selector.Selector
	for _, q := range r.selectors {
		if q.UserID == userID {
			selectors = append(selectors, &q)
		}
	}

	return selectors, nil
}

func (r *Repository) Update(q *selector.Selector) error {
	r.selectors[q.ID] = *q
	return nil
}
