package mysql

import (
	"database/sql"
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
	_, err := r.db.Exec("INSERT INTO balancers (id, owner_id, address) VALUES (?, ?, ?)", n.ID, n.OwnerID, n.Address)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repo) Get(id uuid.UUID) (*balancer.Balancer, error) {
	var n balancer.Balancer

	err := r.db.QueryRow("SELECT id, owner_id, address FROM balancers WHERE id = ?", id).Scan(&n.ID, &n.OwnerID, &n.Address)
	if err != nil {
		return nil, err
	}

	return &n, nil
}

func (r *Repo) ListByOwnerID(ownerID uuid.UUID) ([]*balancer.Balancer, error) {
	var balancers []*balancer.Balancer

	rows, err := r.db.Query("SELECT id, owner_id, address FROM balancers WHERE owner_id = ?", ownerID)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var n balancer.Balancer
		err = rows.Scan(&n.ID, &n.OwnerID, &n.Address)
		if err != nil {
			return nil, err
		}

		balancers = append(balancers, &n)
	}

	return balancers, nil
}

func (r *Repo) FirstWhereAddress(address string) (*balancer.Balancer, error) {
	var n balancer.Balancer

	err := r.db.QueryRow("SELECT id, owner_id, address, shards FROM balancers WHERE address = ?", address).Scan(&n.ID, &n.OwnerID, &n.Address)
	if err != nil {
		return nil, err
	}

	return &n, nil
}

func (r *Repo) Update(n *balancer.Balancer) error {

	_, err := r.db.Exec("UPDATE balancers SET owner_id = ?, address = ? WHERE id = ?", n.OwnerID, n.Address, n.ID)
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
