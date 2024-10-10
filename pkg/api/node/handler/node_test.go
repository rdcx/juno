package handler

import (
	"bytes"
	"juno/pkg/api/node"
	"juno/pkg/api/node/dto"
	"juno/pkg/api/node/repo/mem"
	"juno/pkg/api/node/service"
	"juno/pkg/api/user"
	"net/http/httptest"
	"strings"
	"testing"

	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

func testNodeMatches(t *testing.T, id, ownerID uuid.UUID, address string, shards []int, n *node.Node) bool {
	if n.ID != id {
		t.Errorf("Expected ID %s, got %s", id, n.ID)
		return false
	}

	if n.OwnerID != ownerID {
		t.Errorf("Expected OwnerID %s, got %s", ownerID, n.OwnerID)
		return false
	}

	if n.Address != address {
		t.Errorf("Expected Address %s, got %s", address, n.Address)
		return false
	}

	if len(n.Shards) != len(shards) {
		t.Errorf("Expected %d shards, got %d", len(shards), len(n.Shards))
		return false
	}

	for i, shard := range shards {
		if n.Shards[i] != shard {
			t.Errorf("Expected shard %d, got %d", shard, n.Shards[i])
			return false
		}
	}

	return true
}

func TestGet(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := mem.New()
		svc := service.New(repo)
		handler := New(logrus.New(), svc)

		u := &user.User{
			ID: uuid.New(),
		}

		n := &node.Node{
			ID:      uuid.New(),
			OwnerID: u.ID,
			Address: "example.com:8080",
			Shards:  []int{1, 2, 3},
		}

		err := repo.Create(n)
		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}

		w := httptest.NewRecorder()

		tc, _ := gin.CreateTestContext(w)

		tc.Params = gin.Params{
			{Key: "id", Value: n.ID.String()},
		}

		tc.Set("user", u)

		handler.Get(tc)

		if w.Code != 200 {
			t.Fatalf("Expected status code 200, got %d", w.Code)
		}
		var res dto.GetNodeResponse

		err = json.Unmarshal(w.Body.Bytes(), &res)

		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}

		if res.Status != "success" {
			t.Errorf("Expected success, got %s", res.Status)
		}

		if res.Node == nil {
			t.Fatal("Expected node, got nil")
		}

		resN, err := res.Node.ToDomain()

		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}

		if !testNodeMatches(t, n.ID, n.OwnerID, n.Address, n.Shards, resN) {
			t.Errorf("Node does not match")
		}
	})

	t.Run("error", func(t *testing.T) {
		repo := mem.New()
		svc := service.New(repo)
		handler := New(logrus.New(), svc)

		u := &user.User{
			ID: uuid.New(),
		}

		n := &node.Node{
			ID:      uuid.New(),
			OwnerID: uuid.New(),
			Address: "example.com:8080",
			Shards:  []int{1, 2, 3},
		}

		err := repo.Create(n)
		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}

		w := httptest.NewRecorder()

		tc, _ := gin.CreateTestContext(w)

		tc.Params = gin.Params{
			{Key: "id", Value: n.ID.String()},
		}

		tc.Set("user", u)

		handler.Get(tc)

		if w.Code != 404 {
			t.Fatalf("Expected status code 404, got %d", w.Code)
		}

		var res *dto.GetNodeResponse

		err = json.Unmarshal(w.Body.Bytes(), &res)

		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}

		if res.Message != "node not found" {
			t.Errorf("Expected error node not found, got %s", res.Message)
		}
	})
}

