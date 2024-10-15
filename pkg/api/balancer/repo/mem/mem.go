package mem

import (
	"errors"
	"fmt"
	"juno/pkg/api/balancer"

	"github.com/google/uuid"
)

type Repository struct {
	balancers map[uuid.UUID]*balancer.Balancer
}

func New() *Repository {
	return &Repository{balancers: make(map[uuid.UUID]*balancer.Balancer)}
}

func (r *Repository) Create(n *balancer.Balancer) error {
	if _, ok := r.balancers[n.ID]; ok {
		return errors.New("primary key violation")
	}

	r.balancers[n.ID] = n

	return nil
}

func (r *Repository) Get(id uuid.UUID) (*balancer.Balancer, error) {
	n, ok := r.balancers[id]
	if !ok || n == nil {
		return nil, balancer.ErrNotFound
	}

	return n, nil
}

func (r *Repository) ListByOwnerID(ownerID uuid.UUID) ([]*balancer.Balancer, error) {
	var balancers []*balancer.Balancer

	for _, n := range r.balancers {
		if n.OwnerID == ownerID {
			balancers = append(balancers, n)
		}
	}

	return balancers, nil
}

func (r *Repository) FirstWhereAddress(address string) (*balancer.Balancer, error) {
	for _, n := range r.balancers {
		if n.Address == address {
			return n, nil
		}
	}

	return nil, errors.New("not found")
}

func (r *Repository) Update(n *balancer.Balancer) error {
	if _, ok := r.balancers[n.ID]; !ok {
		return errors.New("not found")
	}

	r.balancers[n.ID] = n

	return nil
}

func (r *Repository) Delete(id uuid.UUID) error {
	if _, ok := r.balancers[id]; !ok {
		return errors.New("not found")
	}

	delete(r.balancers, id)

	fmt.Println("Deleted balancer", id)

	return nil
}
