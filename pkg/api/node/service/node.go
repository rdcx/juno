package service

import (
	"juno/pkg/api/node"
	"juno/pkg/shard"
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

func validateShardAssignments(shardAssignments [][2]int) error {
	for _, s := range shardAssignments {
		if s[0] < 0 || s[1] < 0 {
			return node.ErrInvalidShards
		}

		if s[0]+s[1] > shard.SHARDS {
			return node.ErrInvalidShards
		}
	}

	return nil
}

func (s *Service) AllShardsNodes() (map[int][]*node.Node, error) {
	nodes, err := s.repo.All()

	if err != nil {
		return nil, node.ErrInternal
	}

	shardsNodes := make(map[int][]*node.Node)

	for _, n := range nodes {
		for _, s := range n.ShardAssignments {
			for i := s[0]; i < s[0]+s[1]; i++ {
				shardsNodes[i] = append(shardsNodes[i], n)
			}
		}
	}

	return shardsNodes, nil
}

func (s *Service) Create(ownerID uuid.UUID, addr string, shardAssignments [][2]int) (*node.Node, error) {
	if found, _ := s.repo.FirstWhereAddress(addr); found != nil {
		return nil, node.ErrAddressExists
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

	n := &node.Node{
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
		return node.ErrNotFound
	}

	return s.repo.Delete(n.ID)
}
