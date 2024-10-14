package mem

import (
	"juno/pkg/api/node"
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
		n := node.Node{
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

		node := repo.nodes[n.ID]

		if !testNodeMatches(t, n.ID, n.OwnerID, n.Address, n.Shards, node) {
			t.Errorf("Node does not match")
		}
	})

	t.Run("error", func(t *testing.T) {
		n := node.Node{
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
		n := node.Node{
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

		node, err := repo.Get(n.ID)
		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}

		if !testNodeMatches(t, n.ID, n.OwnerID, n.Address, n.Shards, node) {
			t.Errorf("Node does not match")
		}
	})

	t.Run("error", func(t *testing.T) {
		n := node.Node{
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
		n1 := node.Node{
			ID:      uuid.New(),
			OwnerID: ownerID,
			Address: "http://example.com",
			Shards:  []int{1, 2, 3},
		}
		n2 := node.Node{
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

		nodes, err := repo.ListByOwnerID(ownerID)
		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}

		if len(nodes) != 2 {
			t.Errorf("Expected 2 nodes, got %d", len(nodes))
		}

		if !testNodeMatches(t, n1.ID, n1.OwnerID, n1.Address, n1.Shards, nodes[0]) {
			t.Errorf("Node 1 does not match")
		}

		if !testNodeMatches(t, n2.ID, n2.OwnerID, n2.Address, n2.Shards, nodes[1]) {
			t.Errorf("Node 2 does not match")
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
		n := node.Node{
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

		node := repo.nodes[n.ID]

		if !testNodeMatches(t, n.ID, n.OwnerID, n.Address, n.Shards, node) {
			t.Errorf("Node does not match")
		}
	})

	t.Run("error", func(t *testing.T) {
		n := node.Node{
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
		n := node.Node{
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

		_, ok := repo.nodes[n.ID]
		if ok {
			t.Errorf("Node was not deleted")
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
