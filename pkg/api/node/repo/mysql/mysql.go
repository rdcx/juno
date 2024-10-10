package mysql

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"juno/pkg/api/node"
	"strings"

	"github.com/google/uuid"
)

type Repo struct {
	db *sql.DB
}

func New(db *sql.DB) *Repo {
	return &Repo{db: db}
}

func (r *Repo) Create(n *node.Node) error {
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

func (r *Repo) Get(id uuid.UUID) (*node.Node, error) {
	var n node.Node

	var shards string

	err := r.db.QueryRow("SELECT id, owner_id, address, shards FROM nodes WHERE id = ?", id).Scan(&n.ID, &n.OwnerID, &n.Address, &shards)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(shards), &n.Shards)
	if err != nil {
		return nil, err
	}

	return &n, nil
}

func (r *Repo) FirstWhereAddress(address string) (*node.Node, error) {
	var n node.Node
	var shards string

	err := r.db.QueryRow("SELECT id, owner_id, address, shards FROM nodes WHERE address = ?", address).Scan(&n.ID, &n.OwnerID, &n.Address, &shards)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(shards), &n.Shards)
	if err != nil {
		return nil, err
	}

	return &n, nil
}

func (r *Repo) Update(n *node.Node) error {
	shards := "["
	for _, s := range n.Shards {
		shards += fmt.Sprintf("%d", s) + ","
	}
	shards = strings.TrimSuffix(shards, ",")
	shards = shards + "]"

	_, err := r.db.Exec("UPDATE nodes SET owner_id = ?, address = ?, shards = ? WHERE id = ?", n.OwnerID, n.Address, shards, n.ID)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repo) Delete(id uuid.UUID) error {
	_, err := r.db.Exec("DELETE FROM nodes WHERE id = ?", id)
	if err != nil {
		return err
	}

	return nil
}
