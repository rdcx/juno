package mysql

import (
	"database/sql"
	"juno/pkg/api/node"

	"github.com/google/uuid"
)

type Repo struct {
	db *sql.DB
}

func New(db *sql.DB) *Repo {
	return &Repo{db: db}
}

func (r *Repo) Create(n *node.Node) error {
	_, err := r.db.Exec("INSERT INTO nodes (id, owner_id, address) VALUES (?, ?, ?)", n.ID, n.OwnerID, n.Address)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repo) Get(id uuid.UUID) (*node.Node, error) {
	var n node.Node

	err := r.db.QueryRow("SELECT id, owner_id, address FROM nodes WHERE id = ?", id).Scan(&n.ID, &n.OwnerID, &n.Address)
	if err != nil {
		return nil, err
	}

	return &n, nil
}

func (r *Repo) ListByOwnerID(ownerID uuid.UUID) ([]*node.Node, error) {
	var nodes []*node.Node

	rows, err := r.db.Query("SELECT id, owner_id, address FROM nodes WHERE owner_id = ?", ownerID)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var n node.Node
		err = rows.Scan(&n.ID, &n.OwnerID, &n.Address)
		if err != nil {
			return nil, err
		}

		nodes = append(nodes, &n)
	}

	return nodes, nil
}

func (r *Repo) FirstWhereAddress(address string) (*node.Node, error) {
	var n node.Node

	err := r.db.QueryRow("SELECT id, owner_id, address, shards FROM nodes WHERE address = ?", address).Scan(&n.ID, &n.OwnerID, &n.Address)
	if err != nil {
		return nil, err
	}

	return &n, nil
}

func (r *Repo) Update(n *node.Node) error {

	_, err := r.db.Exec("UPDATE nodes SET owner_id = ?, address = ? WHERE id = ?", n.OwnerID, n.Address, n.ID)
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
