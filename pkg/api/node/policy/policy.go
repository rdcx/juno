package policy

import (
	"context"
	"juno/pkg/api/auth"
	"juno/pkg/api/node"
	"juno/pkg/api/user"
	"juno/pkg/can"
)

type Policy struct{}

func New() Policy {
	return Policy{}
}

func (p Policy) CanUpdate(ctx context.Context, n *node.Node) can.Result {
	authUser, ok := auth.UserFromContext(ctx)

	if !ok {
		return can.Error(user.ErrUserNotFoundInContext)
	}

	if authUser.ID == n.OwnerID {
		return can.Allowed()
	}

	return can.Denied("user not allowed to update node")
}

func (p Policy) CanDelete(ctx context.Context, n *node.Node) can.Result {
	authUser, ok := auth.UserFromContext(ctx)

	if !ok {
		return can.Error(user.ErrUserNotFoundInContext)
	}

	if authUser.ID == n.OwnerID {
		return can.Allowed()
	}

	return can.Denied("user not allowed to delete node")
}

func (p Policy) CanRead(ctx context.Context, n *node.Node) can.Result {
	authUser, ok := auth.UserFromContext(ctx)

	if !ok {
		return can.Error(user.ErrUserNotFoundInContext)
	}

	if authUser.ID == n.OwnerID {
		return can.Allowed()
	}

	return can.Denied("user not allowed to read node")
}

func (p Policy) CanCreate() can.Result {
	return can.Allowed()
}
