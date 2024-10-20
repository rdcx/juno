package mem

import (
	"juno/pkg/api/query"

	"github.com/google/uuid"
)

type Repository struct {
	queries map[uuid.UUID]query.Query
}

func New() *Repository {
	return &Repository{
		queries: make(map[uuid.UUID]query.Query),
	}
}

func (r *Repository) Create(q *query.Query) error {
	r.queries[q.ID] = *q
	return nil
}

func (r *Repository) Get(id uuid.UUID) (*query.Query, error) {
	q, ok := r.queries[id]
	if !ok {
		return nil, query.ErrQueryNotFound
	}

	return &q, nil
}

func (r *Repository) ListByUserID(userID uuid.UUID) ([]*query.Query, error) {
	var queries []*query.Query
	for _, q := range r.queries {
		if q.UserID == userID {
			queries = append(queries, &q)
		}
	}

	return queries, nil
}

func (r *Repository) Update(q *query.Query) error {
	r.queries[q.ID] = *q
	return nil
}
