package service

import (
	"context"
	"errors"
	"io"
	"juno/pkg/node/crawl"
	"net/http"
)

type Service struct{}

func New() *Service {
	return &Service{}
}

func (s *Service) FetchPage(ctx context.Context, url string) (
	body []byte,
	status int,
	finalURL string,
	err error,
) {
	// Create a new HTTP request with the provided context
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return
	}

	// Make the HTTP request
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		// Check if the context was canceled or timed out
		if ctx.Err() != nil {
			err = crawl.ErrContextDone
		}
		return
	}
	defer res.Body.Close()

	// Get the status code and final URL after all redirects
	status = res.StatusCode
	finalURL = res.Request.URL.String()

	// Check for non-2xx status codes and return specific errors
	switch status {
	case 500:
		err = crawl.Err500
	case 429:
		err = crawl.Err429
	case 404:
		err = crawl.Err404
	case 400:
		err = crawl.Err400
	default:
		if status < 200 || status >= 300 {
			err = errors.New("unexpected status code: " + http.StatusText(status))
			return
		}
	}

	if err != nil {
		return
	}

	// Read the response body
	body, err = io.ReadAll(res.Body)
	if err != nil {
		return
	}

	return
}
