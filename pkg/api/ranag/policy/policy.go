package policy

import (
	"context"
	"juno/pkg/api/auth"
	"juno/pkg/api/ranag"
	"juno/pkg/api/user"
	"juno/pkg/can"
)

type Policy struct{}

func New() Policy {
	return Policy{}
}

func (p Policy) CanUpdate(ctx context.Context, n *ranag.Ranag) can.Result {
	authUser, ok := auth.UserFromContext(ctx)

	if !ok {
		return can.Error(user.ErrUserNotFoundInContext)
	}

	if authUser.ID == n.OwnerID {
		return can.Allowed()
	}

	return can.Denied("user not allowed to update ranag")
}

func (p Policy) CanDelete(ctx context.Context, n *ranag.Ranag) can.Result {
	authUser, ok := auth.UserFromContext(ctx)

	if !ok {
		return can.Error(user.ErrUserNotFoundInContext)
	}

	if authUser.ID == n.OwnerID {
		return can.Allowed()
	}

	return can.Denied("user not allowed to delete ranag")
}

func (p Policy) CanRead(ctx context.Context, n *ranag.Ranag) can.Result {
	authUser, ok := auth.UserFromContext(ctx)

	if !ok {
		return can.Error(user.ErrUserNotFoundInContext)
	}

	if authUser.ID == n.OwnerID {
		return can.Allowed()
	}

	return can.Denied("user not allowed to read ranag")
}

func (p Policy) CanList(ctx context.Context, ranags []*ranag.Ranag) can.Result {
	authUser, ok := auth.UserFromContext(ctx)

	if !ok {
		return can.Error(user.ErrUserNotFoundInContext)
	}

	for _, n := range ranags {
		if authUser.ID != n.OwnerID {
			return can.Denied("user not allowed to list ranags")
		}
	}

	return can.Allowed()
}

func (p Policy) CanCreate() can.Result {
	return can.Allowed()
}
