package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"juno/pkg/api/assignment/dto"
	"juno/pkg/api/assignment/policy"
	"juno/pkg/api/assignment/repo/mem"
	"juno/pkg/api/assignment/service"
	"juno/pkg/api/auth"
	"juno/pkg/api/user"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

func TestCreate(t *testing.T) {
	t.Run("creates assignment", func(t *testing.T) {
		repo := mem.New()
		svc := service.New(repo)
		policy := policy.New()
		handler := New(logrus.New(), policy, svc)

		ownerID := uuid.New()
		nodeID := uuid.New()
		offset := 0
		length := 10

		body, err := json.Marshal(dto.CreateAssignmentRequest{
			OwnerID: ownerID.String(),
			NodeID:  nodeID.String(),
			Offset:  offset,
			Length:  length,
		})

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		w := httptest.NewRecorder()

		tc, _ := gin.CreateTestContext(w)

		tc.Request = httptest.NewRequest("POST", "/assignments", bytes.NewBuffer(body)).
			WithContext(auth.WithUser(context.Background(), &user.User{ID: ownerID}))

		handler.Create(tc)

		if w.Code != 200 {
			t.Fatalf("expected status code 200, got %v", w.Code)
		}

		var resp dto.CreateAssignmentResponse

		if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if resp.Assignment.OwnerID != ownerID.String() {
			t.Errorf("expected assignment OwnerID %v, got %v", ownerID, resp.Assignment.OwnerID)
		}

		if resp.Assignment.NodeID != nodeID.String() {
			t.Errorf("expected assignment NodeID %v, got %v", nodeID, resp.Assignment.NodeID)
		}

		if resp.Assignment.Offset != offset {
			t.Errorf("expected assignment Offset %v, got %v", offset, resp.Assignment.Offset)
		}

		if resp.Assignment.Length != length {
			t.Errorf("expected assignment Length %v, got %v", length, resp.Assignment.Length)
		}
	})

	t.Run("panics if user is not in context", func(t *testing.T) {
		repo := mem.New()
		svc := service.New(repo)
		policy := policy.New()
		handler := New(logrus.New(), policy, svc)

		ownerID := uuid.New()
		nodeID := uuid.New()
		offset := 0
		length := 10

		body, _ := json.Marshal(dto.CreateAssignmentRequest{
			OwnerID: ownerID.String(),
			NodeID:  nodeID.String(),
			Offset:  offset,
			Length:  length,
		})

		w := httptest.NewRecorder()

		tc, _ := gin.CreateTestContext(w)

		tc.Request = httptest.NewRequest("POST", "/assignments", bytes.NewBuffer(body))

		defer func() {
			if r := recover(); r == nil {
				t.Errorf("expected panic, got nil")
			}
		}()

		handler.Create(tc)
	})
}

