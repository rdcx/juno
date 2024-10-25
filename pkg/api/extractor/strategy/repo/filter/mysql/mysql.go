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

func (r *Repository) AddFilter(strategyID, filterID uuid.UUID) error {
	_, err := r.db.Exec("INSERT INTO strategy_filters (strategy_id, filter_id) VALUES (?, ?)", strategyID, filterID)

	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) ListFilterIDs(strategyID uuid.UUID) ([]uuid.UUID, error) {
	rows, err := r.db.Query("SELECT filter_id FROM strategy_filters WHERE strategy_id = ?", strategyID)

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

func (r *Repository) RemoveFilter(strategyID, filterID uuid.UUID) error {
	_, err := r.db.Exec("DELETE FROM strategy_filters WHERE strategy_id = ? AND filter_id = ?", strategyID, filterID)

	if err != nil {
		return err
	}

	return nil
}
