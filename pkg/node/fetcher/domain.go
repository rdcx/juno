package fetcher

import "context"

type Service interface {
	FetchPage(ctx context.Context, url string) (body []byte, status int, finalURL string, err error)
}
