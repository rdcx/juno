package service

import (
	"juno/pkg/api/node"
	"juno/pkg/api/node/repo/mem"
	"juno/pkg/api/user"
	"strings"
	"testing"

	"github.com/google/uuid"
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

func TestCreate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := mem.New()
		svc := New(repo)

		u := &user.User{
			ID: uuid.New(),
		}

		addr := "example.com:7000"
		shards := []int{1, 2, 3}

		n, err := svc.Create(u, addr, shards)

		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}

		node, err := repo.Get(n.ID)

		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}

		if !testNodeMatches(t, n.ID, n.OwnerID, addr, shards, node) {
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

		n, err := svc.Create(u, addr, []int{-1, 2, 3, 100 * 10000})

		if err == nil {
			t.Fatal("Expected an error")
		}

		if !strings.Contains(err.Error(), "invalid address") {
			t.Errorf("Expected ErrInvalidAddress, got %v", err)
		}

		if !strings.Contains(err.Error(), "invalid shards") {
			t.Errorf("Expected ErrInvalidShards, got %v", err)
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
			Shards:  []int{1, 2, 3},
		}

		u := &user.User{
			ID: n.OwnerID,
		}

		err := repo.Create(n)

		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}

		addr := "http://example.com"
		shards := []int{4, 5, 6}

		n, err = svc.Create(u, addr, shards)

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
			Shards:  []int{1, 2, 3},
		}

		err := repo.Create(n)
		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}

		node, err := svc.Get(&user.User{ID: n.OwnerID}, n.ID)
		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}

		if !testNodeMatches(t, n.ID, n.OwnerID, n.Address, n.Shards, node) {
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

		_, err := svc.Get(&user.User{ID: n.OwnerID}, n.ID)
		if err != node.ErrNotFound {
			t.Errorf("Expected error, got nil")
		}
	})

	t.Run("error unauthorized", func(t *testing.T) {
		repo := mem.New()
		svc := New(repo)

		n := &node.Node{
			ID:      uuid.New(),
			OwnerID: uuid.New(),
			Address: "http://example.com",
			Shards:  []int{1, 2, 3},
		}

		err := repo.Create(n)
		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}

		_, err = svc.Get(&user.User{ID: uuid.New()}, n.ID)
		if err != node.ErrUnauthorized {
			t.Errorf("Expected error, got nil")
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
			Address: "http://example.com",
			Shards:  []int{1, 2, 3},
		}

		err := repo.Create(n)
		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}

		n.Address = "http://example2.com"
		n.Shards = []int{4, 5, 6}

		err = svc.Update(&user.User{ID: n.OwnerID}, n)
		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}

		node, err := repo.Get(n.ID)
		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}

		if !testNodeMatches(t, n.ID, n.OwnerID, n.Address, n.Shards, node) {
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
			Shards:  []int{1, 2, 3},
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
		nn.Shards = []int{-1, 2, 3, 100 * 10000}

		err = svc.Update(&user.User{ID: n.OwnerID}, nn)

		if err == nil {
			t.Fatal("Expected an error")
		}

		if !strings.Contains(err.Error(), "invalid address") {
			t.Errorf("Expected ErrInvalidAddress, got %v", err)
		}

		if !strings.Contains(err.Error(), "invalid shards") {
			t.Errorf("Expected ErrInvalidShards, got %v", err)
		}

		check, err := repo.Get(n.ID)
		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}

		if !testNodeMatches(t, n.ID, n.OwnerID, "valid.com:8000", []int{1, 2, 3}, check) {
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
			Shards:  []int{1, 2, 3},
		}

		err := svc.Update(&user.User{ID: n.OwnerID}, n)
		if err != node.ErrNotFound {
			t.Errorf("Expected error, got nil")
		}
	})

	t.Run("error unauthorized", func(t *testing.T) {
		repo := mem.New()
		svc := New(repo)

		n := &node.Node{
			ID:      uuid.New(),
			OwnerID: uuid.New(),
			Address: "http://example.com",
			Shards:  []int{1, 2, 3},
		}

		err := repo.Create(n)
		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}

		newN := &node.Node{
			ID:      n.ID,
			OwnerID: uuid.New(),
			Address: "http://example2.com",
			Shards:  []int{4, 5, 6},
		}

		err = svc.Update(&user.User{ID: uuid.New()}, newN)

		if err != node.ErrUnauthorized {
			t.Errorf("Expected error, got nil")
		}

		check, err := repo.Get(n.ID)
		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}

		if !testNodeMatches(t, n.ID, n.OwnerID, "http://example.com", []int{1, 2, 3}, check) {
			t.Errorf("Node does not match")
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
			Shards:  []int{1, 2, 3},
		}

		err := repo.Create(n)
		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}

		err = svc.Delete(&user.User{ID: n.OwnerID}, n.ID)
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

		err := svc.Delete(&user.User{ID: uuid.New()}, uuid.New())
		if err != node.ErrNotFound {
			t.Errorf("Expected error, got nil")
		}
	})

	t.Run("error unauthorized", func(t *testing.T) {
		repo := mem.New()
		svc := New(repo)

		n := &node.Node{
			ID:      uuid.New(),
			OwnerID: uuid.New(),
			Address: "http://example.com",
			Shards:  []int{1, 2, 3},
		}

		err := repo.Create(n)
		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}

		err = svc.Delete(&user.User{ID: uuid.New()}, n.ID)
		if err != node.ErrUnauthorized {
			t.Errorf("Expected error, got nil")
		}

		_, err = repo.Get(n.ID)
		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}
	})
}
