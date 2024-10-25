package service

import (
	"juno/pkg/api/ranag"
	"juno/pkg/shard"
	"juno/pkg/util"
	"strings"

	"github.com/google/uuid"
)

type Service struct {
	repo ranag.Repository
}

func New(repo ranag.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Get(id uuid.UUID) (*ranag.Ranag, error) {
	n, err := s.repo.Get(id)

	if err != nil {
		return nil, ranag.ErrNotFound
	}
	return n, nil
}

func (s *Service) ListByOwnerID(ownerID uuid.UUID) ([]*ranag.Ranag, error) {
	ranags, err := s.repo.ListByOwnerID(ownerID)

	if err != nil {
		return nil, ranag.ErrNotFound
	}

	return ranags, nil
}

func validateAddress(addr string) error {
	addSplit := strings.Split(addr, ":")

	if len(addSplit) != 2 {
		return ranag.ErrInvalidAddress
	}

	host := addSplit[0]
	port := addSplit[1]

	if !util.IsValidHostname(host) || !util.IsValidPort(port) {
		return ranag.ErrInvalidAddress
	}

	return nil
}

func validateShardAssignments(shardAssignments [][2]int) error {
	for _, s := range shardAssignments {
		if s[0] < 0 || s[1] < 0 {
			return ranag.ErrInvalidShards
		}

		if s[0]+s[1] > shard.SHARDS {
			return ranag.ErrInvalidShards
		}
	}

	return nil
}

func (s *Service) AllShardsRanags() (map[int][]*ranag.Ranag, error) {
	ranags, err := s.repo.All()

	if err != nil {
		return nil, ranag.ErrInternal
	}

	shardsRanags := make(map[int][]*ranag.Ranag)

	for _, n := range ranags {
		for _, s := range n.ShardAssignments {
			for i := s[0]; i < s[0]+s[1]; i++ {
				shardsRanags[i] = append(shardsRanags[i], n)
			}
		}
	}

	return shardsRanags, nil
}

func (s *Service) Create(ownerID uuid.UUID, addr string, shardAssignments [][2]int) (*ranag.Ranag, error) {
	if found, _ := s.repo.FirstWhereAddress(addr); found != nil {
		return nil, ranag.ErrAddressExists
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

	n := &ranag.Ranag{
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

func (s *Service) Update(id uuid.UUID, dirty *ranag.Ranag) (*ranag.Ranag, error) {
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
		return nil, ranag.ErrAddressExists
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
		return ranag.ErrNotFound
	}

	return s.repo.Delete(n.ID)
}

func (s *Service) GroupByRange() (map[[2]int][]*ranag.Ranag, error) {
	grouped := make(map[[2]int][]*ranag.Ranag)

	ranags, err := s.repo.All()

	if err != nil {
		return nil, ranag.ErrInternal
	}

	for _, n := range ranags {
		for _, s := range n.ShardAssignments {
			grouped[s] = append(grouped[s], n)
		}
	}

	return grouped, nil
}
