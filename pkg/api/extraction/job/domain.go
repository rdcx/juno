package job

import (
	"context"
	"errors"
	"juno/pkg/api/extraction/extractor"
	"juno/pkg/can"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var (
	ErrNotFound = errors.New("job not found")
)

type JobStatus string

const (
	PendingStatus   JobStatus = "pending"
	RunningStatus   JobStatus = "running"
	CompletedStatus JobStatus = "completed"
	FailedStatus    JobStatus = "failed"
)

type Job struct {
	ID          uuid.UUID `json:"id"`
	UserID      uuid.UUID `json:"user_id"`
	Status      JobStatus `json:"status"`
	ExtractorID uuid.UUID `json:"extractor_id"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewJob(userID uuid.UUID, extractor *extractor.Extractor) *Job {
	return &Job{
		ID:          uuid.New(),
		UserID:      userID,
		Status:      PendingStatus,
		ExtractorID: extractor.ID,

		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

type Repository interface {
	Create(job *Job) error
	Get(id uuid.UUID) (*Job, error)
	ListByUserID(userID uuid.UUID) ([]*Job, error)
	Update(job *Job) error
}

type Service interface {
	Create(userID uuid.UUID, extractorID uuid.UUID) (*Job, error)
	Get(id uuid.UUID) (*Job, error)
	ListByUserID(userID uuid.UUID) ([]*Job, error)
}

type Handler interface {
	Create(c *gin.Context)
	Get(c *gin.Context)
	List(c *gin.Context)
}

type Policy interface {
	CanCreate() can.Result
	CanGet(ctx context.Context, j *Job) can.Result
	CanList(ctx context.Context, jobs []*Job) can.Result
	CanUpdate(ctx context.Context, j *Job) can.Result
}
