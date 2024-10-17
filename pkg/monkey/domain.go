package monkey

type Service interface {
	Execute(src string) ([]byte, error)
}