func TestCreate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := mem.New()
		svc := service.New(repo)
		handler := New(logrus.New(), svc)

		u := &user.User{
			ID: uuid.New(),
		}

		addr := "example.com:8080"
		shards := []int{1, 2, 3}

		// Create a new recorder and test context
		w := httptest.NewRecorder()
		tc, _ := gin.CreateTestContext(w)

		// Set the user in the context
		tc.Set("user", u)

		// Marshal the request body
		body, err := json.Marshal(dto.CreateNodeRequest{
			Address: addr,
			Shards:  shards,
		})
		if err != nil {
			t.Fatalf("Unexpected error: %s", err)
		}

		// Set the request body correctly during request creation
		tc.Request = httptest.NewRequest("POST", "/nodes", bytes.NewReader(body))

		// Call the handler
		handler.Create(tc)

		// Check the status code
		if w.Code != 201 {
			t.Fatalf("Expected status code 201, got %d", w.Code)
		}

		// Parse the response
		var res dto.CreateNodeResponse
		err = json.Unmarshal(w.Body.Bytes(), &res)
		if err != nil {
			t.Fatalf("Unexpected error: %s", err)
		}

		// Validate the response
		if res.Status != "success" {
			t.Errorf("Expected success, got %s", res.Status)
		}

		if res.Node == nil {
			t.Fatal("Expected node, got nil")
		}

		// Convert response node to domain
		resN, err := res.Node.ToDomain()
		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}

		// Check if node matches
		if !testNodeMatches(t, resN.ID, u.ID, addr, shards, resN) {
			t.Errorf("Node does not match")
		}
	})

	t.Run("error", func(t *testing.T) {
		repo := mem.New()
		svc := service.New(repo)
		handler := New(logrus.New(), svc)

		u := &user.User{
			ID: uuid.New(),
		}

		n := &node.Node{
			ID:      uuid.New(),
			OwnerID: u.ID,
			Address: "example.com:8080",
			Shards:  []int{1, 2, 3},
		}

		// Pre-create the node to simulate a duplicate node creation
		err := repo.Create(n)
		if err != nil {
			t.Fatalf("Unexpected error: %s", err)
		}

		// Create a new recorder and test context
		w := httptest.NewRecorder()
		tc, _ := gin.CreateTestContext(w)

		// Set the user in the context
		tc.Set("user", u)

		// Marshal the request body
		body, err := json.Marshal(dto.CreateNodeRequest{
			Address: n.Address,
			Shards:  n.Shards,
		})
		if err != nil {
			t.Fatalf("Unexpected error: %s", err)
		}

		// Set the request body correctly during request creation
		tc.Request = httptest.NewRequest("POST", "/nodes", bytes.NewReader(body))

		// Call the handler
		handler.Create(tc)

		// Check the status code
		if w.Code != 400 {
			t.Fatalf("Expected status code 400, got %d", w.Code)
		}

		// Parse the response
		var res dto.CreateNodeResponse
		err = json.Unmarshal(w.Body.Bytes(), &res)
		if err != nil {
			t.Fatalf("Unexpected error: %s", err)
		}

		// Validate the error message
		if res.Message != "address already exists" {
			t.Errorf("Expected error message 'address already exists', got %s", res.Message)
		}
	})
}

