package balancer

type Service interface {
	SendCrawlRequest(url string) error
	SendBatchedLinks(links []string) error
	ReportURLProcessed(url string, status int) error
}
