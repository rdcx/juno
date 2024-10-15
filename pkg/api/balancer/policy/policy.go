package policy

import (
	"context"
	"juno/pkg/api/auth"
	"juno/pkg/api/balancer"
	"juno/pkg/api/user"
	"juno/pkg/can"
)

type Policy struct{}

func New() Policy {
	return Policy{}
}

func (p Policy) CanUpdate(ctx context.Context, n *balancer.Balancer) can.Result {
	authUser, ok := auth.UserFromContext(ctx)

	if !ok {
		return can.Error(user.ErrUserNotFoundInContext)
	}

	if authUser.ID == n.OwnerID {
		return can.Allowed()
	}

	return can.Denied("user not allowed to update balancer")
}

func (p Policy) CanDelete(ctx context.Context, n *balancer.Balancer) can.Result {
	authUser, ok := auth.UserFromContext(ctx)

	if !ok {
		return can.Error(user.ErrUserNotFoundInContext)
	}

	if authUser.ID == n.OwnerID {
		return can.Allowed()
	}

	return can.Denied("user not allowed to delete balancer")
}

func (p Policy) CanRead(ctx context.Context, n *balancer.Balancer) can.Result {
	authUser, ok := auth.UserFromContext(ctx)

	if !ok {
		return can.Error(user.ErrUserNotFoundInContext)
	}

	if authUser.ID == n.OwnerID {
		return can.Allowed()
	}

	return can.Denied("user not allowed to read balancer")
}

func (p Policy) CanList(ctx context.Context, balancers []*balancer.Balancer) can.Result {
	authUser, ok := auth.UserFromContext(ctx)

	if !ok {
		return can.Error(user.ErrUserNotFoundInContext)
	}

	for _, n := range balancers {
		if authUser.ID != n.OwnerID {
			return can.Denied("user not allowed to list balancers")
		}
	}

	return can.Allowed()
}

func (p Policy) CanCreate() can.Result {
	return can.Allowed()
}
