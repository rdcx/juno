package mysql

import (
	"database/sql"
	"log"
	"testing"

	"juno/pkg/api/extractor/selector"
	"juno/pkg/api/extractor/selector/migration/mysql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
)

func setupDB(t *testing.T) *sql.DB {
	db, err := sql.Open("mysql", "root:juno@tcp(localhost:3306)/selector_test?parseTime=true")

	if err != nil {
		log.Fatal(err)
	}

	err = mysql.ExecuteMigrations(db)

	if err != nil {
		log.Fatal(err)
	}

	return db
}

func TestCreate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		db := setupDB(t)
		repo := New(db)

		j := &selector.Selector{
			ID: uuid.New(),
		}

		err := repo.Create(j)

		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}

		var check selector.Selector

		err = db.QueryRow("SELECT id, user_id, name, value, visibility, created_at, updated_at FROM selectors WHERE id = ?", j.ID).Scan(&check.ID, &check.UserID, &check.Name, &check.Value, &check.Visibility, &check.CreatedAt, &check.UpdatedAt)

		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}

		if check.ID != j.ID {
			t.Errorf("Expected %s, got %s", j.ID, check.ID)
		}

		if check.UserID != j.UserID {
			t.Errorf("Expected %s, got %s", j.UserID, check.UserID)
		}
	})
}

func TestGet(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		db := setupDB(t)
		repo := New(db)

		id := uuid.New()
		_, err := db.Exec("INSERT INTO selectors (id, user_id, name, value, visibility) VALUES (?, ?, ?, ?, ?)", id, uuid.New(), "name", "value", "visibility")

		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}

		j, err := repo.Get(id)

		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}

		if j.ID != id {
			t.Errorf("Expected %s, got %s", id, j.ID)
		}
	})

	t.Run("not found", func(t *testing.T) {
		db := setupDB(t)
		repo := New(db)

		id := uuid.New()

		_, err := repo.Get(id)

		if err == nil {
			t.Errorf("Expected error, got nil")
		}

		if err != selector.ErrNotFound {
			t.Errorf("Expected ErrNotFound, got %v", err)
		}
	})
}

func TestListByUserID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		db := setupDB(t)
		repo := New(db)

		userID := uuid.New()
		_, err := db.Exec("INSERT INTO selectors (id, user_id, name, value, visibility) VALUES (?, ?, ?, ?, ?)", uuid.New(), userID, "name", "value", "visibility")

		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}

		selectors, err := repo.ListByUserID(userID)

		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}

		if len(selectors) != 1 {
			t.Errorf("Expected 1, got %d", len(selectors))
		}
	})
}

func TestUpdate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		db := setupDB(t)
		repo := New(db)

		id := uuid.New()
		_, err := db.Exec("INSERT INTO selectors (id, user_id, name, value, visibility) VALUES (?, ?, ?, ?, ?)", id, uuid.New(), "name", "value", "visibility")

		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}

		j := &selector.Selector{
			ID:         id,
			Name:       "new name",
			Value:      "new value",
			Visibility: "new visibility",
		}

		err = repo.Update(j)

		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}

		var check selector.Selector

		err = db.QueryRow("SELECT id, user_id, name, value, visibility, created_at, updated_at FROM selectors WHERE id = ?", id).Scan(&check.ID, &check.UserID, &check.Name, &check.Value, &check.Visibility, &check.CreatedAt, &check.UpdatedAt)

		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}

		if check.Name != j.Name {
			t.Errorf("Expected %s, got %s", j.Name, check.Name)
		}

		if check.Value != j.Value {
			t.Errorf("Expected %s, got %s", j.Value, check.Value)
		}

		if check.Visibility != j.Visibility {
			t.Errorf("Expected %s, got %s", j.Visibility, check.Visibility)
		}
	})
}

func TestDelete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		db := setupDB(t)
		repo := New(db)

		id := uuid.New()
		_, err := db.Exec("INSERT INTO selectors (id, user_id, name, value, visibility) VALUES (?, ?, ?, ?, ?)", id, uuid.New(), "name", "value", "visibility")

		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}

		err = repo.Delete(id)

		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}

		var check selector.Selector

		err = db.QueryRow("SELECT id, user_id, name, value, visibility, created_at, updated_at FROM selectors WHERE id = ?", id).Scan(&check.ID, &check.UserID, &check.Name, &check.Value, &check.Visibility, &check.CreatedAt, &check.UpdatedAt)

		if err == nil {
			t.Errorf("Expected error, got nil")
		}
	})
}
