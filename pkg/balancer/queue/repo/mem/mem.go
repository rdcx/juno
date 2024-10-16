package mem

import "juno/pkg/balancer/queue"

type Repository struct {
	urls []string
}

func New() *Repository {
	return &Repository{}
}

func (r *Repository) Push(url string) error {
	r.urls = append(r.urls, url)
	return nil
}

func (r *Repository) Pop() (string, error) {
	if len(r.urls) == 0 {
		return "", queue.ErrNoURLsInQueue
	}

	url := r.urls[0]
	r.urls = r.urls[1:]
	return url, nil
}
