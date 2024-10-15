package mem

import (
	"juno/pkg/api/assignment"

	"github.com/google/uuid"
)

type Repository struct {
	assignments map[uuid.UUID]*assignment.Assignment
}

func New() *Repository {
	return &Repository{
		assignments: make(map[uuid.UUID]*assignment.Assignment),
	}
}

func (r *Repository) Get(id uuid.UUID) (*assignment.Assignment, error) {
	a, ok := r.assignments[id]
	if !ok {
		return nil, assignment.ErrNotFound
	}
	return a, nil
}

func (r *Repository) ListByNodeID(nodeID uuid.UUID) ([]*assignment.Assignment, error) {
	var result []*assignment.Assignment
	for _, a := range r.assignments {
		if a.NodeID == nodeID {
			result = append(result, a)
		}
	}
	return result, nil
}

func (r *Repository) Create(a *assignment.Assignment) error {
	r.assignments[a.ID] = a
	return nil
}

func (r *Repository) Update(a *assignment.Assignment) error {
	r.assignments[a.ID] = a
	return nil
}

func (r *Repository) Delete(id uuid.UUID) error {
	delete(r.assignments, id)
	return nil
}
