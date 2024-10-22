package dto

import (
	"juno/pkg/api/extraction/job"
	"time"
)

const (
	SUCCESS = "success"
	ERROR   = "error"
)

type Job struct {
	ID          string `json:"id"`
	UserID      string `json:"user_id"`
	ExtractorID string `json:"extractor_id"`
	Status      string `json:"status"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type CreateJobRequest struct {
	ExtractorID string `json:"extractor_id" validate:"required,uuid"`
}

type CreateJobResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`

	Job *Job `json:"result,omitempty"`
}

func NewJobFromDomain(j *job.Job) *Job {
	return &Job{
		ID:          j.ID.String(),
		UserID:      j.UserID.String(),
		ExtractorID: j.ExtractorID.String(),
		Status:      string(j.Status),
		CreatedAt:   j.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   j.UpdatedAt.Format(time.RFC3339),
	}
}

func NewSuccessCreateJobResponse(job *job.Job) CreateJobResponse {
	return CreateJobResponse{
		Status: SUCCESS,
		Job:    NewJobFromDomain(job),
	}
}

func NewErrorCreateJobResponse(message string) CreateJobResponse {
	return CreateJobResponse{
		Status:  ERROR,
		Message: message,
	}
}

type GetJobResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`

	Job *Job `json:"result,omitempty"`
}

func NewSuccessGetJobResponse(job *job.Job) GetJobResponse {
	return GetJobResponse{
		Status: SUCCESS,
		Job:    NewJobFromDomain(job),
	}
}

func NewErrorGetJobResponse(message string) GetJobResponse {
	return GetJobResponse{
		Status:  ERROR,
		Message: message,
	}
}

type ListJobsResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`

	Jobs []*Job `json:"result,omitempty"`
}

func NewSuccessListJobsResponse(jobs []*job.Job) ListJobsResponse {

	js := make([]*Job, 0, len(jobs))

	for _, j := range jobs {
		js = append(js, NewJobFromDomain(j))
	}
	return ListJobsResponse{
		Status: SUCCESS,
		Jobs:   js,
	}
}

func NewErrorListJobsResponse(message string) ListJobsResponse {
	return ListJobsResponse{
		Status:  ERROR,
		Message: message,
	}
}
