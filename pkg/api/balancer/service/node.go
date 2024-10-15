package service

import (
	"juno/pkg/api/balancer"
	"juno/pkg/util"
	"strings"

	"github.com/google/uuid"

	"juno/pkg/shard"
)

type Service struct {
	repo balancer.Repository
}

func New(repo balancer.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Get(id uuid.UUID) (*balancer.Balancer, error) {
	n, err := s.repo.Get(id)

	if err != nil {
		return nil, balancer.ErrNotFound
	}
	return n, nil
}

func (s *Service) ListByOwnerID(ownerID uuid.UUID) ([]*balancer.Balancer, error) {
	balancers, err := s.repo.ListByOwnerID(ownerID)

	if err != nil {
		return nil, balancer.ErrNotFound
	}

	return balancers, nil
}

func validateShards(shards []int) error {
	for _, s := range shards {
		if s > shard.SHARDS || s < 0 {
			return balancer.ErrInvalidShards
		}
	}
	return nil
}

func validateAddress(addr string) error {
	addSplit := strings.Split(addr, ":")

	if len(addSplit) != 2 {
		return balancer.ErrInvalidAddress
	}

	host := addSplit[0]
	port := addSplit[1]

	if !util.IsValidHostname(host) || !util.IsValidPort(port) {
		return balancer.ErrInvalidAddress
	}

	return nil
}

func (s *Service) Create(ownerID uuid.UUID, addr string, shards []int) (*balancer.Balancer, error) {
	if found, _ := s.repo.FirstWhereAddress(addr); found != nil {
		return nil, balancer.ErrAddressExists
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

	n := &balancer.Balancer{
		ID:      uuid.New(),
		OwnerID: ownerID,
		Address: addr,
		Shards:  shards,
	}

	err := s.repo.Create(n)

	if err != nil {
		return nil, err
	}

	return n, nil
}

func (s *Service) Update(id uuid.UUID, dirty *balancer.Balancer) (*balancer.Balancer, error) {
	n, err := s.repo.Get(id)

	if err != nil {
		return nil, err
	}

	errs := []error{}

	if err := validateShards(dirty.Shards); err != nil {
		errs = append(errs, err)
	}

	if err := validateAddress(dirty.Address); err != nil {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return nil, util.ValidationErrs(errs)
	}

	if found, _ := s.repo.FirstWhereAddress(dirty.Address); found != nil && found.ID != n.ID {
		return nil, balancer.ErrAddressExists
	}

	n.Address = dirty.Address
	n.Shards = dirty.Shards

	err = s.repo.Update(n)

	if err != nil {
		return nil, err
	}

	return n, nil
}

func (s *Service) Delete(id uuid.UUID) error {

	n, err := s.repo.Get(id)

	if err != nil {
		return balancer.ErrNotFound
	}

	return s.repo.Delete(n.ID)
}
