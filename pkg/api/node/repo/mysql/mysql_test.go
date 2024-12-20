package mysql

import (
	"database/sql"
	"encoding/json"
	"juno/pkg/api/node"
	"juno/pkg/api/node/migration/mysql"
	"testing"

	_ "github.com/go-sql-driver/mysql"

	"github.com/google/uuid"
)

func testNodeMatches(t *testing.T, conn *sql.DB, id, ownerID uuid.UUID, address string, expectedShards [][2]int) bool {
	sqlCheck := "SELECT id, owner_id, address, shard_assignments FROM nodes WHERE id = ?"

	row := conn.QueryRow(sqlCheck, id)

	var node node.Node
	var shardAssignments string
	err := row.Scan(&node.ID, &node.OwnerID, &node.Address, &shardAssignments)
	if err != nil {
		t.Errorf("Error getting node: %s", err)
		return false
	}

	// Unmarshal shard_assignments from string
	var shards [][2]int
	if err := json.Unmarshal([]byte(shardAssignments), &shards); err != nil {
		t.Errorf("Error unmarshalling shard assignments: %s", err)
		return false
	}

	node.ShardAssignments = shards

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

	if len(node.ShardAssignments) != len(expectedShards) {
		t.Errorf("Expected %d shard assignments, got %d", len(expectedShards), len(node.ShardAssignments))
		return false
	}

	for i, shard := range expectedShards {
		if node.ShardAssignments[i][0] != shard[0] || node.ShardAssignments[i][1] != shard[1] {
			t.Errorf("Expected shard assignment %v, got %v", shard, node.ShardAssignments[i])
			return false
		}
	}

	return true
}

func TestAll(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		n1 := node.Node{
			ID:      uuid.New(),
			OwnerID: uuid.New(),
			Address: "http://example.com",
			ShardAssignments: [][2]int{
				{0, 1}, {1, 2},
			},
		}

		n2 := node.Node{
			ID:      uuid.New(),
			OwnerID: uuid.New(),
			Address: "http://example.org",
			ShardAssignments: [][2]int{
				{0, 1}, {1, 2},
			},
		}

		// Marshal shard_assignments to JSON string
		shardsJSON1, err := json.Marshal(n1.ShardAssignments)
		if err != nil {
			t.Errorf("Error marshalling shard assignments: %s", err)
		}
		shardsJSON2, err := json.Marshal(n2.ShardAssignments)
		if err != nil {
			t.Errorf("Error marshalling shard assignments: %s", err)
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

		defer conn.Exec("DELETE FROM nodes WHERE id = ?", n1.ID)
		defer conn.Exec("DELETE FROM nodes WHERE id = ?", n2.ID)

		_, err = conn.Exec("INSERT INTO nodes (id, owner_id, address, shard_assignments) VALUES (?, ?, ?, ?)", n1.ID, n1.OwnerID, n1.Address, shardsJSON1)
		if err != nil {
			t.Errorf("Error inserting node: %s", err)
		}

		_, err = conn.Exec("INSERT INTO nodes (id, owner_id, address, shard_assignments) VALUES (?, ?, ?, ?)", n2.ID, n2.OwnerID, n2.Address, shardsJSON2)
		if err != nil {
			t.Errorf("Error inserting node: %s", err)
		}

		repo := New(conn)

		nodes, err := repo.All()

		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}

		if len(nodes) < 2 {
			t.Errorf("Expected 2 nodes, got %d", len(nodes))
		}
	})
}

