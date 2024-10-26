package html

type Service interface {
	ExtractLinks(body []byte) ([]string, error)
	Title(body []byte) (string, error)
	GetSelectorValue(body []byte, selector string) (string, error)
}
