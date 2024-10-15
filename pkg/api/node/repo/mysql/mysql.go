package mysql

import (
	"database/sql"
	"encoding/json"
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

	assignments, err := json.Marshal(n.ShardAssignments)
	if err != nil {
		return err
	}

	_, err = r.db.Exec("INSERT INTO nodes (id, owner_id, address, shard_assignments) VALUES (?, ?, ?, ?)", n.ID, n.OwnerID, n.Address, assignments)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repo) Get(id uuid.UUID) (*node.Node, error) {
	var n node.Node

	var shardAssignmentJson string

	err := r.db.QueryRow("SELECT id, owner_id, address, shard_assignments FROM nodes WHERE id = ?", id).Scan(&n.ID, &n.OwnerID, &n.Address, &shardAssignmentJson)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(shardAssignmentJson), &n.ShardAssignments)
	if err != nil {
		return nil, err
	}

	return &n, nil
}

func (r *Repo) ListByOwnerID(ownerID uuid.UUID) ([]*node.Node, error) {
	var nodes []*node.Node

	rows, err := r.db.Query("SELECT id, owner_id, address, shard_assignments FROM nodes WHERE owner_id = ?", ownerID)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var n node.Node
		var shardAssignmentJson string
		err = rows.Scan(&n.ID, &n.OwnerID, &n.Address, &shardAssignmentJson)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal([]byte(shardAssignmentJson), &n.ShardAssignments)
		if err != nil {
			return nil, err
		}

		nodes = append(nodes, &n)
	}

	return nodes, nil
}

func (r *Repo) FirstWhereAddress(address string) (*node.Node, error) {
	var n node.Node
	var shardAssignmentJson string
	err := r.db.QueryRow("SELECT id, owner_id, address, shard_assignments FROM nodes WHERE address = ?", address).Scan(&n.ID, &n.OwnerID, &n.Address, &shardAssignmentJson)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(shardAssignmentJson), &n.ShardAssignments)
	if err != nil {
		return nil, err
	}

	return &n, nil
}

func (r *Repo) Update(n *node.Node) error {

	assignments, err := json.Marshal(n.ShardAssignments)
	if err != nil {
		return err
	}

	_, err = r.db.Exec("UPDATE nodes SET owner_id = ?, address = ?, shard_assignments = ? WHERE id = ?", n.OwnerID, n.Address, assignments, n.ID)
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
