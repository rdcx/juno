package mysql

import (
	"database/sql"
	"encoding/json"
	"juno/pkg/api/ranag"
	"juno/pkg/api/ranag/migration/mysql"
	"testing"

	_ "github.com/go-sql-driver/mysql"

	"github.com/google/uuid"
)

func testRanagMatches(t *testing.T, conn *sql.DB, id, ownerID uuid.UUID, address string, expectedShards [][2]int) bool {
	sqlCheck := "SELECT id, owner_id, address, shard_assignments FROM ranags WHERE id = ?"

	row := conn.QueryRow(sqlCheck, id)

	var ranag ranag.Ranag
	var shardAssignments string
	err := row.Scan(&ranag.ID, &ranag.OwnerID, &ranag.Address, &shardAssignments)
	if err != nil {
		t.Errorf("Error getting ranag: %s", err)
		return false
	}

	// Unmarshal shard_assignments from string
	var shards [][2]int
	if err := json.Unmarshal([]byte(shardAssignments), &shards); err != nil {
		t.Errorf("Error unmarshalling shard assignments: %s", err)
		return false
	}

	ranag.ShardAssignments = shards

	if ranag.ID != id {
		t.Errorf("Expected ID %s, got %s", id, ranag.ID)
		return false
	}

	if ranag.OwnerID != ownerID {
		t.Errorf("Expected OwnerID %s, got %s", ownerID, ranag.OwnerID)
		return false
	}

	if ranag.Address != address {
		t.Errorf("Expected Address %s, got %s", address, ranag.Address)
		return false
	}

	if len(ranag.ShardAssignments) != len(expectedShards) {
		t.Errorf("Expected %d shard assignments, got %d", len(expectedShards), len(ranag.ShardAssignments))
		return false
	}

	for i, shard := range expectedShards {
		if ranag.ShardAssignments[i][0] != shard[0] || ranag.ShardAssignments[i][1] != shard[1] {
			t.Errorf("Expected shard assignment %v, got %v", shard, ranag.ShardAssignments[i])
			return false
		}
	}

	return true
}

func TestAll(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		n1 := ranag.Ranag{
			ID:      uuid.New(),
			OwnerID: uuid.New(),
			Address: "http://example.com",
			ShardAssignments: [][2]int{
				{0, 1}, {1, 2},
			},
		}

		n2 := ranag.Ranag{
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

		conn, err := sql.Open("mysql", "root:juno@tcp(localhost:3306)/ranag_test?parseTime=true")
		if err != nil {
			t.Errorf("Error connecting to database: %s", err)
		}
		err = mysql.ExecuteMigrations(conn)
		if err != nil {
			t.Errorf("Error executing migrations: %s", err)
		}

		defer conn.Close()

		defer conn.Exec("DELETE FROM ranags WHERE id = ?", n1.ID)
		defer conn.Exec("DELETE FROM ranags WHERE id = ?", n2.ID)

		_, err = conn.Exec("INSERT INTO ranags (id, owner_id, address, shard_assignments) VALUES (?, ?, ?, ?)", n1.ID, n1.OwnerID, n1.Address, shardsJSON1)
		if err != nil {
			t.Errorf("Error inserting ranag: %s", err)
		}

		_, err = conn.Exec("INSERT INTO ranags (id, owner_id, address, shard_assignments) VALUES (?, ?, ?, ?)", n2.ID, n2.OwnerID, n2.Address, shardsJSON2)
		if err != nil {
			t.Errorf("Error inserting ranag: %s", err)
		}

		repo := New(conn)

		ranags, err := repo.All()

		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}

		if len(ranags) < 2 {
			t.Errorf("Expected 2 ranags, got %d", len(ranags))
		}
	})
}

