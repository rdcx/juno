package mem

import "juno/pkg/balancer/robotstxt"

type Repository struct {
	robotsTxt map[string]string
}

func New() *Repository {
	return &Repository{
		robotsTxt: make(map[string]string),
	}
}

func (r *Repository) Get(hostname string) (string, error) {

	if val, ok := r.robotsTxt[hostname]; ok {
		return val, nil
	}

	return "", robotstxt.ErrRobotsTxtNotInCache
}

func (r *Repository) Set(hostname, robotsTxt string) error {
	r.robotsTxt[hostname] = robotsTxt
	return nil
}
