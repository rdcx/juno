package robotstxt

import "errors"

var ErrRobotsTxtNotInCache = errors.New("robots.txt not in cache")
var ErrCoultNotFetchRobotsTxt = errors.New("could not fetch robots.txt")

type Repository interface {
	Get(hostname string) (string, error)
	Set(hostname, robotsTxt string) error
}

type Service interface {
	CanCrawlURL(url string) bool
	FetchRobotsTxt(hostname string) error
}
