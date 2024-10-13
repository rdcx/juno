package mem

import (
	"juno/pkg/api/user"

	"github.com/google/uuid"
)

type Repo struct {
	users map[uuid.UUID]*user.User
}

func New() *Repo {
	return &Repo{
		users: make(map[uuid.UUID]*user.User),
	}
}

func (r *Repo) Create(u *user.User) error {

	r.users[u.ID] = u

	return nil
}

func (r *Repo) Get(id uuid.UUID) (*user.User, error) {
	u, ok := r.users[id]

	if !ok {
		return nil, user.ErrNotFound
	}

	return u, nil
}

func (r *Repo) FirstWhereEmail(email string) (*user.User, error) {
	for _, u := range r.users {
		if u.Email == email {
			return u, nil
		}
	}

	return nil, user.ErrNotFound
}

func (r *Repo) Update(u *user.User) error {
	r.users[u.ID] = u

	return nil
}

func (r *Repo) Delete(id uuid.UUID) error {
	delete(r.users, id)

	return nil
}
