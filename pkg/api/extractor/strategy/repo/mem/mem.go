package mem

import (
	"juno/pkg/api/extractor/strategy"

	"github.com/google/uuid"
)

type Repository struct {
	strategys map[uuid.UUID]*strategy.Strategy
}

func New() *Repository {
	return &Repository{
		strategys: make(map[uuid.UUID]*strategy.Strategy),
	}
}

func (r *Repository) Create(strategy *strategy.Strategy) error {
	r.strategys[strategy.ID] = strategy
	return nil
}

func (r *Repository) Get(id uuid.UUID) (*strategy.Strategy, error) {
	ex, ok := r.strategys[id]
	if !ok {
		return nil, strategy.ErrNotFound
	}
	return ex, nil
}

func (r *Repository) ListByUserID(userID uuid.UUID) ([]*strategy.Strategy, error) {
	var strategys []*strategy.Strategy
	for _, e := range r.strategys {
		if e.UserID == userID {
			strategys = append(strategys, e)
		}
	}
	return strategys, nil
}
