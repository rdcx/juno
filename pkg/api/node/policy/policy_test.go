package policy

import (
	"juno/pkg/api/node"
	"juno/pkg/api/user"
	"testing"

	"github.com/google/uuid"
)

func TestCanCreate(t *testing.T) {
	t.Run("can", func(t *testing.T) {
		u := &user.User{
			ID: uuid.New(),
		}

		if !CanCreate(u) {
			t.Errorf("Expected true")
		}
	})
}

func TestCanRead(t *testing.T) {
	t.Run("can", func(t *testing.T) {
		u := &user.User{
			ID: uuid.New(),
		}

		n := &node.Node{
			OwnerID: u.ID,
		}

		if !CanRead(u, n) {
			t.Errorf("Expected true")
		}
	})

	t.Run("cannot", func(t *testing.T) {
		u := &user.User{
			ID: uuid.New(),
		}

		n := &node.Node{
			OwnerID: uuid.New(),
		}

		if CanRead(u, n) {
			t.Errorf("Expected false")
		}
	})
}

func TestCanUpdate(t *testing.T) {
	t.Run("can", func(t *testing.T) {
		u := &user.User{
			ID: uuid.New(),
		}

		n := &node.Node{
			OwnerID: u.ID,
		}

		if !CanUpdate(u, n) {
			t.Errorf("Expected true")
		}
	})

	t.Run("cannot", func(t *testing.T) {
		u := &user.User{
			ID: uuid.New(),
		}

		n := &node.Node{
			OwnerID: uuid.New(),
		}

		if CanUpdate(u, n) {
			t.Errorf("Expected false")
		}
	})
}

func TestCanDelete(t *testing.T) {
	t.Run("can", func(t *testing.T) {
		u := &user.User{
			ID: uuid.New(),
		}

		n := &node.Node{
			OwnerID: u.ID,
		}

		if !CanDelete(u, n) {
			t.Errorf("Expected true")
		}
	})

	t.Run("cannot", func(t *testing.T) {
		u := &user.User{
			ID: uuid.New(),
		}

		n := &node.Node{
			OwnerID: uuid.New(),
		}

		if CanDelete(u, n) {
			t.Errorf("Expected false")
		}
	})
}
