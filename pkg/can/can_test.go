package can

import (
	"errors"
	"testing"
)

func Test(t *testing.T) {
	t.Run("allow", func(t *testing.T) {
		cr := Result{Allowed: true}

		cr.Allow(func() {
			t.Log("ok")
		})

		cr.Err(func(err error) {
			t.Error("expected err to not be called")
		})

		cr.Deny(func(reason string) {
			t.Error("expected deny to not be called")
		})
	})

	t.Run("deny", func(t *testing.T) {
		cr := Result{Allowed: false, Reason: "reason"}

		cr.Allow(func() {
			t.Error("expected allow to not be called")
		})

		cr.Err(func(err error) {
			t.Error("expected err to not be called")
		})

		cr.Deny(func(reason string) {
			if reason != "reason" {
				t.Errorf("expected reason to be reason, got %s", reason)
			}
			t.Log("failed")
		})
	})

	t.Run("error", func(t *testing.T) {
		err := errors.New("something went wrong")
		cr := Result{Allowed: false, Error: err}

		cr.Allow(func() {
			t.Error("expected allow to not be called")
		})

		cr.Err(func(e error) {
			if e != err {
				t.Errorf("expected error to be ErrNotAllowed, got %v", e)
			}
			t.Log("error")
		})

		cr.Deny(func(reason string) {
			t.Error("expected deny to not be called")
		})
	})

	t.Run("allowed", func(t *testing.T) {
		cr := Allowed()

		cr.Allow(func() {
			t.Log("ok")
		})

		cr.Err(func(err error) {
			t.Error("expected err to not be called")
		})

		cr.Deny(func(reason string) {
			t.Error("expected deny to not be called")
		})
	})

	t.Run("denied", func(t *testing.T) {
		cr := Denied("reason")

		cr.Allow(func() {
			t.Error("expected allow to not be called")
		})

		cr.Err(func(err error) {
			t.Error("expected err to not be called")
		})

		cr.Deny(func(reason string) {
			if reason != "reason" {
				t.Errorf("expected reason to be reason, got %s", reason)
			}
			t.Log("failed")
		})
	})

	t.Run("error", func(t *testing.T) {
		err := errors.New("something went wrong")
		cr := Error(err)

		cr.Allow(func() {
			t.Error("expected allow to not be called")
		})

		cr.Err(func(e error) {
			if e != err {
				t.Errorf("expected error to be ErrNotAllowed, got %v", e)
			}
			t.Log("error")
		})

		cr.Deny(func(reason string) {
			t.Error("expected deny to not be called")
		})
	})
}
