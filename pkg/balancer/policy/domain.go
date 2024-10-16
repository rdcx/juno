package policy

import (
	"errors"
	"time"
)

var ErrPolicyNotFound = errors.New("policy not found")

type Policy struct {
	Hostname      string
	CrawlInterval time.Duration
	LastCrawled   time.Time

	// The number of times the hostname has been crawled
	Crawled int
}

type Repository interface {
	Get(hostname string) (*Policy, error)
	Set(hostname string, policy *Policy) error
}

type Service interface {
	CanCrawl(hostname string) (bool, error)
	Update(policy *Policy) error
}
