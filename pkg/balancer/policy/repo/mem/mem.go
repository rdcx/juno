package mem

import "juno/pkg/balancer/policy"

type Repository struct {
	policies map[string]*policy.HostnamePolicy
}

func New() *Repository {
	return &Repository{
		policies: make(map[string]*policy.HostnamePolicy),
	}
}

func (r *Repository) Get(hostname string) (*policy.HostnamePolicy, error) {
	p, ok := r.policies[hostname]

	if !ok {
		return nil, policy.ErrPolicyNotFound
	}

	return p, nil
}

func (r *Repository) Set(hostname string, policy *policy.HostnamePolicy) error {
	r.policies[hostname] = policy

	return nil
}
