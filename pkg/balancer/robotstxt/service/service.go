package service

import (
	"fmt"
	"io"
	"juno/pkg/balancer/robotstxt"
	"net/http"

	junourl "juno/pkg/url"

	temtorobots "github.com/temoto/robotstxt"
)

type Service struct {
	repo robotstxt.Repository
}

func New(repo robotstxt.Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func isAllowed(robotsTxt string, url string) bool {

	robots, err := temtorobots.FromString(robotsTxt)

	if err != nil {
		return false
	}

	agent := "JunoBot/1.0"

	return robots.TestAgent(url, agent)
}

func (s *Service) fetchRobotsTxt(https bool, hostname string) (string, error) {

	proto := "http"

	if https {
		proto = "https"
	}

	robotsURL := fmt.Sprintf("%s://%s/robots.txt", proto, hostname)

	fmt.Println(robotsURL)

	// try both http and https
	res, err := http.Get(robotsURL)

	if err != nil {
		return "", robotstxt.ErrCoultNotFetchRobotsTxt
	}

	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)

	if err != nil {
		return "", robotstxt.ErrCoultNotFetchRobotsTxt
	}

	return string(data), nil
}

func (s *Service) CanCrawlURL(url string) bool {

	hostname, err := junourl.ToHostname(url)

	if err != nil {
		return false
	}

	rtxt, err := s.repo.Get(hostname)
	if err == robotstxt.ErrRobotsTxtNotInCache {
		isHTTPS := junourl.IsHTTPS(url)
		rtxt, _ = s.fetchRobotsTxt(isHTTPS, hostname)

		err = s.repo.Set(hostname, rtxt)

		if err != nil {
			return false
		}
	}

	return isAllowed(rtxt, url)
}
