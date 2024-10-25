package mem

import (
	"juno/pkg/api/extractor/field"

	"github.com/google/uuid"
)

type Repository struct {
	fields map[uuid.UUID][]*field.Field
}

func New() *Repository {
	return &Repository{
		fields: make(map[uuid.UUID][]*field.Field),
	}
}

func (r *Repository) AddField(strategyID, fieldID uuid.UUID) error {
	r.fields[strategyID] = append(r.fields[strategyID], &field.Field{ID: fieldID})
	return nil
}

func (r *Repository) ListFieldIDs(strategyID uuid.UUID) ([]uuid.UUID, error) {
	var ids []uuid.UUID
	for _, f := range r.fields[strategyID] {
		ids = append(ids, f.ID)
	}

	return ids, nil
}

func (r *Repository) RemoveField(strategyID, fieldID uuid.UUID) error {
	var fields []*field.Field
	for _, f := range r.fields[strategyID] {
		if f.ID != fieldID {
			fields = append(fields, f)
		}
	}
	r.fields[strategyID] = fields
	return nil
}
