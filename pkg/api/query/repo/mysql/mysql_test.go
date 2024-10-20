package mysql

import (
	"database/sql"
	"encoding/json"
	"log"
	"testing"

	"juno/pkg/api/query"
	"juno/pkg/api/query/migration/mysql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
)

func setupDB(t *testing.T) *sql.DB {
	db, err := sql.Open("mysql", "root:juno@tcp(localhost:3306)/query_test?parseTime=true")

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
	db := setupDB(t)
	repo := New(db)

	q := &query.Query{
		ID:                uuid.New(),
		UserID:            uuid.New(),
		QueryType:         query.BasicQueryType,
		BasicQueryVersion: "1.0",
		BasicQuery: &query.BasicQuery{
			Title: &query.StringMatch{
				Value:     "test",
				MatchType: query.ExactStringMatch,
			},
		},
		Status: query.PendingStatus,
	}
	defer db.Exec("DELETE FROM queries WHERE id = ?", q.ID)

	err := repo.Create(q)

	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}

	var check query.Query

	var basicQuery string

	err = db.QueryRow("SELECT id, user_id, status, query_type, basic_query_version, basic_query FROM queries WHERE id = ?", q.ID).Scan(&check.ID, &check.UserID, &check.Status, &check.QueryType, &check.BasicQueryVersion, &basicQuery)

	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}

	check.BasicQuery = &query.BasicQuery{}

	err = json.Unmarshal([]byte(basicQuery), check.BasicQuery)

	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}

	if check.ID != q.ID {
		t.Errorf("Expected %s, got %s", q.ID, check.ID)
	}

	if check.UserID != q.UserID {
		t.Errorf("Expected %s, got %s", q.UserID, check.UserID)
	}

	if check.Status != q.Status {
		t.Errorf("Expected %s, got %s", q.Status, check.Status)
	}

	if check.QueryType != q.QueryType {
		t.Errorf("Expected %s, got %s", q.QueryType, check.QueryType)
	}

	if check.BasicQueryVersion != q.BasicQueryVersion {
		t.Errorf("Expected %s, got %s", q.BasicQueryVersion, check.BasicQueryVersion)
	}

	if check.BasicQuery.Title.Value != q.BasicQuery.Title.Value {
		t.Errorf("Expected %s, got %s", q.BasicQuery.Title.Value, check.BasicQuery.Title.Value)
	}

	if check.BasicQuery.Title.MatchType != q.BasicQuery.Title.MatchType {
		t.Errorf("Expected %s, got %s", q.BasicQuery.Title.MatchType, check.BasicQuery.Title.MatchType)
	}
}

func TestGet(t *testing.T) {
	db := setupDB(t)
	repo := New(db)

	q := &query.Query{
		ID:                uuid.New(),
		UserID:            uuid.New(),
		QueryType:         query.BasicQueryType,
		BasicQueryVersion: "1.0",
		BasicQuery: &query.BasicQuery{
			Title: &query.StringMatch{
				Value:     "test",
				MatchType: query.ExactStringMatch,
			},
		},
		Status: query.PendingStatus,
	}
	defer db.Exec("DELETE FROM queries WHERE id = ?", q.ID)

	err := repo.Create(q)

	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}

	check, err := repo.Get(q.ID)

	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}

	if check.ID != q.ID {
		t.Errorf("Expected %s, got %s", q.ID, check.ID)
	}

	if check.UserID != q.UserID {
		t.Errorf("Expected %s, got %s", q.UserID, check.UserID)
	}

	if check.Status != q.Status {
		t.Errorf("Expected %s, got %s", q.Status, check.Status)
	}

	if check.QueryType != q.QueryType {
		t.Errorf("Expected %s, got %s", q.QueryType, check.QueryType)
	}

	if check.BasicQueryVersion != q.BasicQueryVersion {
		t.Errorf("Expected %s, got %s", q.BasicQueryVersion, check.BasicQueryVersion)
	}

	if check.BasicQuery.Title.Value != q.BasicQuery.Title.Value {
		t.Errorf("Expected %s, got %s", q.BasicQuery.Title.Value, check.BasicQuery.Title.Value)
	}

	if check.BasicQuery.Title.MatchType != q.BasicQuery.Title.MatchType {
		t.Errorf("Expected %s, got %s", q.BasicQuery.Title.MatchType, check.BasicQuery.Title.MatchType)
	}
}

