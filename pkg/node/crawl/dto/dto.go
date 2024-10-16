package dto

const (
	SUCCESS = "success"
	ERROR   = "error"
)

type CrawlRequest struct {
	URL string `json:"url"`
}

type CrawlResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

func NewSuccessCrawlResponse() *CrawlResponse {
	return &CrawlResponse{
		Status: SUCCESS,
	}
}

func NewErrorCrawlResponse(err string) *CrawlResponse {
	return &CrawlResponse{
		Status:  ERROR,
		Message: err,
	}
}
