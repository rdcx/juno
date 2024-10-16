package balancer

type Service interface {
	ReportURLFound(url string) error
	ReportURLProcessed(url string, status int) error
}
