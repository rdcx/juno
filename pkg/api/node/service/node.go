package service

import (
	"juno/pkg/api/node"
	"juno/pkg/api/node/policy"

	"juno/pkg/api/user"

	"github.com/google/uuid"
)

type Service struct {
	repo node.Repository
}

func New(repo node.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Get(u *user.User, id uuid.UUID) (*node.Node, error) {
	n, err := s.repo.Get(id)

	if err != nil {
		return nil, node.ErrNotFound
	}

	if can := policy.CanRead(u, n); !can {
		return nil, node.ErrUnauthorized
	}

	return n, nil
}

func (s *Service) Create(user *user.User, addr string, shards []int) (*node.Node, error) {
	if found, _ := s.repo.FirstWhereAddress(addr); found != nil {
		return nil, node.ErrAddressExists
	}

	n := &node.Node{
		ID:      uuid.New(),
		OwnerID: user.ID,
		Address: addr,
		Shards:  shards,
	}

	err := s.repo.Create(n)

	if err != nil {
		return nil, err
	}

	return n, nil
}

func (s *Service) Update(u *user.User, n *node.Node) error {
	n, err := s.repo.Get(n.ID)

	if err != nil {
		return err
	}

	if can := policy.CanUpdate(u, n); !can {
		return node.ErrUnauthorized
	}

	if found, _ := s.repo.FirstWhereAddress(n.Address); found != nil && found.ID != n.ID {
		return node.ErrAddressExists
	}

	return s.repo.Update(n)
}

func (s *Service) Delete(u *user.User, id uuid.UUID) error {

	n, err := s.repo.Get(id)

	if err != nil {
		return node.ErrNotFound
	}

	if can := policy.CanDelete(u, n); !can {
		return node.ErrUnauthorized
	}

	return s.repo.Delete(id)
}
