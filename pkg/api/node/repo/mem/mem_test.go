package mem

import (
	"juno/pkg/api/node"
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

func TestAll(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		n1 := node.Node{
			ID:      uuid.New(),
			OwnerID: uuid.New(),
			Address: "http://example.com",
		}
		n2 := node.Node{
			ID:      uuid.New(),
			OwnerID: uuid.New(),
			Address: "http://example.org",
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

		nodes, err := repo.All()
		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}

		if len(nodes) != 2 {
			t.Errorf("Expected 2 nodes, got %d", len(nodes))
		}

		for _, n := range nodes {
			if n.ID == n1.ID {
				if !testNodeMatches(t, n1.ID, n1.OwnerID, n1.Address, n) {
					t.Errorf("Node 1 does not match")
				}
			} else if n.ID == n2.ID {
				if !testNodeMatches(t, n2.ID, n2.OwnerID, n2.Address, n) {
					t.Errorf("Node 2 does not match")
				}
			} else {
				t.Errorf("Unexpected node ID: %s", n.ID)
			}
		}
	})
}

func TestCreate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		n := node.Node{
			ID:      uuid.New(),
			OwnerID: uuid.New(),
			Address: "http://example.com",
		}

		repo := New()

		err := repo.Create(&n)

		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}

		node := repo.nodes[n.ID]

		if !testNodeMatches(t, n.ID, n.OwnerID, n.Address, node) {
			t.Errorf("Node does not match")
		}
	})

	t.Run("error", func(t *testing.T) {
		n := node.Node{
			ID:      uuid.New(),
			OwnerID: uuid.New(),
			Address: "http://example.com",
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

		if !testNodeMatches(t, n.ID, n.OwnerID, n.Address, node) {
			t.Errorf("Node does not match")
		}
	})

	t.Run("error", func(t *testing.T) {
		n := node.Node{
			ID:      uuid.New(),
			OwnerID: uuid.New(),
			Address: "http://example.com",
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
		}
		n2 := node.Node{
			ID:      uuid.New(),
			OwnerID: ownerID,
			Address: "http://example.org",
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

		if !testNodeMatches(t, n1.ID, n1.OwnerID, n1.Address, nodes[0]) {
			t.Errorf("Node 1 does not match")
		}

		if !testNodeMatches(t, n2.ID, n2.OwnerID, n2.Address, nodes[1]) {
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
		}

		repo := New()

		err := repo.Create(&n)
		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}

		n.Address = "http://example.org"

		err = repo.Update(&n)
		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}

		node := repo.nodes[n.ID]

		if !testNodeMatches(t, n.ID, n.OwnerID, n.Address, node) {
			t.Errorf("Node does not match")
		}
	})

	t.Run("error", func(t *testing.T) {
		n := node.Node{
			ID:      uuid.New(),
			OwnerID: uuid.New(),
			Address: "http://example.com",
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
