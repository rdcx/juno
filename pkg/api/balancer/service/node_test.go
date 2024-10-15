package service

import (
	"juno/pkg/api/balancer"
	"juno/pkg/api/balancer/repo/mem"
	"juno/pkg/api/user"
	"strings"
	"testing"

	"github.com/google/uuid"
)

func testBalancerMatches(t *testing.T, id, ownerID uuid.UUID, address string, shards []int, n *balancer.Balancer) bool {
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

		n, err := svc.Create(u.ID, addr, shards)

		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}

		balancer, err := repo.Get(n.ID)

		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}

		if !testBalancerMatches(t, n.ID, n.OwnerID, addr, shards, balancer) {
			t.Errorf("Balancer does not match")
		}
	})

	t.Run("validation", func(t *testing.T) {
		repo := mem.New()
		svc := New(repo)

		u := &user.User{
			ID: uuid.New(),
		}

		addr := "bad address"

		n, err := svc.Create(u.ID, addr, []int{-1, 2, 3, 100 * 10000})

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

		n := &balancer.Balancer{
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

		n, err = svc.Create(u.ID, addr, shards)

		if err != balancer.ErrAddressExists {
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

		n := &balancer.Balancer{
			ID:      uuid.New(),
			OwnerID: uuid.New(),
			Address: "example.com:8000",
			Shards:  []int{1, 2, 3},
		}

		err := repo.Create(n)
		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}

		balancer, err := svc.Get(n.ID)
		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}

		if !testBalancerMatches(t, n.ID, n.OwnerID, n.Address, n.Shards, balancer) {
			t.Errorf("Balancer does not match")
		}
	})

	t.Run("error not found", func(t *testing.T) {
		repo := mem.New()
		svc := New(repo)

		n := &balancer.Balancer{
			ID:      uuid.New(),
			OwnerID: uuid.New(),
		}

		_, err := svc.Get(n.ID)
		if err != balancer.ErrNotFound {
			t.Errorf("Expected error, got nil")
		}
	})
}

func TestListByOwnerID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := mem.New()
		svc := New(repo)

		ownerID := uuid.New()

		n1 := &balancer.Balancer{
			ID:      uuid.New(),
			OwnerID: ownerID,
			Address: "example.com:8000",
			Shards:  []int{1, 2, 3},
		}

		err := repo.Create(n1)
		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}

		n2 := &balancer.Balancer{
			ID:      uuid.New(),
			OwnerID: uuid.New(),
			Address: "example2.com:8080",
			Shards:  []int{4, 5, 6},
		}

		err = repo.Create(n2)

		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}

		balancers, err := svc.ListByOwnerID(ownerID)
		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}

		if len(balancers) != 1 {
			t.Errorf("Expected 1 balancers, got %d", len(balancers))
		}

		if !testBalancerMatches(t, n1.ID, n1.OwnerID, n1.Address, n1.Shards, balancers[0]) {
			t.Errorf("Balancer does not match")
		}
	})
}

func TestUpdate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := mem.New()
		svc := New(repo)

		n := &balancer.Balancer{
			ID:      uuid.New(),
			OwnerID: uuid.New(),
			Address: "balancer.com:9392",
			Shards:  []int{1, 2, 3},
		}

		err := repo.Create(n)
		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}

		n.Address = "example2.com:8080"
		n.Shards = []int{4, 5, 6}

		n, err = svc.Update(n.ID, n)
		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}

		balancer, err := repo.Get(n.ID)
		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}

		if !testBalancerMatches(t, n.ID, n.OwnerID, n.Address, n.Shards, balancer) {
			t.Errorf("Balancer does not match")
		}
	})

	t.Run("validation", func(t *testing.T) {
		repo := mem.New()
		svc := New(repo)

		n := &balancer.Balancer{
			ID:      uuid.New(),
			OwnerID: uuid.New(),
			Address: "valid.com:8000",
			Shards:  []int{1, 2, 3},
		}

		err := repo.Create(n)

		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}

		nn := &balancer.Balancer{
			ID:      n.ID,
			OwnerID: n.OwnerID,
		}

		nn.Address = "bad address"
		nn.Shards = []int{-1, 2, 3, 100 * 10000}

		updated, err := svc.Update(n.ID, nn)

		if err == nil {
			t.Fatal("Expected an error")
		}

		if updated != nil {
			t.Fatal("Expected balancer, got nil")
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

		if !testBalancerMatches(t, n.ID, n.OwnerID, "valid.com:8000", []int{1, 2, 3}, check) {
			t.Errorf("Balancer does not match")
		}
	})

	t.Run("error not found", func(t *testing.T) {
		repo := mem.New()
		svc := New(repo)

		n := &balancer.Balancer{
			ID:      uuid.New(),
			OwnerID: uuid.New(),
			Address: "http://example.com",
			Shards:  []int{1, 2, 3},
		}

		n, err := svc.Update(uuid.New(), n)
		if err != balancer.ErrNotFound {
			t.Errorf("Expected error, got nil")
		}
	})
}

func TestDelete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := mem.New()
		svc := New(repo)

		n := &balancer.Balancer{
			ID:      uuid.New(),
			OwnerID: uuid.New(),
			Address: "http://example.com",
			Shards:  []int{1, 2, 3},
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
		if err != balancer.ErrNotFound {
			t.Errorf("Expected error, got nil")
		}
	})

	t.Run("error not found", func(t *testing.T) {
		repo := mem.New()
		svc := New(repo)

		err := svc.Delete(uuid.New())
		if err != balancer.ErrNotFound {
			t.Errorf("Expected error, got nil")
		}
	})
}
