package mysql

import (
	"database/sql"
	"log"
	"testing"

	"juno/pkg/api/extractor/strategy"
	"juno/pkg/api/extractor/strategy/migration/mysql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
)

func setupDB(t *testing.T) *sql.DB {
	db, err := sql.Open("mysql", "root:juno@tcp(localhost:3306)/strategy_test?parseTime=true")

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

		f := &strategy.Strategy{
			ID:     id,
			UserID: uuid.New(),
			Name:   "test",
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
	})

	t.Run("not found", func(t *testing.T) {
		db := setupDB(t)
		defer db.Close()

		repo := New(db)

		id := uuid.New()

		_, err := repo.Get(id)

		if err != strategy.ErrNotFound {
			t.Errorf("expected error to be %v, got %v", strategy.ErrNotFound, err)
		}
	})
}

func TestCreate(t *testing.T) {
	db := setupDB(t)
	defer db.Close()

	repo := New(db)

	f := &strategy.Strategy{
		ID:     uuid.New(),
		UserID: uuid.New(),
		Name:   "test",
	}

	err := repo.Create(f)

	if err != nil {
		t.Fatal(err)
	}

	got, err := repo.Get(f.ID)

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
}

func TestListByUserID(t *testing.T) {
	db := setupDB(t)
	defer db.Close()

	repo := New(db)

	userID := uuid.New()

	f := &strategy.Strategy{
		ID:     uuid.New(),
		UserID: userID,
		Name:   "test",
	}

	err := repo.Create(f)

	if err != nil {
		t.Fatal(err)
	}

	list, err := repo.ListByUserID(userID)

	if err != nil {
		t.Fatal(err)
	}

	if len(list) != 1 {
		t.Errorf("expected 1, got %d", len(list))
	}

	if list[0].ID != f.ID {
		t.Errorf("expected %s, got %s", f.ID, list[0].ID)
	}
}

func TestUpdate(t *testing.T) {
	db := setupDB(t)
	defer db.Close()

	repo := New(db)

	f := &strategy.Strategy{
		ID:     uuid.New(),
		UserID: uuid.New(),
		Name:   "test",
	}

	err := repo.Create(f)

	if err != nil {
		t.Fatal(err)
	}

	f.Name = "updated"

	err = repo.Update(f)

	if err != nil {
		t.Fatal(err)
	}

	got, err := repo.Get(f.ID)

	if err != nil {
		t.Fatal(err)
	}

	if got.Name != f.Name {
		t.Errorf("expected name to be %v, got %v", f.Name, got.Name)
	}
}

func TestDelete(t *testing.T) {
	db := setupDB(t)
	defer db.Close()

	repo := New(db)

	f := &strategy.Strategy{
		ID:     uuid.New(),
		UserID: uuid.New(),
		Name:   "test",
	}

	err := repo.Create(f)

	if err != nil {
		t.Fatal(err)
	}

	err = repo.Delete(f.ID)

	if err != nil {
		t.Fatal(err)
	}

	_, err = repo.Get(f.ID)

	if err != strategy.ErrNotFound {
		t.Errorf("expected error to be %v, got %v", strategy.ErrNotFound, err)
	}
}
