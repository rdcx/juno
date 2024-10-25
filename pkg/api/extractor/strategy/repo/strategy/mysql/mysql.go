package mysql

import (
	"database/sql"
	"juno/pkg/api/extractor/strategy"

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

func (r *Repository) Get(id uuid.UUID) (*strategy.Strategy, error) {
	var s strategy.Strategy

	err := r.db.QueryRow("SELECT id, user_id, name, created_at, updated_at FROM strategies WHERE id = ?", id).Scan(&s.ID, &s.UserID, &s.Name, &s.CreatedAt, &s.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, strategy.ErrNotFound
		}
		return nil, err
	}

	return &s, nil
}

func (r *Repository) Create(s *strategy.Strategy) error {
	_, err := r.db.Exec("INSERT INTO strategies (id, user_id, name) VALUES (?, ?, ?)", s.ID, s.UserID, s.Name)

	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) ListByUserID(userID uuid.UUID) ([]*strategy.Strategy, error) {
	rows, err := r.db.Query("SELECT id, user_id, name, created_at, updated_at FROM strategies WHERE user_id = ?", userID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var strats []*strategy.Strategy

	for rows.Next() {
		var s strategy.Strategy

		err := rows.Scan(&s.ID, &s.UserID, &s.Name, &s.CreatedAt, &s.UpdatedAt)

		if err != nil {
			return nil, err
		}

		strats = append(strats, &s)
	}

	return strats, nil
}

func (r *Repository) Update(s *strategy.Strategy) error {
	_, err := r.db.Exec("UPDATE strategies SET name = ? WHERE id = ?", s.Name, s.ID)

	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) Delete(id uuid.UUID) error {
	_, err := r.db.Exec("DELETE FROM strategies WHERE id = ?", id)

	if err != nil {
		return err
	}

	return nil
}
