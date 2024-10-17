package queue

import "errors"

var ErrNoURLsInQueue = errors.New("no urls in queue")
var ErrProcessQueueCancelled = errors.New("process queue cancelled")

type Service interface {
	Push(url string) error
	Pop() (string, error)
}

type Repository interface {
	Exists(url string) (bool, error)
	Push(url string) error
	Pop() (string, error)
}
