package policy

import (
	"context"
	"juno/pkg/api/assignment"
	"juno/pkg/api/auth"
	"juno/pkg/can"
)

type Policy struct{}

func New() *Policy {
	return &Policy{}
}

func (p Policy) CanCreate() can.Result {
	return can.Allowed()
}

func (p Policy) CanRead(ctx context.Context, a *assignment.Assignment) can.Result {
	authUser, ok := auth.UserFromContext(ctx)

	if !ok {
		return can.Denied("no user in context")
	}

	if a.OwnerID == authUser.ID {
		return can.Allowed()
	}

	return can.Denied("user is not the owner")
}

func (p Policy) CanUpdate(ctx context.Context, a *assignment.Assignment) can.Result {
	authUser, ok := auth.UserFromContext(ctx)

	if !ok {
		return can.Denied("no user in context")
	}

	if a.OwnerID == authUser.ID {
		return can.Allowed()
	}

	return can.Denied("user is not the owner")
}

func (p Policy) CanDelete(ctx context.Context, a *assignment.Assignment) can.Result {
	authUser, ok := auth.UserFromContext(ctx)

	if !ok {
		return can.Denied("no user in context")
	}

	if a.OwnerID == authUser.ID {
		return can.Allowed()
	}

	return can.Denied("user is not the owner")
}

func (p Policy) CanList(ctx context.Context, assignments []*assignment.Assignment) can.Result {
	authUser, ok := auth.UserFromContext(ctx)

	if !ok {
		return can.Denied("no user in context")
	}

	for _, a := range assignments {
		if a.OwnerID != authUser.ID {
			return can.Denied("user is not the owner")
		}
	}

	return can.Allowed()
}
