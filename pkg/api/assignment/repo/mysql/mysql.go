package mysql

import (
	"database/sql"
	"juno/pkg/api/assignment"
	"strings"

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

func (r *Repository) Get(id uuid.UUID) (*assignment.Assignment, error) {
	var a assignment.Assignment

	err := r.db.QueryRow("SELECT id, owner_id, node_id, offset, length FROM assignments WHERE id = ?", id).Scan(&a.ID, &a.OwnerID, &a.NodeID, &a.Offset, &a.Length)

	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return nil, assignment.ErrNotFound
		}
		return nil, err
	}

	return &a, nil
}

func (r *Repository) ListByNodeID(nodeID uuid.UUID) ([]*assignment.Assignment, error) {
	var result []*assignment.Assignment

	rows, err := r.db.Query("SELECT id, owner_id, node_id, offset, length FROM assignments WHERE node_id = ?", nodeID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var a assignment.Assignment
		err := rows.Scan(&a.ID, &a.OwnerID, &a.NodeID, &a.Offset, &a.Length)

		if err != nil {
			return nil, err
		}

		result = append(result, &a)
	}

	return result, nil
}

func (r *Repository) Create(a *assignment.Assignment) error {
	_, err := r.db.Exec("INSERT INTO assignments (id, owner_id, node_id, offset, length) VALUES (?, ?, ?, ?, ?)", a.ID, a.OwnerID, a.NodeID, a.Offset, a.Length)

	return err
}

func (r *Repository) Update(a *assignment.Assignment) error {
	_, err := r.db.Exec("UPDATE assignments SET owner_id = ?, node_id = ?, offset = ?, length = ? WHERE id = ?", a.OwnerID, a.NodeID, a.Offset, a.Length, a.ID)

	return err
}

func (r *Repository) Delete(id uuid.UUID) error {
	_, err := r.db.Exec("DELETE FROM assignments WHERE id = ?", id)

	return err
}
