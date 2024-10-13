package policy

import (
	"context"
	"juno/pkg/api/auth"
	"juno/pkg/api/user"
	"juno/pkg/can"
)

type Policy struct{}

func New() *Policy {
	return &Policy{}
}

func (p Policy) CanCreate() can.Result {
	return can.Allowed()
}

func (p Policy) CanUpdate(c context.Context, u *user.User) can.Result {
	authUser, ok := auth.UserFromContext(c)

	if !ok {
		return can.Error(user.ErrUserNotFoundInContext)
	}

	if authUser.ID == u.ID {
		return can.Allowed()
	}

	return can.Denied("user not allowed to update user")
}

func (p Policy) CanRead(c context.Context, u *user.User) can.Result {
	authUser, ok := auth.UserFromContext(c)

	if !ok {
		return can.Error(user.ErrUserNotFoundInContext)
	}

	if authUser.ID == u.ID {
		return can.Allowed()
	}

	return can.Denied("user not allowed to read user")
}

func (p Policy) CanDelete(c context.Context, u *user.User) can.Result {
	authUser, ok := auth.UserFromContext(c)

	if !ok {
		return can.Error(user.ErrUserNotFoundInContext)
	}

	if authUser.ID == u.ID {
		return can.Allowed()
	}

	return can.Denied("user not allowed to delete user")
}