func TestListByUserID(t *testing.T) {
	db := setupDB(t)
	repo := New(db)

	userID := uuid.New()

	q := &query.Query{
		ID:                uuid.New(),
		UserID:            userID,
		QueryType:         query.BasicQueryType,
		BasicQueryVersion: "1.0",
		BasicQuery: &query.BasicQuery{
			Title: &query.StringMatch{
				Value:     "test",
				MatchType: query.ExactStringMatch,
			},
		},
		Status: query.PendingStatus,
	}
	defer db.Exec("DELETE FROM queries WHERE id = ?", q.ID)

	err := repo.Create(q)

	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}

	queries, err := repo.ListByUserID(userID)

	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}

	if len(queries) != 1 {
		t.Errorf("Expected 1, got %d", len(queries))
	}

	check := queries[0]

	if check.ID != q.ID {
		t.Errorf("Expected %s, got %s", q.ID, check.ID)
	}

	if check.UserID != q.UserID {
		t.Errorf("Expected %s, got %s", q.UserID, check.UserID)
	}

	if check.Status != q.Status {
		t.Errorf("Expected %s, got %s", q.Status, check.Status)
	}

	if check.QueryType != q.QueryType {
		t.Errorf("Expected %s, got %s", q.QueryType, check.QueryType)
	}

	if check.BasicQueryVersion != q.BasicQueryVersion {
		t.Errorf("Expected %s, got %s", q.BasicQueryVersion, check.BasicQueryVersion)
	}

	if check.BasicQuery.Title.Value != q.BasicQuery.Title.Value {
		t.Errorf("Expected %s, got %s", q.BasicQuery.Title.Value, check.BasicQuery.Title.Value)
	}

	if check.BasicQuery.Title.MatchType != q.BasicQuery.Title.MatchType {
		t.Errorf("Expected %s, got %s", q.BasicQuery.Title.MatchType, check.BasicQuery.Title.MatchType)
	}
}

func TestUpdate(t *testing.T) {
	db := setupDB(t)
	repo := New(db)

	q := &query.Query{
		ID:                uuid.New(),
		UserID:            uuid.New(),
		QueryType:         query.BasicQueryType,
		BasicQueryVersion: "1.0",
		BasicQuery: &query.BasicQuery{
			Title: &query.StringMatch{
				Value:     "test",
				MatchType: query.ExactStringMatch,
			},
		},
		Status: query.PendingStatus,
	}
	defer db.Exec("DELETE FROM queries WHERE id = ?", q.ID)

	err := repo.Create(q)

	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}

	q.Status = query.CompletedStatus

	err = repo.Update(q)

	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}

	var check query.Query

	var basicQuery string

	err = db.QueryRow("SELECT id, user_id, status, query_type, basic_query_version, basic_query FROM queries WHERE id = ?", q.ID).Scan(&check.ID, &check.UserID, &check.Status, &check.QueryType, &check.BasicQueryVersion, &basicQuery)

	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}

	check.BasicQuery = &query.BasicQuery{}

	err = json.Unmarshal([]byte(basicQuery), check.BasicQuery)

	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}

	if check.ID != q.ID {
		t.Errorf("Expected %s, got %s", q.ID, check.ID)
	}

	if check.UserID != q.UserID {
		t.Errorf("Expected %s, got %s", q.UserID, check.UserID)
	}

	if check.Status != q.Status {
		t.Errorf("Expected %s, got %s", q.Status, check.Status)
	}

	if check.QueryType != q.QueryType {
		t.Errorf("Expected %s, got %s", q.QueryType, check.QueryType)
	}

	if check.BasicQueryVersion != q.BasicQueryVersion {
		t.Errorf("Expected %s, got %s", q.BasicQueryVersion, check.BasicQueryVersion)
	}

	if check.BasicQuery.Title.Value != q.BasicQuery.Title.Value {
		t.Errorf("Expected %s, got %s", q.BasicQuery.Title.Value, check.BasicQuery.Title.Value)
	}

	if check.BasicQuery.Title.MatchType != q.BasicQuery.Title.MatchType {
		t.Errorf("Expected %s, got %s", q.BasicQuery.Title.MatchType, check.BasicQuery.Title.MatchType)
	}
}
