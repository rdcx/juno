package mem

import (
	"juno/pkg/api/extractor/selector"

	"github.com/google/uuid"
)

type Repository struct {
	selectors map[uuid.UUID][]*selector.Selector
}

func New() *Repository {
	return &Repository{
		selectors: make(map[uuid.UUID][]*selector.Selector),
	}
}

func (r *Repository) AddSelector(strategyID, selectorID uuid.UUID) error {
	r.selectors[strategyID] = append(r.selectors[strategyID], &selector.Selector{ID: selectorID})
	return nil
}

func (r *Repository) ListSelectorIDs(strategyID uuid.UUID) ([]uuid.UUID, error) {
	var ids []uuid.UUID
	for _, f := range r.selectors[strategyID] {
		ids = append(ids, f.ID)
	}

	return ids, nil
}

func (r *Repository) RemoveSelector(strategyID, selectorID uuid.UUID) error {
	var selectors []*selector.Selector
	for _, f := range r.selectors[strategyID] {
		if f.ID != selectorID {
			selectors = append(selectors, f)
		}
	}
	r.selectors[strategyID] = selectors
	return nil
}
