package mysql

import (
	"database/sql"
	"juno/pkg/api/extraction/strategy"
)

type Repository struct {
	db *sql.DB
}

func New(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Create(strategy *strategy.Strategy) error {
	res, err := r.db.Exec("INSERT INTO strategys (id, user_id, name, selector, filters) VALUES (?, ?, ?, ?, ?)", strategy.ID, strategy.UserID, strategy.Name, strategy.Selector, strategy.Filters)

	if err != nil {
		return err
	}

	return nil
}
