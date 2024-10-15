package mem

import (
	"errors"
	"fmt"
	"juno/pkg/api/node"

	"github.com/google/uuid"
)

type Repository struct {
	nodes map[uuid.UUID]*node.Node
}

func New() *Repository {
	return &Repository{nodes: make(map[uuid.UUID]*node.Node)}
}

func (r *Repository) Create(n *node.Node) error {
	if _, ok := r.nodes[n.ID]; ok {
		return errors.New("primary key violation")
	}

	r.nodes[n.ID] = n

	return nil
}

func (r *Repository) All() ([]*node.Node, error) {
	var nodes []*node.Node

	for _, n := range r.nodes {
		nodes = append(nodes, n)
	}

	return nodes, nil
}

func (r *Repository) Get(id uuid.UUID) (*node.Node, error) {
	n, ok := r.nodes[id]
	if !ok || n == nil {
		return nil, node.ErrNotFound
	}

	return n, nil
}

func (r *Repository) ListByOwnerID(ownerID uuid.UUID) ([]*node.Node, error) {
	var nodes []*node.Node

	for _, n := range r.nodes {
		if n.OwnerID == ownerID {
			nodes = append(nodes, n)
		}
	}

	return nodes, nil
}

func (r *Repository) FirstWhereAddress(address string) (*node.Node, error) {
	for _, n := range r.nodes {
		if n.Address == address {
			return n, nil
		}
	}

	return nil, errors.New("not found")
}

func (r *Repository) Update(n *node.Node) error {
	if _, ok := r.nodes[n.ID]; !ok {
		return errors.New("not found")
	}

	r.nodes[n.ID] = n

	return nil
}

func (r *Repository) Delete(id uuid.UUID) error {
	if _, ok := r.nodes[id]; !ok {
		return errors.New("not found")
	}

	delete(r.nodes, id)

	fmt.Println("Deleted node", id)

	return nil
}