func TestUpdate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := mem.New()
		svc := service.New(repo)
		handler := New(logrus.New(), svc)

		u := &user.User{
			ID: uuid.New(),
		}

		n := &node.Node{
			ID:      uuid.New(),
			OwnerID: u.ID,
			Address: "http://example.com",
			Shards:  []int{1, 2, 3},
		}

		err := repo.Create(n)
		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}

		addr := "new.example.com:2000"
		shards := []int{4, 5, 6}

		// Create a new recorder and test context
		w := httptest.NewRecorder()
		tc, _ := gin.CreateTestContext(w)

		// Set the user in the context
		tc.Set("user", u)

		// Marshal the request body
		body, err := json.Marshal(dto.UpdateNodeRequest{
			Address: addr,
			Shards:  shards,
		})
		if err != nil {
			t.Fatalf("Unexpected error: %s", err)
		}

		// Set the request body correctly during request creation
		tc.Request = httptest.NewRequest("PUT", "/nodes/"+n.ID.String(), bytes.NewReader(body))

		// Set the path parameter
		tc.Params = gin.Params{
			{Key: "id", Value: n.ID.String()},
		}

		// Call the handler
		handler.Update(tc)

		// Check the status code
		if w.Code != 200 {
			t.Fatalf("Expected status code 200, got %d", w.Code)
		}

		// Parse the response
		var res dto.UpdateNodeResponse
		err = json.Unmarshal(w.Body.Bytes(), &res)
		if err != nil {
			t.Fatalf("Unexpected error: %s", err)
		}

		// Validate the response
		if res.Status != "success" {
			t.Errorf("Expected success, got %s", res.Status)
		}
	})

	t.Run("error", func(t *testing.T) {
		repo := mem.New()
		svc := service.New(repo)
		handler := New(logrus.New(), svc)

		u := &user.User{
			ID: uuid.New(),
		}

		n := &node.Node{
			ID:      uuid.New(),
			OwnerID: u.ID,
			Address: "valid.com:8000",
			Shards:  []int{1, 2, 3},
		}

		err := repo.Create(n)
		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}

		addr := "invalid"
		shards := []int{4, 5, 6}

		// Create a new recorder and test context
		w := httptest.NewRecorder()
		tc, _ := gin.CreateTestContext(w)

		// Set the user in the context
		tc.Set("user", u)

		// Marshal the request body
		body, err := json.Marshal(dto.UpdateNodeRequest{
			Address: addr,
			Shards:  shards,
		})
		if err != nil {
			t.Fatalf("Unexpected error: %s", err)
		}

		// Set the request body correctly during request creation
		tc.Request = httptest.NewRequest("PUT", "/nodes/"+n.ID.String(), bytes.NewReader(body))

		// Set the path parameter
		tc.Params = gin.Params{
			{Key: "id", Value: n.ID.String()},
		}

		// Call the handler
		handler.Update(tc)

		// Check the status code
		if w.Code != 400 {
			t.Fatalf("Expected status code 400, got %d", w.Code)
		}

		// Parse the response
		var res dto.UpdateNodeResponse
		err = json.Unmarshal(w.Body.Bytes(), &res)
		if err != nil {
			t.Fatalf("Unexpected error: %s", err)
		}

		// Validate the error message
		if !strings.Contains(res.Message, "invalid address") {
			t.Errorf("Expected validation error message 'invalid address', got %s", res.Message)
		}
	})
}

func TestDelete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := mem.New()
		svc := service.New(repo)
		handler := New(logrus.New(), svc)

		u := &user.User{
			ID: uuid.New(),
		}

		n := &node.Node{
			ID:      uuid.New(),
			OwnerID: u.ID,
			Address: "http://example.com",
			Shards:  []int{1, 2, 3},
		}

		err := repo.Create(n)
		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}

		// Create a new recorder and test context
		w := httptest.NewRecorder()
		tc, _ := gin.CreateTestContext(w)

		// Set the user in the context
		tc.Set("user", u)

		// Set the path parameter
		tc.Params = gin.Params{
			{Key: "id", Value: n.ID.String()},
		}

		// Call the handler
		handler.Delete(tc)

		// Check the status code
		if w.Code != 200 {
			t.Fatalf("Expected status code 200, got %d", w.Code)
		}

		// Parse the response
		var res dto.DeleteNodeResponse
		err = json.Unmarshal(w.Body.Bytes(), &res)
		if err != nil {
			t.Fatalf("Unexpected error: %s", err)
		}

		// Validate the response
		if res.Status != "success" {
			t.Errorf("Expected success, got %s", res.Status)
		}
	})

	t.Run("error", func(t *testing.T) {
		repo := mem.New()
		svc := service.New(repo)
		handler := New(logrus.New(), svc)

		u := &user.User{
			ID: uuid.New(),
		}

		n := &node.Node{
			ID:     uuid.New(),
			Shards: []int{1, 2, 3},
		}

		// Create a new recorder and test context
		w := httptest.NewRecorder()
		tc, _ := gin.CreateTestContext(w)

		// Set the user in the context
		tc.Set("user", u)

		// Set the path parameter
		tc.Params = gin.Params{
			{Key: "id", Value: n.ID.String()},
		}

		// Call the handler
		handler.Delete(tc)

		// Check the status code
		if w.Code != 404 {
			t.Fatalf("Expected status code 404, got %d", w.Code)
		}

		// Parse the response
		var res dto.DeleteNodeResponse
		err := json.Unmarshal(w.Body.Bytes(), &res)
		if err != nil {
			t.Fatalf("Unexpected error: %s", err)
		}

		// Validate the error message
		if res.Message != "node not found" {
			t.Errorf("Expected error message 'node not found', got %s", res.Message)
		}
	})
}
