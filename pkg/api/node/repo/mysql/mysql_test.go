package mysql

import (
	"database/sql"
	"encoding/json"
	"juno/pkg/api/node/domain"
	"juno/pkg/api/node/migration/mysql"
	"testing"

	_ "github.com/go-sql-driver/mysql"

	"github.com/google/uuid"
)

func testNodeMatches(t *testing.T, conn *sql.DB, id, ownerID uuid.UUID, address string, shards []int) bool {
	sqlCheck := "SELECT id, owner_id, address, shards FROM nodes WHERE id = ?"

	row := conn.QueryRow(sqlCheck, id)

	var jsonShards string
	var node domain.Node
	err := row.Scan(&node.ID, &node.OwnerID, &node.Address, &jsonShards)
	if err != nil {
		t.Errorf("Error getting node: %s", err)
		return false
	}

	err = json.Unmarshal([]byte(jsonShards), &node.Shards)
	if err != nil {
		t.Errorf("Error unmarshalling shards: %s", err)
		return false
	}

	if node.ID != id {
		t.Errorf("Expected ID %s, got %s", id, node.ID)
		return false
	}

	if node.OwnerID != ownerID {
		t.Errorf("Expected OwnerID %s, got %s", ownerID, node.OwnerID)
		return false
	}

	if node.Address != address {
		t.Errorf("Expected Address %s, got %s", address, node.Address)
		return false
	}

	if len(node.Shards) != len(shards) {
		t.Errorf("Expected %d shards, got %d", len(shards), len(node.Shards))
		return false
	}

	for i, shard := range shards {
		if node.Shards[i] != shard {
			t.Errorf("Expected shard %d, got %d", shard, node.Shards[i])
			return false
		}
	}

	return true
}

func TestCreate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		n := domain.Node{
			ID:      uuid.New(),
			OwnerID: uuid.New(),
			Address: "http://example.com",
			Shards:  []int{1, 2, 3},
		}

		conn, err := sql.Open("mysql", "root:juno@tcp(localhost:3306)/node_test?parseTime=true")
		if err != nil {
			t.Errorf("Error connecting to database: %s", err)
		}
		err = mysql.ExecuteMigrations(conn)
		if err != nil {
			t.Errorf("Error executing migrations: %s", err)
		}

		defer conn.Close()

		defer conn.Exec("DELETE FROM nodes WHERE id = ?", n.ID)

		repo := New(conn)

		err = repo.Create(&n)

		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}

		if !testNodeMatches(t, conn, n.ID, n.OwnerID, n.Address, n.Shards) {
			t.Errorf("Node does not match")
		}
	})

	t.Run("error", func(t *testing.T) {
		// test code
	})
}
