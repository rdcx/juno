package service

import (
	"sync"
	"time"
)

type Policy struct {
	Interval  time.Duration
	LastCrawl time.Time
	Lock      sync.Mutex
}

type PolicyManager struct {
	Policies map[string]*Policy
}

func NewPolicyManager() *PolicyManager {
	return &PolicyManager{
		Policies: make(map[string]*Policy),
	}
}

func (m *PolicyManager) GetPolicy(hostname string) *Policy {
	m.Policies[hostname] = &Policy{}
	return m.Policies[hostname]
}

func (m *PolicyManager) SetPolicy(hostname string, policy *Policy) {
	m.Policies[hostname] = policy
}

func (p *PolicyManager) CanCrawl(hostname string) bool {
	policy := p.GetPolicy(hostname)

	policy.Lock.Lock()
	defer policy.Lock.Unlock()

	if time.Since(policy.LastCrawl) < policy.Interval {
		return false
	}

	policy.LastCrawl = time.Now()

	return true
}
