package service

const CRAWL_TIMEOUT = 10

type Crawler struct {
	Queue         *Queue
	PolicyManager *PolicyManager
}

func NewCrawler() *Crawler {
	return &Crawler{
		Queue:         NewQueue(),
		PolicyManager: NewPolicyManager(),
	}
}
