package mem

import (
	"juno/pkg/api/balancer"
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
		n := balancer.Balancer{
			ID:      uuid.New(),
			OwnerID: uuid.New(),
			Address: "http://example.com",
			Shards:  []int{1, 2, 3},
		}

		repo := New()

		err := repo.Create(&n)

		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}

		balancer := repo.balancers[n.ID]

		if !testBalancerMatches(t, n.ID, n.OwnerID, n.Address, n.Shards, balancer) {
			t.Errorf("Balancer does not match")
		}
	})

	t.Run("error", func(t *testing.T) {
		n := balancer.Balancer{
			ID:      uuid.New(),
			OwnerID: uuid.New(),
			Address: "http://example.com",
			Shards:  []int{1, 2, 3},
		}

		repo := New()

		err := repo.Create(&n)
		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}

		err = repo.Create(&n)
		if err == nil {
			t.Errorf("Expected an error")
		}
	})
}

func TestGet(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		n := balancer.Balancer{
			ID:      uuid.New(),
			OwnerID: uuid.New(),
			Address: "http://example.com",
			Shards:  []int{1, 2, 3},
		}

		repo := New()

		err := repo.Create(&n)
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

	t.Run("error", func(t *testing.T) {
		n := balancer.Balancer{
			ID:      uuid.New(),
			OwnerID: uuid.New(),
			Address: "http://example.com",
			Shards:  []int{1, 2, 3},
		}

		repo := New()

		_, err := repo.Get(n.ID)
		if err == nil {
			t.Errorf("Expected an error")
		}
	})
}

func TestListByOwnerID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ownerID := uuid.New()
		n1 := balancer.Balancer{
			ID:      uuid.New(),
			OwnerID: ownerID,
			Address: "http://example.com",
			Shards:  []int{1, 2, 3},
		}
		n2 := balancer.Balancer{
			ID:      uuid.New(),
			OwnerID: ownerID,
			Address: "http://example.org",
			Shards:  []int{4, 5, 6},
		}

		repo := New()

		err := repo.Create(&n1)
		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}

		err = repo.Create(&n2)
		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}

		balancers, err := repo.ListByOwnerID(ownerID)
		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}

		if len(balancers) != 2 {
			t.Errorf("Expected 2 balancers, got %d", len(balancers))
		}

		if !testBalancerMatches(t, n1.ID, n1.OwnerID, n1.Address, n1.Shards, balancers[0]) {
			t.Errorf("Balancer 1 does not match")
		}

		if !testBalancerMatches(t, n2.ID, n2.OwnerID, n2.Address, n2.Shards, balancers[1]) {
			t.Errorf("Balancer 2 does not match")
		}
	})

	t.Run("error", func(t *testing.T) {
		repo := New()

		_, err := repo.ListByOwnerID(uuid.New())
		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}
	})
}

func TestUpdate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		n := balancer.Balancer{
			ID:      uuid.New(),
			OwnerID: uuid.New(),
			Address: "http://example.com",
			Shards:  []int{1, 2, 3},
		}

		repo := New()

		err := repo.Create(&n)
		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}

		n.Address = "http://example.org"
		n.Shards = []int{4, 5, 6}

		err = repo.Update(&n)
		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}

		balancer := repo.balancers[n.ID]

		if !testBalancerMatches(t, n.ID, n.OwnerID, n.Address, n.Shards, balancer) {
			t.Errorf("Balancer does not match")
		}
	})

	t.Run("error", func(t *testing.T) {
		n := balancer.Balancer{
			ID:      uuid.New(),
			OwnerID: uuid.New(),
			Address: "http://example.com",
			Shards:  []int{1, 2, 3},
		}

		repo := New()

		err := repo.Update(&n)
		if err == nil {
			t.Errorf("Expected an error")
		}
	})
}

func TestDelete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		n := balancer.Balancer{
			ID:      uuid.New(),
			OwnerID: uuid.New(),
			Address: "http://example.com",
			Shards:  []int{1, 2, 3},
		}

		repo := New()

		err := repo.Create(&n)
		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}

		err = repo.Delete(n.ID)
		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}

		_, ok := repo.balancers[n.ID]
		if ok {
			t.Errorf("Balancer was not deleted")
		}
	})

	t.Run("error", func(t *testing.T) {
		repo := New()

		err := repo.Delete(uuid.New())
		if err == nil {
			t.Errorf("Expected an error")
		}
	})
}
