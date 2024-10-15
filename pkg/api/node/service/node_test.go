package service

import (
	"juno/pkg/api/node"
	"juno/pkg/api/node/repo/mem"
	"juno/pkg/api/user"
	"strings"
	"testing"

	"github.com/google/uuid"
)

func testNodeMatches(t *testing.T, id, ownerID uuid.UUID, address string, n *node.Node) bool {
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

	return true
}

func TestCreate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := mem.New()
		svc := New(repo)

		u := &user.User{
			ID: uuid.New(),
		}

		addr := "example.com:7000"

		n, err := svc.Create(u.ID, addr, [][2]int{{0, 1}, {1, 2}})

		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}

		node, err := repo.Get(n.ID)

		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}

		if !testNodeMatches(t, n.ID, n.OwnerID, addr, node) {
			t.Errorf("Node does not match")
		}
	})

	t.Run("validation", func(t *testing.T) {
		repo := mem.New()
		svc := New(repo)

		u := &user.User{
			ID: uuid.New(),
		}

		addr := "bad address"

		n, err := svc.Create(u.ID, addr, [][2]int{{0, 1}, {1, 2}})

		if err == nil {
			t.Fatal("Expected an error")
		}

		if !strings.Contains(err.Error(), "invalid address") {
			t.Errorf("Expected ErrInvalidAddress, got %v", err)
		}

		if n != nil {
			t.Errorf("Expected nil, got %v", n)
		}
	})

	t.Run("error address already exists", func(t *testing.T) {
		repo := mem.New()
		svc := New(repo)

		n := &node.Node{
			ID:      uuid.New(),
			OwnerID: uuid.New(),
			Address: "http://example.com",
		}

		u := &user.User{
			ID: n.OwnerID,
		}

		err := repo.Create(n)

		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}

		addr := "http://example.com"

		n, err = svc.Create(u.ID, addr, [][2]int{{0, 1}, {1, 2}})

		if err != node.ErrAddressExists {
			t.Errorf("Expected error, got nil")
		}

		if n != nil {
			t.Errorf("Expected nil, got %v", n)
		}
	})
}

func TestGet(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := mem.New()
		svc := New(repo)

		n := &node.Node{
			ID:      uuid.New(),
			OwnerID: uuid.New(),
			Address: "example.com:8000",
		}

		err := repo.Create(n)
		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}

		node, err := svc.Get(n.ID)
		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}

		if !testNodeMatches(t, n.ID, n.OwnerID, n.Address, node) {
			t.Errorf("Node does not match")
		}
	})

	t.Run("error not found", func(t *testing.T) {
		repo := mem.New()
		svc := New(repo)

		n := &node.Node{
			ID:      uuid.New(),
			OwnerID: uuid.New(),
		}

		_, err := svc.Get(n.ID)
		if err != node.ErrNotFound {
			t.Errorf("Expected error, got nil")
		}
	})
}

func TestListByOwnerID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := mem.New()
		svc := New(repo)

		ownerID := uuid.New()

		n1 := &node.Node{
			ID:      uuid.New(),
			OwnerID: ownerID,
			Address: "example.com:8000",
		}

		err := repo.Create(n1)
		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}

		n2 := &node.Node{
			ID:      uuid.New(),
			OwnerID: uuid.New(),
			Address: "example2.com:8080",
		}

		err = repo.Create(n2)

		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}

		nodes, err := svc.ListByOwnerID(ownerID)
		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}

		if len(nodes) != 1 {
			t.Errorf("Expected 1 nodes, got %d", len(nodes))
		}

		if !testNodeMatches(t, n1.ID, n1.OwnerID, n1.Address, nodes[0]) {
			t.Errorf("Node does not match")
		}
	})
}

func TestUpdate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := mem.New()
		svc := New(repo)

		n := &node.Node{
			ID:      uuid.New(),
			OwnerID: uuid.New(),
			Address: "node.com:9392",
		}

		err := repo.Create(n)
		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}

		n.Address = "example2.com:8080"

		n, err = svc.Update(n.ID, n)
		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}

		node, err := repo.Get(n.ID)
		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}

		if !testNodeMatches(t, n.ID, n.OwnerID, n.Address, node) {
			t.Errorf("Node does not match")
		}
	})

	t.Run("validation", func(t *testing.T) {
		repo := mem.New()
		svc := New(repo)

		n := &node.Node{
			ID:      uuid.New(),
			OwnerID: uuid.New(),
			Address: "valid.com:8000",
		}

		err := repo.Create(n)

		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}

		nn := &node.Node{
			ID:      n.ID,
			OwnerID: n.OwnerID,
		}

		nn.Address = "bad address"

		updated, err := svc.Update(n.ID, nn)

		if err == nil {
			t.Fatal("Expected an error")
		}

		if updated != nil {
			t.Fatal("Expected node, got nil")
		}

		if !strings.Contains(err.Error(), "invalid address") {
			t.Errorf("Expected ErrInvalidAddress, got %v", err)
		}

		check, err := repo.Get(n.ID)
		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}

		if !testNodeMatches(t, n.ID, n.OwnerID, "valid.com:8000", check) {
			t.Errorf("Node does not match")
		}
	})

	t.Run("error not found", func(t *testing.T) {
		repo := mem.New()
		svc := New(repo)

		n := &node.Node{
			ID:      uuid.New(),
			OwnerID: uuid.New(),
			Address: "http://example.com",
		}

		n, err := svc.Update(uuid.New(), n)
		if err != node.ErrNotFound {
			t.Errorf("Expected error, got nil")
		}
	})
}

func TestDelete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := mem.New()
		svc := New(repo)

		n := &node.Node{
			ID:      uuid.New(),
			OwnerID: uuid.New(),
			Address: "http://example.com",
		}

		err := repo.Create(n)
		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}

		err = svc.Delete(n.ID)
		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}

		_, err = repo.Get(n.ID)
		if err != node.ErrNotFound {
			t.Errorf("Expected error, got nil")
		}
	})

	t.Run("error not found", func(t *testing.T) {
		repo := mem.New()
		svc := New(repo)

		err := svc.Delete(uuid.New())
		if err != node.ErrNotFound {
			t.Errorf("Expected error, got nil")
		}
	})
}