func TestCreate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		n := node.Node{
			ID:      uuid.New(),
			OwnerID: uuid.New(),
			Address: "http://example.com",
			ShardAssignments: [][2]int{
				{0, 1}, {1, 2},
			},
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

		if !testNodeMatches(t, conn, n.ID, n.OwnerID, n.Address, n.ShardAssignments) {
			t.Errorf("Node does not match")
		}
	})

	t.Run("error", func(t *testing.T) {
		n := node.Node{
			ID:      uuid.New(),
			OwnerID: uuid.New(),
			Address: "http://example.com",
			ShardAssignments: [][2]int{
				{0, 1}, {1, 2},
			},
		}

		// Marshal shard_assignments to JSON string
		shardsJSON, err := json.Marshal(n.ShardAssignments)
		if err != nil {
			t.Errorf("Error marshalling shard assignments: %s", err)
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

		_, err = conn.Exec("INSERT INTO nodes (id, owner_id, address, shard_assignments) VALUES (?, ?, ?, ?)", n.ID, n.OwnerID, n.Address, shardsJSON)

		if err != nil {
			t.Errorf("Error inserting node: %s", err)
		}

		repo := New(conn)

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
			ShardAssignments: [][2]int{
				{0, 1}, {1, 2},
			},
		}

		// Marshal shard_assignments to JSON string
		shardsJSON, err := json.Marshal(n.ShardAssignments)
		if err != nil {
			t.Errorf("Error marshalling shard assignments: %s", err)
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

		_, err = conn.Exec("INSERT INTO nodes (id, owner_id, address, shard_assignments) VALUES (?, ?, ?, ?)", n.ID, n.OwnerID, n.Address, shardsJSON)

		if err != nil {
			t.Errorf("Error inserting node: %s", err)
		}

		repo := New(conn)

		node, err := repo.Get(n.ID)

		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}

		if node.ID != n.ID {
			t.Errorf("Expected ID %s, got %s", n.ID, node.ID)
		}

		if node.OwnerID != n.OwnerID {
			t.Errorf("Expected OwnerID %s, got %s", n.OwnerID, node.OwnerID)
		}

		if node.Address != n.Address {
			t.Errorf("Expected Address %s, got %s", n.Address, node.Address)
		}
	})

	t.Run("error", func(t *testing.T) {
		conn, err := sql.Open("mysql", "root:juno@tcp(localhost:3306)/node_test?parseTime=true")
		if err != nil {
			t.Errorf("Error connecting to database: %s", err)
		}

		defer conn.Close()

		repo := New(conn)

		_, err = repo.Get(uuid.New())

		if err == nil {
			t.Errorf("Expected an error")
		}

		if err.Error() != "sql: no rows in result set" {
			t.Errorf("Unexpected error: %s", err)
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
			ShardAssignments: [][2]int{
				{0, 1}, {1, 2},
			},
		}

		n2 := node.Node{
			ID:      uuid.New(),
			OwnerID: uuid.New(),
			Address: "http://example.org",
			ShardAssignments: [][2]int{
				{0, 1}, {1, 2},
			},
		}

		// Marshal shard_assignments to JSON string
		shardsJSON1, err := json.Marshal(n1.ShardAssignments)
		if err != nil {
			t.Errorf("Error marshalling shard assignments: %s", err)
		}
		shardsJSON2, err := json.Marshal(n2.ShardAssignments)
		if err != nil {
			t.Errorf("Error marshalling shard assignments: %s", err)
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

		defer conn.Exec("DELETE FROM nodes WHERE id = ?", n1.ID)
		defer conn.Exec("DELETE FROM nodes WHERE id = ?", n2.ID)

		_, err = conn.Exec("INSERT INTO nodes (id, owner_id, address, shard_assignments) VALUES (?, ?, ?, ?)", n1.ID, n1.OwnerID, n1.Address, shardsJSON1)
		if err != nil {
			t.Errorf("Error inserting node: %s", err)
		}

		_, err = conn.Exec("INSERT INTO nodes (id, owner_id, address, shard_assignments) VALUES (?, ?, ?, ?)", n2.ID, n2.OwnerID, n2.Address, shardsJSON2)
		if err != nil {
			t.Errorf("Error inserting node: %s", err)
		}

		repo := New(conn)

		nodes, err := repo.ListByOwnerID(ownerID)

		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}

		if len(nodes) != 1 {
			t.Errorf("Expected 1 nodes, got %d", len(nodes))
		}

		if nodes[0].ID != n1.ID {
			t.Errorf("Expected ID %s, got %s", n1.ID, nodes[0].ID)
		}

		if nodes[0].OwnerID != n1.OwnerID {
			t.Errorf("Expected OwnerID %s, got %s", n1.OwnerID, nodes[0].OwnerID)
		}

		if nodes[0].Address != n1.Address {
			t.Errorf("Expected Address %s, got %s", n1.Address, nodes[0].Address)
		}
	})
}

func TestUpdate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		n := node.Node{
			ID:      uuid.New(),
			OwnerID: uuid.New(),
			Address: "http://example.com",
			ShardAssignments: [][2]int{
				{0, 1}, {1, 2},
			},
		}

		// Marshal shard_assignments to JSON string
		shardsJSON, err := json.Marshal(n.ShardAssignments)
		if err != nil {
			t.Errorf("Error marshalling shard assignments: %s", err)
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

		_, err = conn.Exec("INSERT INTO nodes (id, owner_id, address, shard_assignments) VALUES (?, ?, ?, ?)", n.ID, n.OwnerID, n.Address, shardsJSON)

		if err != nil {
			t.Errorf("Error inserting node: %s", err)
		}

		repo := New(conn)

		n.Address = "http://example.org"

		err = repo.Update(&n)

		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}

		if !testNodeMatches(t, conn, n.ID, n.OwnerID, n.Address, n.ShardAssignments) {
			t.Errorf("Node does not match")
		}
	})
}

func TestDelete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		n := node.Node{
			ID:      uuid.New(),
			OwnerID: uuid.New(),
			Address: "http://example.com",
			ShardAssignments: [][2]int{
				{0, 1}, {1, 2},
			},
		}

		// Marshal shard_assignments to JSON string
		shardsJSON, err := json.Marshal(n.ShardAssignments)
		if err != nil {
			t.Errorf("Error marshalling shard assignments: %s", err)
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

		_, err = conn.Exec("INSERT INTO nodes (id, owner_id, address, shard_assignments) VALUES (?, ?, ?, ?)", n.ID, n.OwnerID, n.Address, shardsJSON)

		if err != nil {
			t.Errorf("Error inserting node: %s", err)
		}

		repo := New(conn)

		err = repo.Delete(n.ID)

		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}

		row := conn.QueryRow("SELECT COUNT(*) FROM nodes WHERE id = ?", n.ID)

		var count int
		err = row.Scan(&count)
		if err != nil {
			t.Errorf("Error getting count: %s", err)
		}

		if count != 0 {
			t.Errorf("Expected 0 rows, got %d", count)
		}
	})
}