func TestGet(t *testing.T) {
	t.Run("returns assignment by ID", func(t *testing.T) {
		repo := mem.New()
		svc := service.New(repo)
		policy := policy.New()
		handler := New(logrus.New(), policy, svc)

		ownerID := uuid.New()
		nodeID := uuid.New()
		offset := 0
		length := 10

		assignment, err := svc.Create(ownerID, nodeID, offset, length)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		w := httptest.NewRecorder()

		tc, _ := gin.CreateTestContext(w)

		tc.Params = gin.Params{{Key: "id", Value: assignment.ID.String()}}

		tc.Request = httptest.NewRequest("GET", "/assignments/"+assignment.ID.String(), nil).
			WithContext(auth.WithUser(context.Background(), &user.User{ID: ownerID}))

		handler.Get(tc)

		if w.Code != 200 {
			t.Fatalf("expected status code 200, got %v", w.Code)
		}

		var resp dto.GetAssignmentResponse

		if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if resp.Assignment.ID != assignment.ID.String() {
			t.Errorf("expected assignment ID %v, got %v", assignment.ID, resp.Assignment.ID)
		}

		if resp.Assignment.OwnerID != ownerID.String() {
			t.Errorf("expected assignment OwnerID %v, got %v", ownerID, resp.Assignment.OwnerID)
		}

		if resp.Assignment.NodeID != nodeID.String() {
			t.Errorf("expected assignment NodeID %v, got %v", nodeID, resp.Assignment.NodeID)
		}

		if resp.Assignment.Offset != offset {
			t.Errorf("expected assignment Offset %v, got %v", offset, resp.Assignment.Offset)
		}

		if resp.Assignment.Length != length {
			t.Errorf("expected assignment Length %v, got %v", length, resp.Assignment.Length)
		}
	})

	t.Run("returns error if assignment not found", func(t *testing.T) {
		repo := mem.New()
		svc := service.New(repo)
		policy := policy.New()
		handler := New(logrus.New(), policy, svc)

		w := httptest.NewRecorder()

		tc, _ := gin.CreateTestContext(w)

		tc.Params = gin.Params{{Key: "id", Value: uuid.New().String()}}

		tc.Request = httptest.NewRequest("GET", "/assignments/"+uuid.New().String(), nil).
			WithContext(auth.WithUser(context.Background(), &user.User{ID: uuid.New()}))

		handler.Get(tc)

		if w.Code != 404 {
			t.Fatalf("expected status code 404, got %v", w.Code)
		}
	})

	t.Run("forbidden if user is not owner", func(t *testing.T) {
		repo := mem.New()
		svc := service.New(repo)
		policy := policy.New()
		handler := New(logrus.New(), policy, svc)

		ownerID := uuid.New()
		nodeID := uuid.New()
		offset := 0
		length := 10

		assignment, err := svc.Create(ownerID, nodeID, offset, length)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		w := httptest.NewRecorder()

		tc, _ := gin.CreateTestContext(w)

		tc.Params = gin.Params{{Key: "id", Value: assignment.ID.String()}}

		tc.Request = httptest.NewRequest("GET", "/assignments/"+assignment.ID.String(), nil).
			WithContext(auth.WithUser(context.Background(), &user.User{ID: uuid.New()}))

		handler.Get(tc)

		if w.Code != 403 {
			t.Fatalf("expected status code 401, got %v", w.Code)
		}
	})
}

func TestListByNodeID(t *testing.T) {
	t.Run("returns assignments for entity ID", func(t *testing.T) {
		repo := mem.New()
		svc := service.New(repo)
		policy := policy.New()
		handler := New(logrus.New(), policy, svc)

		ownerID := uuid.New()
		nodeID := uuid.New()
		offset := 0
		length := 10

		assignments := []*dto.Assignment{
			{
				OwnerID: ownerID.String(),
				NodeID:  nodeID.String(),
				Offset:  offset,
				Length:  length,
			},
		}

		for _, a := range assignments {
			_, err := svc.Create(uuid.MustParse(a.OwnerID), uuid.MustParse(a.NodeID), a.Offset, a.Length)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		}

		w := httptest.NewRecorder()

		tc, _ := gin.CreateTestContext(w)

		tc.Params = gin.Params{{Key: "id", Value: nodeID.String()}}

		tc.Request = httptest.NewRequest("GET", "/entities/"+nodeID.String()+"/assignments", nil).
			WithContext(auth.WithUser(context.Background(), &user.User{ID: ownerID}))

		handler.ListByNodeID(tc)

		if w.Code != 200 {
			t.Fatalf("expected status code 200, got %v", w.Code)
		}

		var resp dto.ListAssignmentsResponse

		if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(resp.Assignments) != len(assignments) {
			t.Errorf("expected %v assignments, got %v", len(assignments), len(resp.Assignments))
		}

		for i, a := range assignments {

			if resp.Assignments[i].OwnerID != a.OwnerID {
				t.Errorf("expected assignment OwnerID %v, got %v", a.OwnerID, resp.Assignments[i].OwnerID)
			}

			if resp.Assignments[i].NodeID != a.NodeID {
				t.Errorf("expected assignment NodeID %v, got %v", a.NodeID, resp.Assignments[i].NodeID)
			}

			if resp.Assignments[i].Offset != a.Offset {
				t.Errorf("expected assignment Offset %v, got %v", a.Offset, resp.Assignments[i].Offset)
			}

			if resp.Assignments[i].Length != a.Length {
				t.Errorf("expected assignment Length %v, got %v", a.Length, resp.Assignments[i].Length)
			}
		}
	})
}

