package policy

import (
	"errors"
	"time"
)

const DefaultCrawlInterval = 30 * time.Second

var ErrPolicyNotFound = errors.New("policy not found")

type CrawlPolicy struct {
	Hostname      string
	CrawlInterval time.Duration
	LastCrawled   time.Time

	// The number of times the hostname has been crawled
	TimesCrawled int
}

func New(hostname string) *CrawlPolicy {
	return &CrawlPolicy{
		Hostname:      hostname,
		CrawlInterval: DefaultCrawlInterval,
		LastCrawled:   time.Time{},
		TimesCrawled:  0,
	}
}

type Repository interface {
	Get(hostname string) (*CrawlPolicy, error)
	Set(hostname string, policy *CrawlPolicy) error
}

type Service interface {
	Get(hostname string) (*CrawlPolicy, error)
	Set(hostname string, policy *CrawlPolicy) error
	CanCrawl(p *CrawlPolicy) bool
	RecordCrawl(p *CrawlPolicy) error
}
