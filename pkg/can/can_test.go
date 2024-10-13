package can

import (
	"errors"
	"testing"
)

func Test(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		cr := Result{Allowed: true}

		cr.Ok(func() {
			t.Log("ok")
		})

		cr.Err(func(err error) {
			t.Error("expected err to not be called")
		})

		cr.Failed(func(reason string) {
			t.Error("expected failed to not be called")
		})
	})

	t.Run("failure", func(t *testing.T) {
		cr := Result{Allowed: false, Reason: "reason"}

		cr.Ok(func() {
			t.Error("expected ok to not be called")
		})

		cr.Err(func(err error) {
			t.Error("expected err to not be called")
		})

		cr.Failed(func(reason string) {
			if reason != "reason" {
				t.Errorf("expected reason to be reason, got %s", reason)
			}
			t.Log("failed")
		})
	})

	t.Run("error", func(t *testing.T) {
		err := errors.New("something went wrong")
		cr := Result{Allowed: false, Error: err}

		cr.Yes(func() {
			t.Error("expected yes to not be called")
		})

		cr.Err(func(e error) {
			if e != err {
				t.Errorf("expected error to be ErrNotAllowed, got %v", e)
			}
			t.Log("error")
		})

		cr.No(func(reason string) {
			t.Error("expected no to not be called")
		})
	})
}
