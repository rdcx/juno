package mysql

import (
	"database/sql"

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

func (r *Repository) AddField(strategyID, fieldID uuid.UUID) error {
	_, err := r.db.Exec("INSERT INTO strategy_fields (strategy_id, field_id) VALUES (?, ?)", strategyID, fieldID)

	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) ListFieldIDs(strategyID uuid.UUID) ([]uuid.UUID, error) {
	rows, err := r.db.Query("SELECT field_id FROM strategy_fields WHERE strategy_id = ?", strategyID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var ids []uuid.UUID

	for rows.Next() {
		var id uuid.UUID

		err := rows.Scan(&id)

		if err != nil {
			return nil, err
		}

		ids = append(ids, id)
	}

	return ids, nil
}

func (r *Repository) RemoveField(strategyID, fieldID uuid.UUID) error {
	_, err := r.db.Exec("DELETE FROM strategy_fields WHERE strategy_id = ? AND field_id = ?", strategyID, fieldID)

	if err != nil {
		return err
	}

	return nil
}
