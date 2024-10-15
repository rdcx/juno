package mysql

import (
	"database/sql"
	"juno/pkg/api/node"
	"juno/pkg/api/node/migration/mysql"
	"testing"

	_ "github.com/go-sql-driver/mysql"

	"github.com/google/uuid"
)

func testNodeMatches(t *testing.T, conn *sql.DB, id, ownerID uuid.UUID, address string) bool {
	sqlCheck := "SELECT id, owner_id, address FROM nodes WHERE id = ?"

	row := conn.QueryRow(sqlCheck, id)

	var node node.Node
	err := row.Scan(&node.ID, &node.OwnerID, &node.Address)
	if err != nil {
		t.Errorf("Error getting node: %s", err)
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

	return true
}

func TestCreate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		n := node.Node{
			ID:      uuid.New(),
			OwnerID: uuid.New(),
			Address: "http://example.com",
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

		if !testNodeMatches(t, conn, n.ID, n.OwnerID, n.Address) {
			t.Errorf("Node does not match")
		}
	})

	t.Run("error", func(t *testing.T) {
		n := node.Node{
			ID:      uuid.New(),
			OwnerID: uuid.New(),
			Address: "http://example.com",
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

		_, err = conn.Exec("INSERT INTO nodes (id, owner_id, address) VALUES (?, ?, ?)", n.ID, n.OwnerID, n.Address)

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

		_, err = conn.Exec("INSERT INTO nodes (id, owner_id, address) VALUES (?, ?, ?)", n.ID, n.OwnerID, n.Address)

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
		}

		n2 := node.Node{
			ID:      uuid.New(),
			OwnerID: ownerID,
			Address: "http://example.org",
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

		_, err = conn.Exec("INSERT INTO nodes (id, owner_id, address) VALUES (?, ?, ?)", n1.ID, n1.OwnerID, n1.Address)
		if err != nil {
			t.Errorf("Error inserting node: %s", err)
		}

		_, err = conn.Exec("INSERT INTO nodes (id, owner_id, address) VALUES (?, ?, ?)", n2.ID, n2.OwnerID, n2.Address)
		if err != nil {
			t.Errorf("Error inserting node: %s", err)
		}

		repo := New(conn)

		nodes, err := repo.ListByOwnerID(ownerID)

		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}

		if len(nodes) != 2 {
			t.Errorf("Expected 2 nodes, got %d", len(nodes))
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

		if nodes[1].ID != n2.ID {
			t.Errorf("Expected ID %s, got %s", n2.ID, nodes[1].ID)
		}

		if nodes[1].OwnerID != n2.OwnerID {
			t.Errorf("Expected OwnerID %s, got %s", n2.OwnerID, nodes[1].OwnerID)
		}

		if nodes[1].Address != n2.Address {
			t.Errorf("Expected Address %s, got %s", n2.Address, nodes[1].Address)
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

		_, err = conn.Exec("INSERT INTO nodes (id, owner_id, address) VALUES (?, ?, ?)", n.ID, n.OwnerID, n.Address)

		if err != nil {
			t.Errorf("Error inserting node: %s", err)
		}

		repo := New(conn)

		n.Address = "http://example.org"

		err = repo.Update(&n)

		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}

		if !testNodeMatches(t, conn, n.ID, n.OwnerID, n.Address) {
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

		_, err = conn.Exec("INSERT INTO nodes (id, owner_id, address) VALUES (?, ?, ?)", n.ID, n.OwnerID, n.Address)

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
