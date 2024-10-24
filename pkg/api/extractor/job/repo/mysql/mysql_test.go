package mysql

import (
	"database/sql"
	"log"
	"testing"

	"juno/pkg/api/extractor/job"
	"juno/pkg/api/extractor/job/migration/mysql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
)

func setupDB(t *testing.T) *sql.DB {
	db, err := sql.Open("mysql", "root:juno@tcp(localhost:3306)/extractor_test?parseTime=true")

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

		j := &job.Job{
			ID: uuid.New(),
		}

		err := repo.Create(j)

		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}

		var check job.Job

		err = db.QueryRow("SELECT id, user_id, strategy_id, status, created_at, updated_at FROM jobs WHERE id = ?", j.ID).Scan(&check.ID, &check.UserID, &check.StrategyID, &check.Status, &check.CreatedAt, &check.UpdatedAt)

		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}

		if check.ID != j.ID {
			t.Errorf("Expected %s, got %s", j.ID, check.ID)
		}

		if check.UserID != j.UserID {
			t.Errorf("Expected %s, got %s", j.UserID, check.UserID)
		}

		if check.Status != j.Status {
			t.Errorf("Expected %s, got %s", j.Status, check.Status)
		}

		if check.StrategyID != j.StrategyID {
			t.Errorf("Expected %s, got %s", j.StrategyID, check.StrategyID)
		}

		if check.Status != j.Status {
			t.Errorf("Expected %s, got %s", j.Status, check.Status)
		}
	})
}

func TestGet(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		db := setupDB(t)
		repo := New(db)

		j := &job.Job{
			ID: uuid.New(),
		}

		err := repo.Create(j)

		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}

		check, err := repo.Get(j.ID)

		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}

		if check.ID != j.ID {
			t.Errorf("Expected %s, got %s", j.ID, check.ID)
		}

		if check.UserID != j.UserID {
			t.Errorf("Expected %s, got %s", j.UserID, check.UserID)
		}

		if check.Status != j.Status {
			t.Errorf("Expected %s, got %s", j.Status, check.Status)
		}

		if check.StrategyID != j.StrategyID {
			t.Errorf("Expected %s, got %s", j.StrategyID, check.StrategyID)
		}
	})
}

func TestListByUserID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		db := setupDB(t)
		repo := New(db)

		userID := uuid.New()

		j1 := &job.Job{
			ID:     uuid.New(),
			UserID: userID,
		}

		j2 := &job.Job{
			ID:     uuid.New(),
			UserID: uuid.New(),
		}

		err := repo.Create(j1)

		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}

		err = repo.Create(j2)

		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}

		list, err := repo.ListByUserID(userID)

		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}

		if len(list) != 1 {
			t.Errorf("Expected 2, got %d", len(list))
		}

		if list[0].ID != j1.ID {
			t.Errorf("Expected %s, got %s", j1.ID, list[0].ID)
		}
	})
}

func TestUpdate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		db := setupDB(t)
		repo := New(db)

		j := &job.Job{
			ID: uuid.New(),
		}

		err := repo.Create(j)

		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}

		j.Status = job.CompletedStatus

		err = repo.Update(j)

		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}

		var check job.Job

		err = db.QueryRow("SELECT id, user_id, strategy_id, status, created_at, updated_at FROM jobs WHERE id = ?", j.ID).Scan(&check.ID, &check.UserID, &check.StrategyID, &check.Status, &check.CreatedAt, &check.UpdatedAt)

		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}

		if check.Status != j.Status {
			t.Errorf("Expected %s, got %s", j.Status, check.Status)
		}
	})
}
