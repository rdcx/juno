package mysql

import (
	"database/sql"
	"log"
	"testing"

	"juno/pkg/api/extractor/field"
	"juno/pkg/api/extractor/field/migration/mysql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
)

func setupDB(t *testing.T) *sql.DB {
	db, err := sql.Open("mysql", "root:juno@tcp(localhost:3306)/field_test?parseTime=true")

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

		f := &field.Field{
			ID:         id,
			UserID:     uuid.New(),
			SelectorID: uuid.New(),
			Name:       "test",
			Type:       field.FieldTypeFloat,
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

		if got.SelectorID != f.SelectorID {
			t.Errorf("expected selector id to be %v, got %v", f.SelectorID, got.SelectorID)
		}

		if got.Name != f.Name {
			t.Errorf("expected name to be %v, got %v", f.Name, got.Name)
		}

		if got.Type != f.Type {
			t.Errorf("expected type to be %v, got %v", f.Type, got.Type)
		}
	})

	t.Run("not found", func(t *testing.T) {
		db := setupDB(t)
		defer db.Close()

		repo := New(db)

		id := uuid.New()

		_, err := repo.Get(id)

		if err != field.ErrNotFound {
			t.Errorf("expected error to be %v, got %v", field.ErrNotFound, err)
		}
	})
}

func TestCreate(t *testing.T) {
	db := setupDB(t)
	defer db.Close()

	repo := New(db)

	id := uuid.New()

	f := &field.Field{
		ID:         id,
		UserID:     uuid.New(),
		SelectorID: uuid.New(),
		Name:       "test",
		Type:       field.FieldTypeFloat,
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

	if got.SelectorID != f.SelectorID {
		t.Errorf("expected selector id to be %v, got %v", f.SelectorID, got.SelectorID)
	}

	if got.Name != f.Name {
		t.Errorf("expected name to be %v, got %v", f.Name, got.Name)
	}

	if got.Type != f.Type {
		t.Errorf("expected type to be %v, got %v", f.Type, got.Type)
	}
}

func TestListByUserID(t *testing.T) {
	db := setupDB(t)
	defer db.Close()

	repo := New(db)

	userID := uuid.New()

	f1 := &field.Field{
		ID:         uuid.New(),
		UserID:     userID,
		SelectorID: uuid.New(),
		Name:       "test1",
		Type:       field.FieldTypeFloat,
	}

	err := repo.Create(f1)

	if err != nil {
		t.Fatal(err)
	}

	f2 := &field.Field{
		ID:         uuid.New(),
		UserID:     uuid.New(),
		SelectorID: uuid.New(),
		Name:       "test2",
		Type:       field.FieldTypeFloat,
	}

	err = repo.Create(f2)

	if err != nil {
		t.Fatal(err)
	}

	list, err := repo.ListByUserID(userID)

	if err != nil {
		t.Fatal(err)
	}

	if len(list) != 1 {
		t.Errorf("expected list length to be 2, got %d", len(list))
	}

	if list[0].ID != f1.ID {
		t.Errorf("expected id to be %v, got %v", f1.ID, list[0].ID)
	}
}

func TestListBySelectorID(t *testing.T) {
	db := setupDB(t)
	defer db.Close()

	repo := New(db)

	selectorID := uuid.New()

	f1 := &field.Field{
		ID:         uuid.New(),
		UserID:     uuid.New(),
		SelectorID: selectorID,
		Name:       "test1",
		Type:       field.FieldTypeFloat,
	}

	err := repo.Create(f1)

	if err != nil {
		t.Fatal(err)
	}

	f2 := &field.Field{
		ID:         uuid.New(),
		UserID:     uuid.New(),
		SelectorID: uuid.New(),
		Name:       "test2",
		Type:       field.FieldTypeFloat,
	}

	err = repo.Create(f2)

	if err != nil {
		t.Fatal(err)
	}

	list, err := repo.ListBySelectorID(selectorID)

	if err != nil {
		t.Fatal(err)
	}

	if len(list) != 1 {
		t.Errorf("expected list length to be 2, got %d", len(list))
	}

	if list[0].ID != f1.ID {
		t.Errorf("expected id to be %v, got %v", f1.ID, list[0].ID)
	}
}

func TestUpdate(t *testing.T) {
	db := setupDB(t)
	defer db.Close()

	repo := New(db)

	id := uuid.New()

	f := &field.Field{
		ID:         id,
		UserID:     uuid.New(),
		SelectorID: uuid.New(),
		Name:       "test",
		Type:       field.FieldTypeFloat,
	}

	err := repo.Create(f)

	if err != nil {
		t.Fatal(err)
	}

	f.Type = field.FieldTypeInteger

	err = repo.Update(f)

	if err != nil {
		t.Fatal(err)
	}

	got, err := repo.Get(id)

	if err != nil {
		t.Fatal(err)
	}

	if got.Type != field.FieldTypeInteger {
		t.Errorf("expected type to be %v, got %v", field.FieldTypeInteger, got.Type)
	}
}

func TestDelete(t *testing.T) {
	db := setupDB(t)
	defer db.Close()

	repo := New(db)

	id := uuid.New()

	f := &field.Field{
		ID:         id,
		UserID:     uuid.New(),
		SelectorID: uuid.New(),
		Name:       "test",
		Type:       field.FieldTypeFloat,
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

	if err != field.ErrNotFound {
		t.Errorf("expected error to be %v, got %v", field.ErrNotFound, err)
	}
}
