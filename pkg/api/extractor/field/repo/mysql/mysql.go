package mysql

import (
	"database/sql"
	"juno/pkg/api/extractor/field"

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

func (r *Repository) Get(id uuid.UUID) (*field.Field, error) {
	var f field.Field

	err := r.db.QueryRow("SELECT id, user_id, selector_id, name, type, created_at, updated_at FROM fields WHERE id = ?", id).Scan(&f.ID, &f.UserID, &f.SelectorID, &f.Name, &f.Type, &f.CreatedAt, &f.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, field.ErrNotFound
		}
		return nil, err
	}

	return &f, nil
}

func (r *Repository) Create(f *field.Field) error {
	_, err := r.db.Exec("INSERT INTO fields (id, user_id, selector_id, name, type) VALUES (?, ?, ?, ?, ?)", f.ID, f.UserID, f.SelectorID, f.Name, f.Type)

	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) ListByUserID(userID uuid.UUID) ([]*field.Field, error) {
	rows, err := r.db.Query("SELECT id, user_id, selector_id, name, type, created_at, updated_at FROM fields WHERE user_id = ?", userID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var fields []*field.Field

	for rows.Next() {
		var f field.Field

		err := rows.Scan(&f.ID, &f.UserID, &f.SelectorID, &f.Name, &f.Type, &f.CreatedAt, &f.UpdatedAt)

		if err != nil {
			return nil, err
		}

		fields = append(fields, &f)
	}

	return fields, nil
}

func (r *Repository) ListBySelectorID(selectorID uuid.UUID) ([]*field.Field, error) {
	rows, err := r.db.Query("SELECT id, user_id, selector_id, name, type, created_at, updated_at FROM fields WHERE selector_id = ?", selectorID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var fields []*field.Field

	for rows.Next() {
		var f field.Field

		err := rows.Scan(&f.ID, &f.UserID, &f.SelectorID, &f.Name, &f.Type, &f.CreatedAt, &f.UpdatedAt)

		if err != nil {
			return nil, err
		}

		fields = append(fields, &f)
	}

	return fields, nil
}

func (r *Repository) Update(f *field.Field) error {
	_, err := r.db.Exec("UPDATE fields SET name = ?, type = ? WHERE id = ?", f.Name, f.Type, f.ID)

	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) Delete(id uuid.UUID) error {
	_, err := r.db.Exec("DELETE FROM fields WHERE id = ?", id)

	if err != nil {
		return err
	}

	return nil
}
