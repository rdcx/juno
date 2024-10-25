package mem

import (
	"errors"
	"fmt"
	"juno/pkg/api/ranag"

	"github.com/google/uuid"
)

type Repository struct {
	ranags map[uuid.UUID]*ranag.Ranag
}

func New() *Repository {
	return &Repository{ranags: make(map[uuid.UUID]*ranag.Ranag)}
}

func (r *Repository) Create(n *ranag.Ranag) error {
	if _, ok := r.ranags[n.ID]; ok {
		return errors.New("primary key violation")
	}

	r.ranags[n.ID] = n

	return nil
}

func (r *Repository) All() ([]*ranag.Ranag, error) {
	var ranags []*ranag.Ranag

	for _, n := range r.ranags {
		ranags = append(ranags, n)
	}

	return ranags, nil
}

func (r *Repository) Get(id uuid.UUID) (*ranag.Ranag, error) {
	n, ok := r.ranags[id]
	if !ok || n == nil {
		return nil, ranag.ErrNotFound
	}

	return n, nil
}

func (r *Repository) ListByOwnerID(ownerID uuid.UUID) ([]*ranag.Ranag, error) {
	var ranags []*ranag.Ranag

	for _, n := range r.ranags {
		if n.OwnerID == ownerID {
			ranags = append(ranags, n)
		}
	}

	return ranags, nil
}

func (r *Repository) FirstWhereAddress(address string) (*ranag.Ranag, error) {
	for _, n := range r.ranags {
		if n.Address == address {
			return n, nil
		}
	}

	return nil, errors.New("not found")
}

func (r *Repository) Update(n *ranag.Ranag) error {
	if _, ok := r.ranags[n.ID]; !ok {
		return errors.New("not found")
	}

	r.ranags[n.ID] = n

	return nil
}

func (r *Repository) Delete(id uuid.UUID) error {
	if _, ok := r.ranags[id]; !ok {
		return errors.New("not found")
	}

	delete(r.ranags, id)

	fmt.Println("Deleted ranag", id)

	return nil
}
