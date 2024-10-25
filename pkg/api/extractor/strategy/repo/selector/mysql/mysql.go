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

func (r *Repository) AddSelector(strategyID, selectorID uuid.UUID) error {
	_, err := r.db.Exec("INSERT INTO strategy_selectors (strategy_id, selector_id) VALUES (?, ?)", strategyID, selectorID)

	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) ListSelectorIDs(strategyID uuid.UUID) ([]uuid.UUID, error) {
	rows, err := r.db.Query("SELECT selector_id FROM strategy_selectors WHERE strategy_id = ?", strategyID)

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

func (r *Repository) RemoveSelector(strategyID, selectorID uuid.UUID) error {
	_, err := r.db.Exec("DELETE FROM strategy_selectors WHERE strategy_id = ? AND selector_id = ?", strategyID, selectorID)

	if err != nil {
		return err
	}

	return nil
}
