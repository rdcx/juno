package dto

type CrawlRequest struct {
	URL string `json:"url"`
}

type CrawlResponse struct {
	Status string `json:"status"`
}
