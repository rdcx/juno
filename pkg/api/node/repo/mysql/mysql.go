package mysql

import (
	"database/sql"
	"fmt"
	"juno/pkg/api/node/domain"
	"strings"
)

type Repo struct {
	db *sql.DB
}

func New(db *sql.DB) *Repo {
	return &Repo{db: db}
}

func (r *Repo) Create(n *domain.Node) error {
	shards := "["
	for _, s := range n.Shards {
		shards += fmt.Sprintf("%d", s) + ","
	}
	shards = strings.TrimSuffix(shards, ",")
	shards = shards + "]"

	_, err := r.db.Exec("INSERT INTO nodes (id, owner_id, address, shards) VALUES (?, ?, ?, ?)", n.ID, n.OwnerID, n.Address, shards)
	if err != nil {
		return err
	}

	return nil
}
