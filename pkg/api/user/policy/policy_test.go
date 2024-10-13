package policy

import (
	"juno/pkg/api/user"
	"testing"

	"github.com/google/uuid"
)

func TestCreate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		result := CanCreate()

		if result.Allowed == false {
			t.Errorf("expected allowed to be true, got false")
		}
	})
}

func TestUpdate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		u := &user.User{ID: uuid.New()}
		result := CanUpdate(u, u)

		if result.Allowed == false {
			t.Errorf("expected allowed to be true, got false")
		}
	})

	t.Run("failure", func(t *testing.T) {
		u := &user.User{ID: uuid.New()}
		result := CanUpdate(u, &user.User{ID: uuid.New()})

		if result.Allowed == true {
			t.Errorf("expected allowed to be false, got true")
		}
	})
}

func TestRead(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		u := &user.User{ID: uuid.New()}
		result := CanRead(u, u)

		if result.Allowed == false {
			t.Errorf("expected allowed to be true, got false")
		}
	})

	t.Run("failure", func(t *testing.T) {
		u := &user.User{ID: uuid.New()}
		result := CanRead(u, &user.User{ID: uuid.New()})

		if result.Allowed == true {
			t.Errorf("expected allowed to be false, got true")
		}
	})
}

func TestDelete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		u := &user.User{ID: uuid.New()}
		result := CanDelete(u, u)

		if result.Allowed == false {
			t.Errorf("expected allowed to be true, got false")
		}
	})

	t.Run("failure", func(t *testing.T) {
		u := &user.User{ID: uuid.New()}
		result := CanDelete(u, &user.User{ID: uuid.New()})

		if result.Allowed == true {
			t.Errorf("expected allowed to be false, got true")
		}
	})
}
