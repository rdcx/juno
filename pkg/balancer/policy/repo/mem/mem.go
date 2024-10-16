package mem

import "juno/pkg/balancer/policy"

type Repository struct {
	policies map[string]*policy.CrawlPolicy
}

func New() *Repository {
	return &Repository{
		policies: make(map[string]*policy.CrawlPolicy),
	}
}

func (r *Repository) Get(hostname string) (*policy.CrawlPolicy, error) {
	p, ok := r.policies[hostname]

	if !ok {
		return nil, policy.ErrPolicyNotFound
	}

	return p, nil
}

func (r *Repository) Set(hostname string, policy *policy.CrawlPolicy) error {
	r.policies[hostname] = policy

	return nil
}
