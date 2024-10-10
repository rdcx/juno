package service

import (
	"juno/pkg/api/node"
	"juno/pkg/api/node/policy"
	"juno/pkg/util"
	"strings"

	"juno/pkg/api/user"

	"github.com/google/uuid"

	sharddomain "juno/pkg/shard/domain"
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

func validateShards(shards []int) error {
	for _, s := range shards {
		if s > sharddomain.SHARDS || s < 0 {
			return node.ErrInvalidShards
		}
	}
	return nil
}

func validateAddress(addr string) error {
	addSplit := strings.Split(addr, ":")

	if len(addSplit) != 2 {
		return node.ErrInvalidAddress
	}

	host := addSplit[0]
	port := addSplit[1]

	if !util.IsValidHostname(host) || !util.IsValidPort(port) {
		return node.ErrInvalidAddress
	}

	return nil
}

func (s *Service) Create(user *user.User, addr string, shards []int) (*node.Node, error) {
	if found, _ := s.repo.FirstWhereAddress(addr); found != nil {
		return nil, node.ErrAddressExists
	}

	errs := []error{}

	if err := validateShards(shards); err != nil {
		errs = append(errs, err)
	}

	if err := validateAddress(addr); err != nil {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return nil, util.ValidationErrs(errs)
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

func (s *Service) Update(u *user.User, dirty *node.Node) error {
	n, err := s.repo.Get(dirty.ID)

	if err != nil {
		return err
	}

	if can := policy.CanUpdate(u, n); !can {
		return node.ErrUnauthorized
	}

	errs := []error{}

	if err := validateShards(dirty.Shards); err != nil {
		errs = append(errs, err)
	}

	if err := validateAddress(dirty.Address); err != nil {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return util.ValidationErrs(errs)
	}

	if found, _ := s.repo.FirstWhereAddress(n.Address); found != nil && found.ID != n.ID {
		return node.ErrAddressExists
	}

	n.Address = dirty.Address
	n.Shards = dirty.Shards

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
