package service

import (
	"juno/pkg/api/balancer"
	"juno/pkg/shard"
	"juno/pkg/util"
	"strings"

	"github.com/google/uuid"
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

func validateShardAssignments(shardAssignments [][2]int) error {
	for _, s := range shardAssignments {
		if s[0] < 0 || s[1] < 0 {
			return balancer.ErrInvalidShards
		}

		if s[0]+s[1] > shard.SHARDS {
			return balancer.ErrInvalidShards
		}
	}

	return nil
}

func (s *Service) AllShardsBalancers() (map[int][]*balancer.Balancer, error) {
	balancers, err := s.repo.All()

	if err != nil {
		return nil, balancer.ErrInternal
	}

	shardsBalancers := make(map[int][]*balancer.Balancer)

	for _, n := range balancers {
		for _, s := range n.ShardAssignments {
			for i := s[0]; i < s[0]+s[1]; i++ {
				shardsBalancers[i] = append(shardsBalancers[i], n)
			}
		}
	}

	return shardsBalancers, nil
}

func (s *Service) Create(ownerID uuid.UUID, addr string, shardAssignments [][2]int) (*balancer.Balancer, error) {
	if found, _ := s.repo.FirstWhereAddress(addr); found != nil {
		return nil, balancer.ErrAddressExists
	}

	errs := []error{}

	if err := validateShardAssignments(shardAssignments); err != nil {
		errs = append(errs, err)
	}

	if err := validateAddress(addr); err != nil {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return nil, util.ValidationErrs(errs)
	}

	n := &balancer.Balancer{
		ID:               uuid.New(),
		OwnerID:          ownerID,
		Address:          addr,
		ShardAssignments: shardAssignments,
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
	n.ShardAssignments = dirty.ShardAssignments

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
