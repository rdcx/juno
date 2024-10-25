package mysql

import (
	"database/sql"
	"log"
	"testing"

	"juno/pkg/api/extractor/filter"
	"juno/pkg/api/extractor/filter/migration/mysql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
)

func setupDB(t *testing.T) *sql.DB {
	db, err := sql.Open("mysql", "root:juno@tcp(localhost:3306)/filter_test?parseTime=true")

	if err != nil {
		log.Fatal(err)
	}

	err = mysql.ExecuteMigrations(db)

	if err != nil {
		t.Fatal(err)
	}

	return db
}

func TestGet(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		db := setupDB(t)
		defer db.Close()

		repo := New(db)

		id := uuid.New()

		f := &filter.Filter{
			ID:     id,
			UserID: uuid.New(),
			Name:   "test",
			Type:   filter.FilterTypeStringEquals,
			Value:  "test",
		}

		err := repo.Create(f)

		if err != nil {
			t.Fatal(err)
		}

		got, err := repo.Get(id)

		if err != nil {
			t.Fatal(err)
		}

		if got.ID != f.ID {
			t.Errorf("expected id to be %v, got %v", f.ID, got.ID)
		}

		if got.UserID != f.UserID {
			t.Errorf("expected user id to be %v, got %v", f.UserID, got.UserID)
		}

		if got.Name != f.Name {
			t.Errorf("expected name to be %v, got %v", f.Name, got.Name)
		}

		if got.Type != f.Type {
			t.Errorf("expected type to be %v, got %v", f.Type, got.Type)
		}

		if got.Value != f.Value {
			t.Errorf("expected value to be %v, got %v", f.Value, got.Value)
		}
	})

	t.Run("not found", func(t *testing.T) {
		db := setupDB(t)
		defer db.Close()

		repo := New(db)

		id := uuid.New()

		_, err := repo.Get(id)

		if err != filter.ErrNotFound {
			t.Errorf("expected error to be %v, got %v", filter.ErrNotFound, err)
		}
	})
}

func TestCreate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		db := setupDB(t)
		defer db.Close()

		repo := New(db)

		f := &filter.Filter{
			ID:     uuid.New(),
			UserID: uuid.New(),
			Name:   "test",
			Type:   filter.FilterTypeStringEquals,
			Value:  "test",
		}

		err := repo.Create(f)

		if err != nil {
			t.Fatal(err)
		}
	})
}

func TestUpdate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		db := setupDB(t)
		defer db.Close()

		repo := New(db)

		id := uuid.New()

		f := &filter.Filter{
			ID:     id,
			UserID: uuid.New(),
			Name:   "test",
			Type:   filter.FilterTypeStringEquals,
			Value:  "test",
		}

		err := repo.Create(f)

		if err != nil {
			t.Fatal(err)
		}

		f.Name = "test2"
		f.Type = filter.FilterTypeStringContains
		f.Value = "test2"

		err = repo.Update(f)

		if err != nil {
			t.Fatal(err)
		}

		got, err := repo.Get(id)

		if err != nil {
			t.Fatal(err)
		}

		if got.Name != f.Name {
			t.Errorf("expected name to be %v, got %v", f.Name, got.Name)
		}

		if got.Type != f.Type {
			t.Errorf("expected type to be %v, got %v", f.Type, got.Type)
		}

		if got.Value != f.Value {
			t.Errorf("expected value to be %v, got %v", f.Value, got.Value)
		}
	})
}

func TestDelete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		db := setupDB(t)
		defer db.Close()

		repo := New(db)

		id := uuid.New()

		f := &filter.Filter{
			ID:     id,
			UserID: uuid.New(),
			Name:   "test",
			Type:   filter.FilterTypeStringEquals,
			Value:  "test",
		}

		err := repo.Create(f)

		if err != nil {
			t.Fatal(err)
		}

		err = repo.Delete(id)

		if err != nil {
			t.Fatal(err)
		}

		_, err = repo.Get(id)

		if err != filter.ErrNotFound {
			t.Errorf("expected error to be %v, got %v", filter.ErrNotFound, err)
		}
	})
}

func TestListByUserID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		db := setupDB(t)
		defer db.Close()

		repo := New(db)

		userID := uuid.New()

		f1 := &filter.Filter{
			ID:     uuid.New(),
			UserID: userID,
			Name:   "test",
			Type:   filter.FilterTypeStringEquals,
			Value:  "test",
		}

		err := repo.Create(f1)

		if err != nil {
			t.Fatal(err)
		}

		f2 := &filter.Filter{
			ID:     uuid.New(),
			UserID: uuid.New(),
			Name:   "test2",
			Type:   filter.FilterTypeStringContains,
			Value:  "test2",
		}

		err = repo.Create(f2)

		if err != nil {
			t.Fatal(err)
		}

		filters, err := repo.ListByUserID(userID)

		if err != nil {
			t.Fatal(err)
		}

		if len(filters) != 1 {
			t.Errorf("expected length to be 2, got %d", len(filters))
		}

		if filters[0].ID != f1.ID {
			t.Errorf("expected id to be %v, got %v", f1.ID, filters[0].ID)
		}
	})
}
