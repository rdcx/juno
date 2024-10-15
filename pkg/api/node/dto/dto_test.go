package dto

import "testing"

func TestNodeToDomain(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		n := Node{
			ID:      "00000000-0000-0000-0000-000000000000",
			OwnerID: "00000000-0000-0000-0000-000000000001",
			Address: "http://node.com",
			Offset:  100,
			Shards:  50,
		}

		d, err := n.ToDomain()

		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}

		if d.ID.String() != n.ID {
			t.Errorf("Expected ID %s, got %s", n.ID, d.ID.String())
		}

		if d.OwnerID.String() != n.OwnerID {
			t.Errorf("Expected OwnerID %s, got %s", n.OwnerID, d.OwnerID.String())
		}

		if d.Address != n.Address {
			t.Errorf("Expected Address %s, got %s", n.Address, d.Address)
		}

		if d.Offset != n.Offset {
			t.Errorf("Expected Offset %d, got %d", n.Offset, d.Offset)
		}

		if d.Shards != n.Shards {
			t.Errorf("Expected %d shards, got %d", n.Shards, d.Shards)
		}
	})

	t.Run("invalid id", func(t *testing.T) {
		n := Node{
			ID:      "00000000-0000-0000-0000-000000000000",
			OwnerID: "00000000-0000-0000-0000-000000000001",
			Address: "http://node.com",
			Offset:  1,
			Shards:  3,
		}

		n.ID = "invalid"

		_, err := n.ToDomain()

		if err == nil {
			t.Errorf("Expected an error")
		}
	})

	t.Run("invalid owner id", func(t *testing.T) {
		n := Node{
			ID:      "00000000-0000-0000-0000-000000000000",
			OwnerID: "00000000-0000-0000-0000-000000000001",
			Address: "http://node.com",
			Offset:  0,
			Shards:  1000,
		}

		n.OwnerID = "invalid"

		_, err := n.ToDomain()

		if err == nil {
			t.Errorf("Expected an error")
		}
	})
}
