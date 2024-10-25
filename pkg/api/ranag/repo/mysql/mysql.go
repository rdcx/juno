package mysql

import (
	"database/sql"
	"encoding/json"
	"juno/pkg/api/ranag"

	"github.com/google/uuid"
)

type Repo struct {
	db *sql.DB
}

func New(db *sql.DB) *Repo {
	return &Repo{db: db}
}

func (r *Repo) Create(n *ranag.Ranag) error {

	assignments, err := json.Marshal(n.ShardAssignments)
	if err != nil {
		return err
	}

	_, err = r.db.Exec("INSERT INTO ranags (id, owner_id, address, shard_assignments) VALUES (?, ?, ?, ?)", n.ID, n.OwnerID, n.Address, string(assignments))
	if err != nil {
		return err
	}

	return nil
}

func (r *Repo) All() ([]*ranag.Ranag, error) {
	var ranags []*ranag.Ranag

	rows, err := r.db.Query("SELECT id, owner_id, address, shard_assignments FROM ranags")

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var n ranag.Ranag
		var shardAssignmentJson string
		err = rows.Scan(&n.ID, &n.OwnerID, &n.Address, &shardAssignmentJson)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal([]byte(shardAssignmentJson), &n.ShardAssignments)
		if err != nil {
			return nil, err
		}

		ranags = append(ranags, &n)
	}

	return ranags, nil
}

func (r *Repo) Get(id uuid.UUID) (*ranag.Ranag, error) {
	var n ranag.Ranag

	var shardAssignmentJson string

	err := r.db.QueryRow("SELECT id, owner_id, address, shard_assignments FROM ranags WHERE id = ?", id).Scan(&n.ID, &n.OwnerID, &n.Address, &shardAssignmentJson)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(shardAssignmentJson), &n.ShardAssignments)
	if err != nil {
		return nil, err
	}

	return &n, nil
}

func (r *Repo) ListByOwnerID(ownerID uuid.UUID) ([]*ranag.Ranag, error) {
	var ranags []*ranag.Ranag

	rows, err := r.db.Query("SELECT id, owner_id, address, shard_assignments FROM ranags WHERE owner_id = ?", ownerID)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var n ranag.Ranag
		var shardAssignmentJson string
		err = rows.Scan(&n.ID, &n.OwnerID, &n.Address, &shardAssignmentJson)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal([]byte(shardAssignmentJson), &n.ShardAssignments)
		if err != nil {
			return nil, err
		}

		ranags = append(ranags, &n)
	}

	return ranags, nil
}

func (r *Repo) FirstWhereAddress(address string) (*ranag.Ranag, error) {
	var n ranag.Ranag
	var shardAssignmentJson string
	err := r.db.QueryRow("SELECT id, owner_id, address, shard_assignments FROM ranags WHERE address = ?", address).Scan(&n.ID, &n.OwnerID, &n.Address, &shardAssignmentJson)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(shardAssignmentJson), &n.ShardAssignments)
	if err != nil {
		return nil, err
	}

	return &n, nil
}

func (r *Repo) Update(n *ranag.Ranag) error {

	assignments, err := json.Marshal(n.ShardAssignments)
	if err != nil {
		return err
	}

	_, err = r.db.Exec("UPDATE ranags SET owner_id = ?, address = ?, shard_assignments = ? WHERE id = ?", n.OwnerID, n.Address, string(assignments), n.ID)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repo) Delete(id uuid.UUID) error {
	_, err := r.db.Exec("DELETE FROM ranags WHERE id = ?", id)
	if err != nil {
		return err
	}

	return nil
}