func TestUpdate(t *testing.T) {
	t.Run("updates assignment", func(t *testing.T) {
		repo := mem.New()
		svc := service.New(repo)
		policy := policy.New()
		handler := New(logrus.New(), policy, svc)

		ownerID := uuid.New()
		nodeID := uuid.New()
		offset := 0
		length := 10

		assignment, err := svc.Create(ownerID, nodeID, offset, length)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		newOffset := 5
		newLength := 15

		body, err := json.Marshal(dto.UpdateAssignmentRequest{
			Offset: newOffset,
			Length: newLength,
		})

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		w := httptest.NewRecorder()

		tc, _ := gin.CreateTestContext(w)

		tc.Params = gin.Params{{Key: "id", Value: assignment.ID.String()}}

		tc.Request = httptest.NewRequest("PUT", "/assignments/"+assignment.ID.String(), bytes.NewBuffer(body)).
			WithContext(auth.WithUser(context.Background(), &user.User{ID: ownerID}))

		handler.Update(tc)

		if w.Code != 200 {
			t.Fatalf("expected status code 200, got %v", w.Code)
		}

		var resp dto.UpdateAssignmentResponse

		if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if resp.Assignment.ID != assignment.ID.String() {
			t.Errorf("expected assignment ID %v, got %v", assignment.ID, resp.Assignment.ID)
		}

		if resp.Assignment.OwnerID != ownerID.String() {
			t.Errorf("expected assignment OwnerID %v, got %v", ownerID, resp.Assignment.OwnerID)
		}

		if resp.Assignment.NodeID != nodeID.String() {
			t.Errorf("expected assignment NodeID %v, got %v", nodeID, resp.Assignment.NodeID)
		}

		if resp.Assignment.Offset != newOffset {
			t.Errorf("expected assignment Offset %v, got %v", newOffset, resp.Assignment.Offset)
		}

		if resp.Assignment.Length != newLength {
			t.Errorf("expected assignment Length %v, got %v", newLength, resp.Assignment.Length)
		}
	})
}

func TestDelete(t *testing.T) {
	t.Run("deletes assignment", func(t *testing.T) {
		repo := mem.New()
		svc := service.New(repo)
		policy := policy.New()
		handler := New(logrus.New(), policy, svc)

		ownerID := uuid.New()
		nodeID := uuid.New()
		offset := 0
		length := 10

		assignment, err := svc.Create(ownerID, nodeID, offset, length)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		w := httptest.NewRecorder()

		tc, _ := gin.CreateTestContext(w)

		tc.Params = gin.Params{{Key: "id", Value: assignment.ID.String()}}

		tc.Request = httptest.NewRequest("DELETE", "/assignments/"+assignment.ID.String(), nil).
			WithContext(auth.WithUser(context.Background(), &user.User{ID: ownerID}))

		handler.Delete(tc)

		if w.Code != 200 {
			t.Fatalf("expected status code 200, got %v", w.Code)
		}

		var resp dto.DeleteAssignmentResponse

		if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if resp.Status != dto.SUCCESS {
			t.Errorf("expected status %v, got %v", dto.SUCCESS, resp.Status)
		}

		// check if assignment was deleted
		_, err = svc.Get(assignment.ID)

		if err == nil {
			t.Fatalf("expected error, got nil")
		}
	})
}
