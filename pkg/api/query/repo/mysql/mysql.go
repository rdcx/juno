package mysql

import (
	"database/sql"
	"encoding/json"
	"juno/pkg/api/query"

	"github.com/google/uuid"
)

type Repository struct {
	db *sql.DB
}

func New(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Get(id uuid.UUID) (*query.Query, error) {
	var q query.Query

	var basicQuery string

	err := r.db.QueryRow("SELECT id, user_id, status, query_type, basic_query_version, basic_query, created_at, updated_at FROM queries WHERE id = ?", id).Scan(&q.ID, &q.UserID, &q.Status, &q.QueryType, &q.BasicQueryVersion, &basicQuery, &q.CreatedAt, &q.UpdatedAt)
	if err != nil {
		return nil, err
	}

	q.BasicQuery = &query.BasicQuery{}

	err = json.Unmarshal([]byte(basicQuery), q.BasicQuery)
	if err != nil {
		return nil, err
	}

	return &q, nil
}

func (r *Repository) ListByUserID(userID uuid.UUID) ([]*query.Query, error) {
	rows, err := r.db.Query("SELECT id, user_id, status, query_type, basic_query_version, basic_query, created_at, updated_at FROM queries WHERE user_id = ?", userID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	queries := make([]*query.Query, 0)

	for rows.Next() {
		var q query.Query
		var basicQuery string

		err := rows.Scan(&q.ID, &q.UserID, &q.Status, &q.QueryType, &q.BasicQueryVersion, &basicQuery, &q.CreatedAt, &q.UpdatedAt)
		if err != nil {
			return nil, err
		}

		q.BasicQuery = &query.BasicQuery{}

		err = json.Unmarshal([]byte(basicQuery), q.BasicQuery)
		if err != nil {
			return nil, err
		}

		queries = append(queries, &q)
	}

	return queries, nil
}

func (r *Repository) Create(q *query.Query) error {

	basicQuery, err := json.Marshal(q.BasicQuery)

	if err != nil {
		return err
	}

	_, err = r.db.Exec("INSERT INTO queries (id, user_id, status, query_type, basic_query_version, basic_query) VALUES (?, ?, ?, ?, ?, ?)", q.ID, q.UserID, q.Status, q.QueryType, q.BasicQueryVersion, basicQuery)
	return err
}

func (r *Repository) Update(q *query.Query) error {
	basicQuery, err := json.Marshal(q.BasicQuery)

	if err != nil {
		return err
	}

	_, err = r.db.Exec("UPDATE queries SET status = ?, basic_query = ? WHERE id = ?", q.Status, basicQuery, q.ID)
	return err
}
