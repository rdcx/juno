package policy

import (
	"context"
	"juno/pkg/api/auth"
	"juno/pkg/api/extractor/selector"
	"juno/pkg/can"
)

type Policy struct{}

func New() *Policy {
	return &Policy{}
}

func (p Policy) CanCreate() can.Result {
	return can.Allowed()
}

func (p Policy) CanRead(ctx context.Context, sel *selector.Selector) can.Result {
	u := auth.MustUserFromContext(ctx)

	if sel.UserID == u.ID {
		return can.Allowed()
	}

	return can.Denied("selector does not belong to user")
}

func (p Policy) CanUpdate(ctx context.Context, sel *selector.Selector) can.Result {
	u := auth.MustUserFromContext(ctx)

	if sel.UserID == u.ID {
		return can.Allowed()
	}

	return can.Denied("selector does not belong to user")
}

func (p Policy) CanDelete(ctx context.Context, sel *selector.Selector) can.Result {
	u := auth.MustUserFromContext(ctx)

	if sel.UserID == u.ID {
		return can.Allowed()
	}

	return can.Denied("selector does not belong to user")
}

func (p Policy) CanList(ctx context.Context, sels []*selector.Selector) can.Result {
	u := auth.MustUserFromContext(ctx)

	for _, sel := range sels {
		if sel.UserID != u.ID {
			return can.Denied("selector does not belong to user")
		}
	}

	return can.Allowed()
}
