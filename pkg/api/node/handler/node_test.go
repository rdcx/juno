package handler

import (
	"bytes"
	"context"
	"juno/pkg/api/auth"
	"juno/pkg/api/node"
	"juno/pkg/api/node/dto"
	"juno/pkg/api/node/policy"
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

func testNodeMatches(t *testing.T, id, ownerID uuid.UUID, address string, shardAssignments [][2]int, n *node.Node) bool {
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

	if len(n.ShardAssignments) != len(shardAssignments) {
		t.Errorf("Expected %d shard assignments, got %d", len(shardAssignments), len(n.ShardAssignments))
		return false
	}

	for i, s := range shardAssignments {
		if n.ShardAssignments[i][0] != s[0] || n.ShardAssignments[i][1] != s[1] {
			t.Errorf("Expected shard assignment %v, got %v", s, n.ShardAssignments[i])
			return false
		}
	}

	return true
}

func TestAllShardsNodes(t *testing.T) {
	repo := mem.New()
	svc := service.New(repo)
	handler := New(logrus.New(), policy.New(), svc)

	n1 := &node.Node{
		ID:               uuid.New(),
		OwnerID:          uuid.New(),
		Address:          "http://example.com",
		ShardAssignments: [][2]int{{0, 1000}, {1000, 1000}},
	}

	err := repo.Create(n1)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	n2 := &node.Node{
		ID:               uuid.New(),
		OwnerID:          uuid.New(),
		Address:          "node.com:8080",
		ShardAssignments: [][2]int{{0, 1000}, {1000, 1000}},
	}

	err = repo.Create(n2)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	w := httptest.NewRecorder()
	tc, _ := gin.CreateTestContext(w)

	tc.Request = httptest.NewRequest("GET", "/nodes", nil)

	handler.AllShardsNodes(tc)

	var res dto.AllShardsNodesResponse

	err = json.Unmarshal(w.Body.Bytes(), &res)

	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	if res.Status != "success" {
		t.Errorf("Expected success, got %s", res.Status)
	}

	if len(res.Shards) != 2000 {
		t.Errorf("Expected 2000 shards, got %d", len(res.Shards))
	}

	for i := 0; i < 2000; i++ {
		if res.Shards[i][0] != n1.Address {
			t.Errorf("Expected address %s for shard %d, got %s", n1.Address, i, res.Shards[i][0])
		}
	}
}

func TestGet(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := mem.New()
		svc := service.New(repo)
		policy := policy.New()
		handler := New(logrus.New(), policy, svc)

		u := &user.User{
			ID: uuid.New(),
		}

		n := &node.Node{
			ID:      uuid.New(),
			OwnerID: u.ID,
			Address: "example.com:8080",
			ShardAssignments: [][2]int{
				{0, 100},
				{100, 200},
			},
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

		tc.Request = httptest.NewRequest("GET", "/nodes/"+n.ID.String(), nil).WithContext(
			auth.WithUser(context.Background(), u),
		)

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

		if !testNodeMatches(t, n.ID, n.OwnerID, n.Address, n.ShardAssignments, resN) {
			t.Errorf("Node does not match")
		}
	})

	t.Run("error", func(t *testing.T) {
		repo := mem.New()
		svc := service.New(repo)
		policy := policy.New()
		handler := New(logrus.New(), policy, svc)

		u := &user.User{
			ID: uuid.New(),
		}

		n := &node.Node{
			ID:      uuid.New(),
			OwnerID: uuid.New(),
			Address: "example.com:8080",
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

		tc.Request = httptest.NewRequest("GET", "/nodes/"+n.ID.String(), nil).WithContext(
			auth.WithUser(context.Background(), u),
		)

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

func TestList(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := mem.New()
		svc := service.New(repo)
		policy := policy.New()
		handler := New(logrus.New(), policy, svc)

		u := &user.User{
			ID: uuid.New(),
		}

		n1 := &node.Node{
			ID:      uuid.New(),
			OwnerID: u.ID,
			Address: "example.com:8080",
		}

		n2 := &node.Node{
			ID:      uuid.New(),
			OwnerID: uuid.New(),
			Address: "example.com:8081",
			ShardAssignments: [][2]int{
				{0, 100},
				{100, 200},
			},
		}

		err := repo.Create(n1)
		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}

		err = repo.Create(n2)
		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}

		w := httptest.NewRecorder()

		tc, _ := gin.CreateTestContext(w)

		tc.Request = httptest.NewRequest("GET", "/nodes", nil).WithContext(
			auth.WithUser(context.Background(), u),
		)

		handler.List(tc)

		if w.Code != 200 {
			t.Fatalf("Expected status code 200, got %d", w.Code)
		}

		var res dto.ListNodesResponse

		err = json.Unmarshal(w.Body.Bytes(), &res)

		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}

		if res.Status != "success" {
			t.Errorf("Expected success, got %s", res.Status)
		}

		if len(res.Nodes) != 1 {
			t.Fatalf("Expected 1 nodes, got %d", len(res.Nodes))
		}

		nodeCheck, err := res.Nodes[0].ToDomain()

		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}

		if !testNodeMatches(t, n1.ID, n1.OwnerID, n1.Address, n1.ShardAssignments, nodeCheck) {
			t.Errorf("Node 1 does not match")
		}
	})
}

