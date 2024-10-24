package mysql

import (
	"database/sql"
	"juno/pkg/api/extractor/selector"

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

func (r *Repository) Get(id uuid.UUID) (*selector.Selector, error) {
	var j selector.Selector

	err := r.db.QueryRow("SELECT id, user_id, name, value, visibility, created_at, updated_at FROM selectors WHERE id = ?", id).Scan(&j.ID, &j.UserID, &j.Name, &j.Value, &j.Visibility, &j.CreatedAt, &j.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, selector.ErrNotFound
		}
		return nil, err
	}

	return &j, nil
}

func (r *Repository) Create(j *selector.Selector) error {
	_, err := r.db.Exec("INSERT INTO selectors (id, user_id, name, value, visibility) VALUES (?, ?, ?, ?, ?)", j.ID, j.UserID, j.Name, j.Value, j.Visibility)

	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) ListByUserID(userID uuid.UUID) ([]*selector.Selector, error) {
	rows, err := r.db.Query("SELECT id, user_id, name, value, visibility, created_at, updated_at FROM selectors WHERE user_id = ?", userID)

	if err != nil {
		return nil, err
	}

	var selectors []*selector.Selector

	for rows.Next() {
		var j selector.Selector

		err := rows.Scan(&j.ID, &j.UserID, &j.Name, &j.Value, &j.Visibility, &j.CreatedAt, &j.UpdatedAt)

		if err != nil {
			return nil, err
		}

		selectors = append(selectors, &j)
	}

	return selectors, nil
}

func (r *Repository) Update(j *selector.Selector) error {
	_, err := r.db.Exec("UPDATE selectors SET name = ?, value = ?, visibility = ? WHERE id = ?", j.Name, j.Value, j.Visibility, j.ID)

	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) Delete(id uuid.UUID) error {
	_, err := r.db.Exec("DELETE FROM selectors WHERE id = ?", id)

	if err != nil {
		return err
	}

	return nil
}
