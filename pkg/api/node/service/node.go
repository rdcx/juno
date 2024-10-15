package service

import (
	"juno/pkg/api/node"
	"juno/pkg/util"
	"strings"

	"github.com/google/uuid"
)

type Service struct {
	repo node.Repository
}

func New(repo node.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Get(id uuid.UUID) (*node.Node, error) {
	n, err := s.repo.Get(id)

	if err != nil {
		return nil, node.ErrNotFound
	}
	return n, nil
}

func (s *Service) ListByOwnerID(ownerID uuid.UUID) ([]*node.Node, error) {
	nodes, err := s.repo.ListByOwnerID(ownerID)

	if err != nil {
		return nil, node.ErrNotFound
	}

	return nodes, nil
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

func (s *Service) Create(ownerID uuid.UUID, addr string) (*node.Node, error) {
	if found, _ := s.repo.FirstWhereAddress(addr); found != nil {
		return nil, node.ErrAddressExists
	}

	errs := []error{}

	if err := validateAddress(addr); err != nil {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return nil, util.ValidationErrs(errs)
	}

	n := &node.Node{
		ID:      uuid.New(),
		OwnerID: ownerID,
		Address: addr,
	}

	err := s.repo.Create(n)

	if err != nil {
		return nil, err
	}

	return n, nil
}

func (s *Service) Update(id uuid.UUID, dirty *node.Node) (*node.Node, error) {
	n, err := s.repo.Get(id)

	if err != nil {
		return nil, err
	}

	errs := []error{}

	if err := validateAddress(dirty.Address); err != nil {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return nil, util.ValidationErrs(errs)
	}

	if found, _ := s.repo.FirstWhereAddress(dirty.Address); found != nil && found.ID != n.ID {
		return nil, node.ErrAddressExists
	}

	n.Address = dirty.Address

	err = s.repo.Update(n)

	if err != nil {
		return nil, err
	}

	return n, nil
}

func (s *Service) Delete(id uuid.UUID) error {

	n, err := s.repo.Get(id)

	if err != nil {
		return node.ErrNotFound
	}

	return s.repo.Delete(n.ID)
}