func TestCreate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := mem.New()
		svc := service.New(repo)
		policy := policy.New()
		handler := New(logrus.New(), policy, svc)

		u := &user.User{
			ID: uuid.New(),
		}

		addr := "example.com:8080"
		shardAssignments := [][2]int{
			{0, 100},
			{100, 200},
		}

		// Create a new recorder and test context
		w := httptest.NewRecorder()
		tc, _ := gin.CreateTestContext(w)

		// Marshal the request body
		body, err := json.Marshal(dto.CreateNodeRequest{
			Address:          addr,
			ShardAssignments: shardAssignments,
		})

		if err != nil {
			t.Fatalf("Unexpected error: %s", err)
		}

		// Set the request body correctly during request creation
		tc.Request = httptest.NewRequest("POST", "/nodes", bytes.NewReader(body)).WithContext(
			auth.WithUser(context.Background(), u),
		)

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
		if !testNodeMatches(t, resN.ID, u.ID, addr, shardAssignments, resN) {
			t.Errorf("Node does not match")
		}
	})

	t.Run("error", func(t *testing.T) {
		repo := mem.New()
		svc := service.New(repo)
		policy := policy.New()
		handler := New(logrus.New(), policy, svc)

		u := &user.User{
			ID: uuid.New(),
		}

		n := &node.Node{
			ID:      uuid.New(),
			OwnerID: u.ID,
			Address: "example.com:8080",
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
		})
		if err != nil {
			t.Fatalf("Unexpected error: %s", err)
		}

		// Set the request body correctly during request creation
		tc.Request = httptest.NewRequest("POST", "/nodes", bytes.NewReader(body)).
			WithContext(auth.WithUser(context.Background(), u))

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
		policy := policy.New()
		handler := New(logrus.New(), policy, svc)

		u := &user.User{
			ID: uuid.New(),
		}

		n := &node.Node{
			ID:      uuid.New(),
			OwnerID: u.ID,
			Address: "http://example.com",
			ShardAssignments: [][2]int{
				{0, 100},
				{100, 200},
			},
		}

		err := repo.Create(n)
		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}

		addr := "new.example.com:2000"
		offset := 1000
		shards := 2000
		shardAssignments := [][2]int{{offset, shards}}

		// Create a new recorder and test context
		w := httptest.NewRecorder()
		tc, _ := gin.CreateTestContext(w)

		// Set the user in the context
		tc.Set("user", u)

		// Marshal the request body
		body, err := json.Marshal(dto.UpdateNodeRequest{
			Address:          addr,
			ShardAssignments: shardAssignments,
		})
		if err != nil {
			t.Fatalf("Unexpected error: %s", err)
		}

		// Set the request body correctly during request creation
		tc.Request = httptest.NewRequest("PUT", "/nodes/"+n.ID.String(), bytes.NewReader(body)).
			WithContext(auth.WithUser(context.Background(), u))

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

		if res.Node == nil {
			t.Fatal("Expected node, got nil")
		}

		// Convert response node to domain
		resN, err := res.Node.ToDomain()

		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}

		// Check if node matches
		if !testNodeMatches(t, resN.ID, u.ID, addr, shardAssignments, resN) {
			t.Errorf("Node does not match")
		}
	})

	t.Run("error", func(t *testing.T) {
		repo := mem.New()
		svc := service.New(repo)
		policy := policy.New()
		handler := New(logrus.New(), policy, svc)

		u := &user.User{
			ID: uuid.New(),
		}

		n := &node.Node{
			ID:      uuid.New(),
			OwnerID: u.ID,
			Address: "valid.com:8000",
		}

		err := repo.Create(n)
		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}

		addr := "invalid"
		offset := 0
		shards := 100

		// Create a new recorder and test context
		w := httptest.NewRecorder()
		tc, _ := gin.CreateTestContext(w)

		// Set the user in the context
		tc.Set("user", u)

		// Marshal the request body
		body, err := json.Marshal(dto.UpdateNodeRequest{
			Address:          addr,
			ShardAssignments: [][2]int{{offset, shards}},
		})
		if err != nil {
			t.Fatalf("Unexpected error: %s", err)
		}

		// Set the request body correctly during request creation
		tc.Request = httptest.NewRequest("PUT", "/nodes/"+n.ID.String(), bytes.NewReader(body)).
			WithContext(auth.WithUser(context.Background(), u))

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
		policy := policy.New()
		handler := New(logrus.New(), policy, svc)

		u := &user.User{
			ID: uuid.New(),
		}

		n := &node.Node{
			ID:      uuid.New(),
			OwnerID: u.ID,
			Address: "http://example.com",
		}

		err := repo.Create(n)
		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}

		// Create a new recorder and test context
		w := httptest.NewRecorder()
		tc, _ := gin.CreateTestContext(w)

		// Set the user in the context
		tc.Request = httptest.NewRequest("DELETE", "/nodes/"+n.ID.String(), nil).WithContext(
			auth.WithUser(context.Background(), u),
		)

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
		policy := policy.New()

		handler := New(logrus.New(), policy, svc)

		u := &user.User{
			ID: uuid.New(),
		}

		n := &node.Node{
			ID: uuid.New(),
		}

		// Create a new recorder and test context
		w := httptest.NewRecorder()
		tc, _ := gin.CreateTestContext(w)

		// Set the user in the context
		tc.Request = httptest.NewRequest("DELETE", "/nodes/"+n.ID.String(), nil).WithContext(
			auth.WithUser(context.Background(), u),
		)

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