func TestCreate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		n := ranag.Ranag{
			ID:      uuid.New(),
			OwnerID: uuid.New(),
			Address: "http://example.com",
			ShardAssignments: [][2]int{
				{0, 1}, {1, 2},
			},
		}

		conn, err := sql.Open("mysql", "root:juno@tcp(localhost:3306)/ranag_test?parseTime=true")
		if err != nil {
			t.Errorf("Error connecting to database: %s", err)
		}
		err = mysql.ExecuteMigrations(conn)
		if err != nil {
			t.Errorf("Error executing migrations: %s", err)
		}

		defer conn.Close()

		defer conn.Exec("DELETE FROM ranags WHERE id = ?", n.ID)

		repo := New(conn)

		err = repo.Create(&n)

		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}

		if !testRanagMatches(t, conn, n.ID, n.OwnerID, n.Address, n.ShardAssignments) {
			t.Errorf("Ranag does not match")
		}
	})

	t.Run("error", func(t *testing.T) {
		n := ranag.Ranag{
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

		conn, err := sql.Open("mysql", "root:juno@tcp(localhost:3306)/ranag_test?parseTime=true")
		if err != nil {
			t.Errorf("Error connecting to database: %s", err)
		}
		err = mysql.ExecuteMigrations(conn)
		if err != nil {
			t.Errorf("Error executing migrations: %s", err)
		}

		defer conn.Close()

		_, err = conn.Exec("INSERT INTO ranags (id, owner_id, address, shard_assignments) VALUES (?, ?, ?, ?)", n.ID, n.OwnerID, n.Address, shardsJSON)

		if err != nil {
			t.Errorf("Error inserting ranag: %s", err)
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
		n := ranag.Ranag{
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

		conn, err := sql.Open("mysql", "root:juno@tcp(localhost:3306)/ranag_test?parseTime=true")
		if err != nil {
			t.Errorf("Error connecting to database: %s", err)
		}
		err = mysql.ExecuteMigrations(conn)
		if err != nil {
			t.Errorf("Error executing migrations: %s", err)
		}

		defer conn.Close()

		defer conn.Exec("DELETE FROM ranags WHERE id = ?", n.ID)

		_, err = conn.Exec("INSERT INTO ranags (id, owner_id, address, shard_assignments) VALUES (?, ?, ?, ?)", n.ID, n.OwnerID, n.Address, shardsJSON)

		if err != nil {
			t.Errorf("Error inserting ranag: %s", err)
		}

		repo := New(conn)

		ranag, err := repo.Get(n.ID)

		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}

		if ranag.ID != n.ID {
			t.Errorf("Expected ID %s, got %s", n.ID, ranag.ID)
		}

		if ranag.OwnerID != n.OwnerID {
			t.Errorf("Expected OwnerID %s, got %s", n.OwnerID, ranag.OwnerID)
		}

		if ranag.Address != n.Address {
			t.Errorf("Expected Address %s, got %s", n.Address, ranag.Address)
		}
	})

	t.Run("error", func(t *testing.T) {
		conn, err := sql.Open("mysql", "root:juno@tcp(localhost:3306)/ranag_test?parseTime=true")
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

		n1 := ranag.Ranag{
			ID:      uuid.New(),
			OwnerID: ownerID,
			Address: "http://example.com",
			ShardAssignments: [][2]int{
				{0, 1}, {1, 2},
			},
		}

		n2 := ranag.Ranag{
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

		conn, err := sql.Open("mysql", "root:juno@tcp(localhost:3306)/ranag_test?parseTime=true")
		if err != nil {
			t.Errorf("Error connecting to database: %s", err)
		}
		err = mysql.ExecuteMigrations(conn)
		if err != nil {
			t.Errorf("Error executing migrations: %s", err)
		}

		defer conn.Close()

		defer conn.Exec("DELETE FROM ranags WHERE id = ?", n1.ID)
		defer conn.Exec("DELETE FROM ranags WHERE id = ?", n2.ID)

		_, err = conn.Exec("INSERT INTO ranags (id, owner_id, address, shard_assignments) VALUES (?, ?, ?, ?)", n1.ID, n1.OwnerID, n1.Address, shardsJSON1)
		if err != nil {
			t.Errorf("Error inserting ranag: %s", err)
		}

		_, err = conn.Exec("INSERT INTO ranags (id, owner_id, address, shard_assignments) VALUES (?, ?, ?, ?)", n2.ID, n2.OwnerID, n2.Address, shardsJSON2)
		if err != nil {
			t.Errorf("Error inserting ranag: %s", err)
		}

		repo := New(conn)

		ranags, err := repo.ListByOwnerID(ownerID)

		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}

		if len(ranags) != 1 {
			t.Errorf("Expected 1 ranags, got %d", len(ranags))
		}

		if ranags[0].ID != n1.ID {
			t.Errorf("Expected ID %s, got %s", n1.ID, ranags[0].ID)
		}

		if ranags[0].OwnerID != n1.OwnerID {
			t.Errorf("Expected OwnerID %s, got %s", n1.OwnerID, ranags[0].OwnerID)
		}

		if ranags[0].Address != n1.Address {
			t.Errorf("Expected Address %s, got %s", n1.Address, ranags[0].Address)
		}
	})
}

func TestUpdate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		n := ranag.Ranag{
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

		conn, err := sql.Open("mysql", "root:juno@tcp(localhost:3306)/ranag_test?parseTime=true")
		if err != nil {
			t.Errorf("Error connecting to database: %s", err)
		}
		err = mysql.ExecuteMigrations(conn)
		if err != nil {
			t.Errorf("Error executing migrations: %s", err)
		}

		defer conn.Close()

		defer conn.Exec("DELETE FROM ranags WHERE id = ?", n.ID)

		_, err = conn.Exec("INSERT INTO ranags (id, owner_id, address, shard_assignments) VALUES (?, ?, ?, ?)", n.ID, n.OwnerID, n.Address, shardsJSON)

		if err != nil {
			t.Errorf("Error inserting ranag: %s", err)
		}

		repo := New(conn)

		n.Address = "http://example.org"

		err = repo.Update(&n)

		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}

		if !testRanagMatches(t, conn, n.ID, n.OwnerID, n.Address, n.ShardAssignments) {
			t.Errorf("Ranag does not match")
		}
	})
}

func TestDelete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		n := ranag.Ranag{
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

		conn, err := sql.Open("mysql", "root:juno@tcp(localhost:3306)/ranag_test?parseTime=true")
		if err != nil {
			t.Errorf("Error connecting to database: %s", err)
		}
		err = mysql.ExecuteMigrations(conn)
		if err != nil {
			t.Errorf("Error executing migrations: %s", err)
		}

		defer conn.Close()

		_, err = conn.Exec("INSERT INTO ranags (id, owner_id, address, shard_assignments) VALUES (?, ?, ?, ?)", n.ID, n.OwnerID, n.Address, shardsJSON)

		if err != nil {
			t.Errorf("Error inserting ranag: %s", err)
		}

		repo := New(conn)

		err = repo.Delete(n.ID)

		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}

		row := conn.QueryRow("SELECT COUNT(*) FROM ranags WHERE id = ?", n.ID)

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
