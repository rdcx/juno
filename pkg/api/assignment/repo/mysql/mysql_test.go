package mysql

import (
	"database/sql"
	"juno/pkg/api/assignment"
	"juno/pkg/api/assignment/migration/mysql"
	"testing"

	_ "github.com/go-sql-driver/mysql"

	"github.com/google/uuid"
)

func testDB(t *testing.T) *sql.DB {
	db, err := sql.Open("mysql", "root:juno@tcp(localhost:3306)/assignment_test?parseTime=true")
	if err != nil {
		t.Fatalf("could not connect to mysql: %v", err)
	}

	err = mysql.ExecuteMigrations(db)
	if err != nil {
		t.Fatalf("could not execute migrations: %v", err)
	}

	return db
}

func TestGet(t *testing.T) {
	t.Run("returns assignment by ID", func(t *testing.T) {
		db := testDB(t)
		repo := New(db)
		a := &assignment.Assignment{
			ID:       uuid.New(),
			OwnerID:  uuid.New(),
			EntityID: uuid.New(),
			Offset:   0,
			Length:   10,
		}

		repo.Create(a)
		defer db.Exec("DELETE FROM assignments WHERE id = ?", a.ID)

		result, err := repo.Get(a.ID)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if result.ID != a.ID {
			t.Errorf("expected assignment ID %v, got %v", a.ID, result.ID)
		}

		if result.OwnerID != a.OwnerID {
			t.Errorf("expected assignment OwnerID %v, got %v", a.OwnerID, result.OwnerID)
		}

		if result.EntityID != a.EntityID {
			t.Errorf("expected assignment EntityID %v, got %v", a.EntityID, result.EntityID)
		}

		if result.Offset != a.Offset {
			t.Errorf("expected assignment Offset %v, got %v", a.Offset, result.Offset)
		}

		if result.Length != a.Length {
			t.Errorf("expected assignment Length %v, got %v", a.Length, result.Length)
		}

	})

	t.Run("returns error if assignment not found", func(t *testing.T) {
		db := testDB(t)
		repo := New(db)
		id := uuid.New()

		_, err := repo.Get(id)

		if err != assignment.ErrNotFound {
			t.Errorf("expected error %v, got %v", assignment.ErrNotFound, err)
		}
	})
}

func TestListByEntityID(t *testing.T) {
	t.Run("returns assignments for entity ID", func(t *testing.T) {
		db := testDB(t)
		repo := New(db)
		entityID := uuid.New()
		assignments := []*assignment.Assignment{
			{
				ID:       uuid.New(),
				OwnerID:  uuid.New(),
				EntityID: entityID,
				Offset:   0,
				Length:   10,
			},
			{
				ID:       uuid.New(),
				OwnerID:  uuid.New(),
				EntityID: uuid.New(),
				Offset:   10,
				Length:   10,
			},
		}

		for _, a := range assignments {
			err := repo.Create(a)
			defer db.Exec("DELETE FROM assignments WHERE id = ?", a.ID)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		}

		result, err := repo.ListByEntityID(entityID)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(result) != 1 {
			t.Fatalf("expected 1 assignments, got %d", len(result))
		}

		for i, a := range result {
			if a.ID != assignments[i].ID {
				t.Errorf("expected assignment ID %v, got %v", assignments[i].ID, a.ID)
			}
		}
	})
}

func TestCreate(t *testing.T) {
	t.Run("creates assignment", func(t *testing.T) {
		db := testDB(t)
		repo := New(db)
		a := &assignment.Assignment{
			ID:       uuid.New(),
			OwnerID:  uuid.New(),
			EntityID: uuid.New(),
			Offset:   0,
			Length:   10,
		}

		defer db.Exec("DELETE FROM assignments WHERE id = ?", a.ID)
		err := repo.Create(a)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		var check assignment.Assignment

		err = db.QueryRow("SELECT id, owner_id, entity_id, offset, length FROM assignments WHERE id = ?", a.ID).Scan(&check.ID, &check.OwnerID, &check.EntityID, &check.Offset, &check.Length)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if check.ID != a.ID {
			t.Errorf("expected assignment ID %v, got %v", a.ID, check.ID)
		}

		if check.OwnerID != a.OwnerID {
			t.Errorf("expected assignment OwnerID %v, got %v", a.OwnerID, check.OwnerID)
		}

		if check.EntityID != a.EntityID {
			t.Errorf("expected assignment EntityID %v, got %v", a.EntityID, check.EntityID)
		}

		if check.Offset != a.Offset {
			t.Errorf("expected assignment Offset %v, got %v", a.Offset, check.Offset)
		}

		if check.Length != a.Length {
			t.Errorf("expected assignment Length %v, got %v", a.Length, check.Length)
		}
	})
}

func TestUpdate(t *testing.T) {
	t.Run("updates assignment", func(t *testing.T) {
		db := testDB(t)
		repo := New(db)
		a := &assignment.Assignment{
			ID:       uuid.New(),
			OwnerID:  uuid.New(),
			EntityID: uuid.New(),
			Offset:   0,
			Length:   10,
		}

		repo.Create(a)
		defer db.Exec("DELETE FROM assignments WHERE id = ?", a.ID)

		a.Offset = 10
		a.Length = 20

		err := repo.Update(a)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		var check assignment.Assignment

		err = db.QueryRow("SELECT id, owner_id, entity_id, offset, length FROM assignments WHERE id = ?", a.ID).Scan(&check.ID, &check.OwnerID, &check.EntityID, &check.Offset, &check.Length)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if check.ID != a.ID {
			t.Errorf("expected assignment ID %v, got %v", a.ID, check.ID)
		}

		if check.OwnerID != a.OwnerID {
			t.Errorf("expected assignment OwnerID %v, got %v", a.OwnerID, check.OwnerID)
		}

		if check.EntityID != a.EntityID {
			t.Errorf("expected assignment EntityID %v, got %v", a.EntityID, check.EntityID)
		}

		if check.Offset != a.Offset {
			t.Errorf("expected assignment Offset %v, got %v", a.Offset, check.Offset)
		}

		if check.Length != a.Length {
			t.Errorf("expected assignment Length %v, got %v", a.Length, check.Length)
		}

	})
}

func TestDelete(t *testing.T) {
	t.Run("deletes assignment", func(t *testing.T) {
		db := testDB(t)
		repo := New(db)
		a := &assignment.Assignment{
			ID:       uuid.New(),
			OwnerID:  uuid.New(),
			EntityID: uuid.New(),
			Offset:   0,
			Length:   10,
		}

		repo.Create(a)

		err := repo.Delete(a.ID)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		var count int
		err = db.QueryRow("SELECT COUNT(*) FROM assignments WHERE id = ?", a.ID).Scan(&count)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if count != 0 {
			t.Fatalf("expected 0 assignments, got %d", count)
		}
	})
}
