package policy

import (
	"context"
	"juno/pkg/api/auth"
	"juno/pkg/api/extractor/job"
	"juno/pkg/can"
)

type Policy struct{}

func New() *Policy {
	return &Policy{}
}

func (p *Policy) CanCreate() can.Result {
	return can.Allowed()
}

func (p *Policy) CanUpdate(ctx context.Context, j *job.Job) can.Result {
	user := auth.MustUserFromContext(ctx)

	if j.UserID != user.ID {
		return can.Denied("job does not belong to user")
	}

	return can.Allowed()
}

func (p *Policy) CanGet(ctx context.Context, j *job.Job) can.Result {
	user := auth.MustUserFromContext(ctx)

	if j.UserID != user.ID {
		return can.Denied("job does not belong to user")
	}

	return can.Allowed()
}

func (p *Policy) CanList(ctx context.Context, jobs []*job.Job) can.Result {
	user := auth.MustUserFromContext(ctx)

	for _, j := range jobs {
		if j.UserID != user.ID {
			return can.Denied("job does not belong to user")
		}
	}

	return can.Allowed()
}
