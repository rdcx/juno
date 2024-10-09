package crawler

import (
	"context"
	"errors"
	"io"
	"net/http"
)

var (
	Err500         = errors.New("fetch returned 500")
	Err429         = errors.New("fetch returned 429")
	Err404         = errors.New("fetch returned 404")
	Err400         = errors.New("fetch returned 400")
	ErrContextDone = errors.New("context was canceled or timed out")
)

func FetchPage(ctx context.Context, url string) (
	status int,
	finalURL string,
	body []byte,
	err error,
) {
	// Create a new HTTP request with the provided context
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return 0, "", nil, err
	}

	// Make the HTTP request
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		// Check if the context was canceled or timed out
		if ctx.Err() != nil {
			err = ErrContextDone
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
		err = Err500
	case 429:
		err = Err429
	case 404:
		err = Err404
	case 400:
		err = Err400
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
