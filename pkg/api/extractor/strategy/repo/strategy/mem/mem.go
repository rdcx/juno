package mem

import (
	"juno/pkg/api/extractor/strategy"

	"github.com/google/uuid"
)

type Repository struct {
	strategys map[uuid.UUID]strategy.Strategy
}

func New() *Repository {
	return &Repository{
		strategys: make(map[uuid.UUID]strategy.Strategy),
	}
}

func (r *Repository) Create(q *strategy.Strategy) error {
	r.strategys[q.ID] = *q
	return nil
}

func (r *Repository) Get(id uuid.UUID) (*strategy.Strategy, error) {
	q, ok := r.strategys[id]
	if !ok {
		return nil, strategy.ErrNotFound
	}

	return &q, nil
}

func (r *Repository) ListByUserID(userID uuid.UUID) ([]*strategy.Strategy, error) {
	var strategys []*strategy.Strategy
	for _, q := range r.strategys {
		if q.UserID == userID {
			strategys = append(strategys, &q)
		}
	}

	return strategys, nil
}

func (r *Repository) Update(q *strategy.Strategy) error {
	r.strategys[q.ID] = *q
	return nil
}
