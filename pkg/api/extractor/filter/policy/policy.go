package policy

import (
	"context"
	"juno/pkg/api/auth"
	"juno/pkg/api/extractor/filter"
	"juno/pkg/can"
)

type Policy struct{}

func New() *Policy {
	return &Policy{}
}

func (p Policy) CanCreate() can.Result {
	return can.Allowed()
}

func (p Policy) CanRead(ctx context.Context, fil *filter.Filter) can.Result {
	u := auth.MustUserFromContext(ctx)

	if fil.UserID == u.ID {
		return can.Allowed()
	}

	return can.Denied("filter does not belong to user")
}

func (p Policy) CanUpdate(ctx context.Context, fil *filter.Filter) can.Result {
	u := auth.MustUserFromContext(ctx)

	if fil.UserID == u.ID {
		return can.Allowed()
	}

	return can.Denied("filter does not belong to user")
}

func (p Policy) CanDelete(ctx context.Context, fil *filter.Filter) can.Result {
	u := auth.MustUserFromContext(ctx)

	if fil.UserID == u.ID {
		return can.Allowed()
	}

	return can.Denied("filter does not belong to user")
}

func (p Policy) CanList(ctx context.Context, fils []*filter.Filter) can.Result {
	u := auth.MustUserFromContext(ctx)

	for _, fil := range fils {
		if fil.UserID != u.ID {
			return can.Denied("filter does not belong to user")
		}
	}

	return can.Allowed()
}
