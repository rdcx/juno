package service

import (
	"juno/pkg/api/assignment"

	"github.com/google/uuid"
)

type Service struct {
	assignmentRepo assignment.Repository
}

func New(assignmentRepo assignment.Repository) *Service {
	return &Service{
		assignmentRepo: assignmentRepo,
	}
}

func (s *Service) Create(ownerID, nodeID uuid.UUID, offset, length int) (*assignment.Assignment, error) {

	assignment := &assignment.Assignment{
		ID:      uuid.New(),
		OwnerID: ownerID,
		NodeID:  nodeID,
		Offset:  offset,
		Length:  length,
	}

	err := s.assignmentRepo.Create(assignment)

	if err != nil {
		return nil, err
	}

	return assignment, nil
}

func (s *Service) Get(id uuid.UUID) (*assignment.Assignment, error) {
	return s.assignmentRepo.Get(id)
}

func (s *Service) ListByNodeID(nodeID uuid.UUID) ([]*assignment.Assignment, error) {
	return s.assignmentRepo.ListByNodeID(nodeID)
}

func (s *Service) Update(id uuid.UUID, offset, length int) (*assignment.Assignment, error) {
	assignment, err := s.assignmentRepo.Get(id)

	if err != nil {
		return nil, err
	}

	assignment.Offset = offset
	assignment.Length = length

	err = s.assignmentRepo.Update(assignment)

	if err != nil {
		return nil, err
	}

	return assignment, nil
}

func (s *Service) Delete(id uuid.UUID) error {
	return s.assignmentRepo.Delete(id)
}
