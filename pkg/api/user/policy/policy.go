package policy

import (
	"context"
	"juno/pkg/api/auth"
	"juno/pkg/api/user"
	"juno/pkg/can"
)

func CanCreate() can.Result {
	return can.Allowed()
}

func CanUpdate(authUser *user.User, u *user.User) can.Result {
	if authUser.ID == u.ID {
		return can.Allowed()
	}

	return can.Denied("user not allowed to update user")
}

func CanRead(c context.Context, u *user.User) can.Result {
	authUser, ok := auth.UserFromContext(c)

	if !ok {
		return can.Error(user.ErrUserNotFoundInContext)
	}

	if authUser.ID == u.ID {
		return can.Allowed()
	}

	return can.Denied("user not allowed to read user")
}

func CanDelete(authUser *user.User, u *user.User) can.Result {
	if authUser.ID == u.ID {
		return can.Allowed()
	}

	return can.Denied("user not allowed to delete user")
}
