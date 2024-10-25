package mem

import (
	"juno/pkg/api/extractor/filter"

	"github.com/google/uuid"
)

type Repository struct {
	filters map[uuid.UUID][]*filter.Filter
}

func New() *Repository {
	return &Repository{
		filters: make(map[uuid.UUID][]*filter.Filter),
	}
}

func (r *Repository) AddFilter(strategyID, filterID uuid.UUID) error {
	r.filters[strategyID] = append(r.filters[strategyID], &filter.Filter{ID: filterID})
	return nil
}

func (r *Repository) ListFilterIDs(strategyID uuid.UUID) ([]uuid.UUID, error) {
	var ids []uuid.UUID
	for _, f := range r.filters[strategyID] {
		ids = append(ids, f.ID)
	}

	return ids, nil
}

func (r *Repository) RemoveFilter(strategyID, filterID uuid.UUID) error {
	var filters []*filter.Filter
	for _, f := range r.filters[strategyID] {
		if f.ID != filterID {
			filters = append(filters, f)
		}
	}
	r.filters[strategyID] = filters
	return nil
}
