package mysql

import (
	"database/sql"
	"encoding/json"
	"juno/pkg/api/balancer"
	"juno/pkg/api/balancer/migration/mysql"
	"testing"

	_ "github.com/go-sql-driver/mysql"

	"github.com/google/uuid"
)

func testBalancerMatches(t *testing.T, conn *sql.DB, id, ownerID uuid.UUID, address string, expectedShards [][2]int) bool {
	sqlCheck := "SELECT id, owner_id, address, shard_assignments FROM balancers WHERE id = ?"

	row := conn.QueryRow(sqlCheck, id)

	var balancer balancer.Balancer
	var shardAssignments string
	err := row.Scan(&balancer.ID, &balancer.OwnerID, &balancer.Address, &shardAssignments)
	if err != nil {
		t.Errorf("Error getting balancer: %s", err)
		return false
	}

	// Unmarshal shard_assignments from string
	var shards [][2]int
	if err := json.Unmarshal([]byte(shardAssignments), &shards); err != nil {
		t.Errorf("Error unmarshalling shard assignments: %s", err)
		return false
	}

	balancer.ShardAssignments = shards

	if balancer.ID != id {
		t.Errorf("Expected ID %s, got %s", id, balancer.ID)
		return false
	}

	if balancer.OwnerID != ownerID {
		t.Errorf("Expected OwnerID %s, got %s", ownerID, balancer.OwnerID)
		return false
	}

	if balancer.Address != address {
		t.Errorf("Expected Address %s, got %s", address, balancer.Address)
		return false
	}

	if len(balancer.ShardAssignments) != len(expectedShards) {
		t.Errorf("Expected %d shard assignments, got %d", len(expectedShards), len(balancer.ShardAssignments))
		return false
	}

	for i, shard := range expectedShards {
		if balancer.ShardAssignments[i][0] != shard[0] || balancer.ShardAssignments[i][1] != shard[1] {
			t.Errorf("Expected shard assignment %v, got %v", shard, balancer.ShardAssignments[i])
			return false
		}
	}

	return true
}

func TestAll(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		n1 := balancer.Balancer{
			ID:      uuid.New(),
			OwnerID: uuid.New(),
			Address: "http://example.com",
			ShardAssignments: [][2]int{
				{0, 1}, {1, 2},
			},
		}

		n2 := balancer.Balancer{
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

		conn, err := sql.Open("mysql", "root:juno@tcp(localhost:3306)/balancer_test?parseTime=true")
		if err != nil {
			t.Errorf("Error connecting to database: %s", err)
		}
		err = mysql.ExecuteMigrations(conn)
		if err != nil {
			t.Errorf("Error executing migrations: %s", err)
		}

		defer conn.Close()

		defer conn.Exec("DELETE FROM balancers WHERE id = ?", n1.ID)
		defer conn.Exec("DELETE FROM balancers WHERE id = ?", n2.ID)

		_, err = conn.Exec("INSERT INTO balancers (id, owner_id, address, shard_assignments) VALUES (?, ?, ?, ?)", n1.ID, n1.OwnerID, n1.Address, shardsJSON1)
		if err != nil {
			t.Errorf("Error inserting balancer: %s", err)
		}

		_, err = conn.Exec("INSERT INTO balancers (id, owner_id, address, shard_assignments) VALUES (?, ?, ?, ?)", n2.ID, n2.OwnerID, n2.Address, shardsJSON2)
		if err != nil {
			t.Errorf("Error inserting balancer: %s", err)
		}

		repo := New(conn)

		balancers, err := repo.All()

		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}

		if len(balancers) < 2 {
			t.Errorf("Expected 2 balancers, got %d", len(balancers))
		}
	})
}

