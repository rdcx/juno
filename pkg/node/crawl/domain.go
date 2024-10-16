package crawl

import "errors"

var (
	Err500         = errors.New("fetch returned 500")
	Err429         = errors.New("fetch returned 429")
	Err404         = errors.New("fetch returned 404")
	Err400         = errors.New("fetch returned 400")
	ErrContextDone = errors.New("context was canceled or timed out")
)
var ErrFailedCrawlRequest = errors.New("failed to send crawl request")

type Service interface {
	Crawl(url string) error
}
