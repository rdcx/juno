package mysql

import (
	"database/sql"
	"juno/pkg/api/extractor/filter"

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

func (r *Repository) Get(id uuid.UUID) (*filter.Filter, error) {
	var f filter.Filter

	err := r.db.QueryRow("SELECT id, user_id, field_id, name, type, value, created_at, updated_at FROM filters WHERE id = ?", id).Scan(&f.ID, &f.UserID, &f.FieldID, &f.Name, &f.Type, &f.Value, &f.CreatedAt, &f.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, filter.ErrNotFound
		}
		return nil, err
	}

	return &f, nil
}

func (r *Repository) Create(f *filter.Filter) error {
	_, err := r.db.Exec("INSERT INTO filters (id, user_id, field_id, name, type, value) VALUES (?, ?, ?, ?, ?, ?)", f.ID, f.UserID, f.FieldID, f.Name, f.Type, f.Value)

	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) Update(f *filter.Filter) error {
	_, err := r.db.Exec("UPDATE filters SET name = ?, type = ?, value = ? WHERE id = ?", f.Name, f.Type, f.Value, f.ID)

	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) Delete(id uuid.UUID) error {
	_, err := r.db.Exec("DELETE FROM filters WHERE id = ?", id)

	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) ListByUserID(userID uuid.UUID) ([]*filter.Filter, error) {
	rows, err := r.db.Query("SELECT id, user_id, field_id, name, type, value, created_at, updated_at FROM filters WHERE user_id = ?", userID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var filters []*filter.Filter

	for rows.Next() {
		var f filter.Filter
		err := rows.Scan(&f.ID, &f.UserID, &f.FieldID, &f.Name, &f.Type, &f.Value, &f.CreatedAt, &f.UpdatedAt)

		if err != nil {
			return nil, err
		}

		filters = append(filters, &f)
	}

	return filters, nil
}
