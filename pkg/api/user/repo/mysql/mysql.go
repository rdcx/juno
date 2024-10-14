package mysql

import (
	"database/sql"
	"juno/pkg/api/user"

	"github.com/google/uuid"
)

type Repo struct {
	db *sql.DB
}

func New(db *sql.DB) *Repo {
	return &Repo{
		db: db,
	}
}

func (r *Repo) Create(u *user.User) error {
	_, err := r.db.Exec("INSERT INTO users (id, email, password) VALUES (?, ?, ?)", u.ID, u.Email, u.Password)

	return err
}

func (r *Repo) Get(id uuid.UUID) (*user.User, error) {
	var u user.User

	err := r.db.QueryRow("SELECT id, email, password FROM users WHERE id = ?", id).Scan(&u.ID, &u.Email, &u.Password)

	if err != nil {
		return nil, user.ErrNotFound
	}

	return &u, nil
}

func (r *Repo) FirstWhereEmail(email string) (*user.User, error) {
	var u user.User

	err := r.db.QueryRow("SELECT id, email, password FROM users WHERE email = ?", email).Scan(&u.ID, &u.Email, &u.Password)

	if err != nil {
		return nil, user.ErrNotFound
	}

	return &u, nil
}

func (r *Repo) Update(u *user.User) error {
	_, err := r.db.Exec("UPDATE users SET email = ?, password = ? WHERE id = ?", u.Email, u.Password, u.ID)

	return err
}

func (r *Repo) Delete(id uuid.UUID) error {
	_, err := r.db.Exec("DELETE FROM users WHERE id = ?", id)

	return err
}
