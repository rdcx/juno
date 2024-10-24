package mysql

import (
	"database/sql"
	"juno/pkg/api/extractor/job"

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

func (r *Repository) Get(id uuid.UUID) (*job.Job, error) {
	var j job.Job

	err := r.db.QueryRow("SELECT id, user_id, strategy_id, status, created_at, updated_at FROM jobs WHERE id = ?", id).Scan(&j.ID, &j.UserID, &j.StrategyID, &j.Status, &j.CreatedAt, &j.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, job.ErrNotFound
		}
		return nil, err
	}

	return &j, nil
}

func (r *Repository) Create(j *job.Job) error {
	_, err := r.db.Exec("INSERT INTO jobs (id, user_id, strategy_id, status) VALUES (?, ?, ?, ?)", j.ID, j.UserID, j.StrategyID, j.Status)

	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) ListByUserID(userID uuid.UUID) ([]*job.Job, error) {
	rows, err := r.db.Query("SELECT id, user_id, strategy_id, status, created_at, updated_at FROM jobs WHERE user_id = ?", userID)

	if err != nil {
		return nil, err
	}

	var jobs []*job.Job

	for rows.Next() {
		var j job.Job

		err := rows.Scan(&j.ID, &j.UserID, &j.StrategyID, &j.Status, &j.CreatedAt, &j.UpdatedAt)

		if err != nil {
			return nil, err
		}

		jobs = append(jobs, &j)
	}

	return jobs, nil
}

func (r *Repository) Update(j *job.Job) error {
	_, err := r.db.Exec("UPDATE jobs SET status = ? WHERE id = ?", j.Status, j.ID)

	if err != nil {
		return err
	}

	return nil
}
