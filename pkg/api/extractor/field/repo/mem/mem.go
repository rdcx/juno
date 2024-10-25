package mem

import (
	"juno/pkg/api/extractor/field"

	"github.com/google/uuid"
)

type Repository struct {
	fields map[uuid.UUID]field.Field
}

func New() *Repository {
	return &Repository{
		fields: make(map[uuid.UUID]field.Field),
	}
}

func (r *Repository) Create(q *field.Field) error {
	r.fields[q.ID] = *q
	return nil
}

func (r *Repository) Get(id uuid.UUID) (*field.Field, error) {
	q, ok := r.fields[id]
	if !ok {
		return nil, field.ErrNotFound
	}

	return &q, nil
}

func (r *Repository) ListByUserID(userID uuid.UUID) ([]*field.Field, error) {
	var fields []*field.Field
	for _, q := range r.fields {
		if q.UserID == userID {
			fields = append(fields, &q)
		}
	}

	return fields, nil
}

func (r *Repository) ListBySelectorID(selectorID uuid.UUID) ([]*field.Field, error) {
	var fields []*field.Field
	for _, q := range r.fields {
		if q.SelectorID == selectorID {
			fields = append(fields, &q)
		}
	}

	return fields, nil
}

func (r *Repository) Update(q *field.Field) error {
	r.fields[q.ID] = *q
	return nil
}