func TestCreate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		n := balancer.Balancer{
			ID:      uuid.New(),
			OwnerID: uuid.New(),
			Address: "http://example.com",
			ShardAssignments: [][2]int{
				{0, 1}, {1, 2},
			},
		}

		conn, err := sql.Open("mysql", "root:juno@tcp(localhost:3306)/balancer_test?parseTime=true")
		if err != nil {
			t.Errorf("Error connecting to database: %s", err)
		}
		err = mysql.ExecuteMigrations(conn)
		if err != nil {
			t.Errorf("Error executing migrations: %s", err)
		}

		defer conn.Close()

		defer conn.Exec("DELETE FROM balancers WHERE id = ?", n.ID)

		repo := New(conn)

		err = repo.Create(&n)

		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}

		if !testBalancerMatches(t, conn, n.ID, n.OwnerID, n.Address, n.ShardAssignments) {
			t.Errorf("Balancer does not match")
		}
	})

	t.Run("error", func(t *testing.T) {
		n := balancer.Balancer{
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

		conn, err := sql.Open("mysql", "root:juno@tcp(localhost:3306)/balancer_test?parseTime=true")
		if err != nil {
			t.Errorf("Error connecting to database: %s", err)
		}
		err = mysql.ExecuteMigrations(conn)
		if err != nil {
			t.Errorf("Error executing migrations: %s", err)
		}

		defer conn.Close()

		_, err = conn.Exec("INSERT INTO balancers (id, owner_id, address, shard_assignments) VALUES (?, ?, ?, ?)", n.ID, n.OwnerID, n.Address, shardsJSON)

		if err != nil {
			t.Errorf("Error inserting balancer: %s", err)
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
		n := balancer.Balancer{
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

		conn, err := sql.Open("mysql", "root:juno@tcp(localhost:3306)/balancer_test?parseTime=true")
		if err != nil {
			t.Errorf("Error connecting to database: %s", err)
		}
		err = mysql.ExecuteMigrations(conn)
		if err != nil {
			t.Errorf("Error executing migrations: %s", err)
		}

		defer conn.Close()

		defer conn.Exec("DELETE FROM balancers WHERE id = ?", n.ID)

		_, err = conn.Exec("INSERT INTO balancers (id, owner_id, address, shard_assignments) VALUES (?, ?, ?, ?)", n.ID, n.OwnerID, n.Address, shardsJSON)

		if err != nil {
			t.Errorf("Error inserting balancer: %s", err)
		}

		repo := New(conn)

		balancer, err := repo.Get(n.ID)

		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}

		if balancer.ID != n.ID {
			t.Errorf("Expected ID %s, got %s", n.ID, balancer.ID)
		}

		if balancer.OwnerID != n.OwnerID {
			t.Errorf("Expected OwnerID %s, got %s", n.OwnerID, balancer.OwnerID)
		}

		if balancer.Address != n.Address {
			t.Errorf("Expected Address %s, got %s", n.Address, balancer.Address)
		}
	})

	t.Run("error", func(t *testing.T) {
		conn, err := sql.Open("mysql", "root:juno@tcp(localhost:3306)/balancer_test?parseTime=true")
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

		n1 := balancer.Balancer{
			ID:      uuid.New(),
			OwnerID: ownerID,
			Address: "http://example.com",
			ShardAssignments: [][2]int{
				{0, 1}, {1, 2},
			},
		}

		n2 := balancer.Balancer{
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

		conn, err := sql.Open("mysql", "root:juno@tcp(localhost:3306)/balancer_test?parseTime=true")
		if err != nil {
			t.Errorf("Error connecting to database: %s", err)
		}
		err = mysql.ExecuteMigrations(conn)
		if err != nil {
			t.Errorf("Error executing migrations: %s", err)
		}

		defer conn.Close()

		defer conn.Exec("DELETE FROM balancers WHERE id = ?", n1.ID)
		defer conn.Exec("DELETE FROM balancers WHERE id = ?", n2.ID)

		_, err = conn.Exec("INSERT INTO balancers (id, owner_id, address, shard_assignments) VALUES (?, ?, ?, ?)", n1.ID, n1.OwnerID, n1.Address, shardsJSON1)
		if err != nil {
			t.Errorf("Error inserting balancer: %s", err)
		}

		_, err = conn.Exec("INSERT INTO balancers (id, owner_id, address, shard_assignments) VALUES (?, ?, ?, ?)", n2.ID, n2.OwnerID, n2.Address, shardsJSON2)
		if err != nil {
			t.Errorf("Error inserting balancer: %s", err)
		}

		repo := New(conn)

		balancers, err := repo.ListByOwnerID(ownerID)

		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}

		if len(balancers) != 1 {
			t.Errorf("Expected 1 balancers, got %d", len(balancers))
		}

		if balancers[0].ID != n1.ID {
			t.Errorf("Expected ID %s, got %s", n1.ID, balancers[0].ID)
		}

		if balancers[0].OwnerID != n1.OwnerID {
			t.Errorf("Expected OwnerID %s, got %s", n1.OwnerID, balancers[0].OwnerID)
		}

		if balancers[0].Address != n1.Address {
			t.Errorf("Expected Address %s, got %s", n1.Address, balancers[0].Address)
		}
	})
}

func TestUpdate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		n := balancer.Balancer{
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

		conn, err := sql.Open("mysql", "root:juno@tcp(localhost:3306)/balancer_test?parseTime=true")
		if err != nil {
			t.Errorf("Error connecting to database: %s", err)
		}
		err = mysql.ExecuteMigrations(conn)
		if err != nil {
			t.Errorf("Error executing migrations: %s", err)
		}

		defer conn.Close()

		defer conn.Exec("DELETE FROM balancers WHERE id = ?", n.ID)

		_, err = conn.Exec("INSERT INTO balancers (id, owner_id, address, shard_assignments) VALUES (?, ?, ?, ?)", n.ID, n.OwnerID, n.Address, shardsJSON)

		if err != nil {
			t.Errorf("Error inserting balancer: %s", err)
		}

		repo := New(conn)

		n.Address = "http://example.org"

		err = repo.Update(&n)

		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}

		if !testBalancerMatches(t, conn, n.ID, n.OwnerID, n.Address, n.ShardAssignments) {
			t.Errorf("Balancer does not match")
		}
	})
}

func TestDelete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		n := balancer.Balancer{
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

		conn, err := sql.Open("mysql", "root:juno@tcp(localhost:3306)/balancer_test?parseTime=true")
		if err != nil {
			t.Errorf("Error connecting to database: %s", err)
		}
		err = mysql.ExecuteMigrations(conn)
		if err != nil {
			t.Errorf("Error executing migrations: %s", err)
		}

		defer conn.Close()

		_, err = conn.Exec("INSERT INTO balancers (id, owner_id, address, shard_assignments) VALUES (?, ?, ?, ?)", n.ID, n.OwnerID, n.Address, shardsJSON)

		if err != nil {
			t.Errorf("Error inserting balancer: %s", err)
		}

		repo := New(conn)

		err = repo.Delete(n.ID)

		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}

		row := conn.QueryRow("SELECT COUNT(*) FROM balancers WHERE id = ?", n.ID)

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
