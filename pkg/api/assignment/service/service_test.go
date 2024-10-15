package service

import (
	"juno/pkg/api/assignment/repo/mem"
	"testing"

	"github.com/google/uuid"
)

func TestCreate(t *testing.T) {
	t.Run("creates assignment", func(t *testing.T) {
		repo := mem.New()
		svc := New(repo)
		ownerID := uuid.New()
		nodeID := uuid.New()
		offset := 0
		length := 10

		assignment, err := svc.Create(ownerID, nodeID, offset, length)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if assignment.OwnerID != ownerID {
			t.Errorf("expected assignment OwnerID %v, got %v", ownerID, assignment.OwnerID)
		}

		if assignment.NodeID != nodeID {
			t.Errorf("expected assignment NodeID %v, got %v", nodeID, assignment.NodeID)
		}

		if assignment.Offset != offset {
			t.Errorf("expected assignment Offset %v, got %v", offset, assignment.Offset)
		}

		if assignment.Length != length {
			t.Errorf("expected assignment Length %v, got %v", length, assignment.Length)
		}

	})
}
