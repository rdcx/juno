package mem

import (
	"juno/pkg/api/extractor/job"

	"github.com/google/uuid"
)

type Repository struct {
	jobs map[uuid.UUID]job.Job
}

func New() *Repository {
	return &Repository{
		jobs: make(map[uuid.UUID]job.Job),
	}
}

func (r *Repository) Create(q *job.Job) error {
	r.jobs[q.ID] = *q
	return nil
}

func (r *Repository) Get(id uuid.UUID) (*job.Job, error) {
	q, ok := r.jobs[id]
	if !ok {
		return nil, job.ErrNotFound
	}

	return &q, nil
}

func (r *Repository) ListByStatus(status job.JobStatus) ([]*job.Job, error) {
	var jobs []*job.Job
	for _, q := range r.jobs {
		if q.Status == status {
			jobs = append(jobs, &q)
		}
	}

	return jobs, nil
}

func (r *Repository) ListByUserID(userID uuid.UUID) ([]*job.Job, error) {
	var jobs []*job.Job
	for _, q := range r.jobs {
		if q.UserID == userID {
			jobs = append(jobs, &q)
		}
	}

	return jobs, nil
}

func (r *Repository) Update(q *job.Job) error {
	r.jobs[q.ID] = *q
	return nil
}
