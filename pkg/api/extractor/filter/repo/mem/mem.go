package mem

import (
	"juno/pkg/api/extractor/filter"

	"github.com/google/uuid"
)

type Repository struct {
	filters map[uuid.UUID]filter.Filter
}

func New() *Repository {
	return &Repository{
		filters: make(map[uuid.UUID]filter.Filter),
	}
}

func (r *Repository) Create(q *filter.Filter) error {
	r.filters[q.ID] = *q
	return nil
}

func (r *Repository) Get(id uuid.UUID) (*filter.Filter, error) {
	q, ok := r.filters[id]
	if !ok {
		return nil, filter.ErrNotFound
	}

	return &q, nil
}

func (r *Repository) ListByUserID(userID uuid.UUID) ([]*filter.Filter, error) {
	var filters []*filter.Filter
	for _, q := range r.filters {
		if q.UserID == userID {
			filters = append(filters, &q)
		}
	}

	return filters, nil
}

func (r *Repository) Update(q *filter.Filter) error {
	r.filters[q.ID] = *q
	return nil
}
