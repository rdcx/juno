package mysql

import (
	"database/sql"
	"encoding/json"
	"juno/pkg/api/balancer"

	"github.com/google/uuid"
)

type Repo struct {
	db *sql.DB
}

func New(db *sql.DB) *Repo {
	return &Repo{db: db}
}

func (r *Repo) Create(n *balancer.Balancer) error {

	assignments, err := json.Marshal(n.ShardAssignments)
	if err != nil {
		return err
	}

	_, err = r.db.Exec("INSERT INTO balancers (id, owner_id, address, shard_assignments) VALUES (?, ?, ?, ?)", n.ID, n.OwnerID, n.Address, string(assignments))
	if err != nil {
		return err
	}

	return nil
}

func (r *Repo) All() ([]*balancer.Balancer, error) {
	var balancers []*balancer.Balancer

	rows, err := r.db.Query("SELECT id, owner_id, address, shard_assignments FROM balancers")

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var n balancer.Balancer
		var shardAssignmentJson string
		err = rows.Scan(&n.ID, &n.OwnerID, &n.Address, &shardAssignmentJson)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal([]byte(shardAssignmentJson), &n.ShardAssignments)
		if err != nil {
			return nil, err
		}

		balancers = append(balancers, &n)
	}

	return balancers, nil
}

func (r *Repo) Get(id uuid.UUID) (*balancer.Balancer, error) {
	var n balancer.Balancer

	var shardAssignmentJson string

	err := r.db.QueryRow("SELECT id, owner_id, address, shard_assignments FROM balancers WHERE id = ?", id).Scan(&n.ID, &n.OwnerID, &n.Address, &shardAssignmentJson)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(shardAssignmentJson), &n.ShardAssignments)
	if err != nil {
		return nil, err
	}

	return &n, nil
}

func (r *Repo) ListByOwnerID(ownerID uuid.UUID) ([]*balancer.Balancer, error) {
	var balancers []*balancer.Balancer

	rows, err := r.db.Query("SELECT id, owner_id, address, shard_assignments FROM balancers WHERE owner_id = ?", ownerID)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var n balancer.Balancer
		var shardAssignmentJson string
		err = rows.Scan(&n.ID, &n.OwnerID, &n.Address, &shardAssignmentJson)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal([]byte(shardAssignmentJson), &n.ShardAssignments)
		if err != nil {
			return nil, err
		}

		balancers = append(balancers, &n)
	}

	return balancers, nil
}

func (r *Repo) FirstWhereAddress(address string) (*balancer.Balancer, error) {
	var n balancer.Balancer
	var shardAssignmentJson string
	err := r.db.QueryRow("SELECT id, owner_id, address, shard_assignments FROM balancers WHERE address = ?", address).Scan(&n.ID, &n.OwnerID, &n.Address, &shardAssignmentJson)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(shardAssignmentJson), &n.ShardAssignments)
	if err != nil {
		return nil, err
	}

	return &n, nil
}

func (r *Repo) Update(n *balancer.Balancer) error {

	assignments, err := json.Marshal(n.ShardAssignments)
	if err != nil {
		return err
	}

	_, err = r.db.Exec("UPDATE balancers SET owner_id = ?, address = ?, shard_assignments = ? WHERE id = ?", n.OwnerID, n.Address, string(assignments), n.ID)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repo) Delete(id uuid.UUID) error {
	_, err := r.db.Exec("DELETE FROM balancers WHERE id = ?", id)
	if err != nil {
		return err
	}

	return nil
}
