package dto

const (
	SUCCESS = "success"
	OK      = "ok"
	ERROR   = "error"
)

type CrawlRequest struct {
	URL string `json:"url"`
}

type CrawlResponse struct {
	Status string `json:"status"`
}

func NewOKCrawlResponse() CrawlResponse {
	return CrawlResponse{
		Status: OK,
	}
}

func NewErrorCrawlResponse() CrawlResponse {
	return CrawlResponse{
		Status: ERROR,
	}
}

type CrawlURLsRequest struct {
	URLs []string `json:"urls"`
}

type CrawlURLsResponse struct {
	Status string `json:"status"`
}

func NewOKCrawlURLsResponse() CrawlURLsResponse {
	return CrawlURLsResponse{
		Status: OK,
	}
}

func NewErrorCrawlURLsResponse() CrawlURLsResponse {
	return CrawlURLsResponse{
		Status: ERROR,
	}
}
