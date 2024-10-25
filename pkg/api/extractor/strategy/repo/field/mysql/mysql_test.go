package mysql

import (
	"database/sql"
	"juno/pkg/api/extractor/strategy/migration/mysql"
	"log"
	"testing"

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

func TestRepo(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		db := setupDB(t)
		defer db.Close()

		repo := New(db)

		id := uuid.New()
		fieldID := uuid.New()

		err := repo.AddField(id, fieldID)

		if err != nil {
			t.Fatal(err)
		}

		ids, err := repo.ListFieldIDs(id)

		if err != nil {
			t.Fatal(err)
		}

		if len(ids) != 1 {
			t.Errorf("expected 1, got %d", len(ids))
		}

		if ids[0] != fieldID {
			t.Errorf("expected %s, got %s", fieldID, ids[0])
		}

		err = repo.RemoveField(id, fieldID)

		if err != nil {
			t.Fatal(err)
		}

		ids, err = repo.ListFieldIDs(id)

		if err != nil {
			t.Fatal(err)
		}

		if len(ids) != 0 {
			t.Errorf("expected 0, got %d", len(ids))
		}
	})
}
