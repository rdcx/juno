package balancer

type Service interface {
	SendCrawlRequest(url string) error
	ReportURLProcessed(url string, status int) error
}
