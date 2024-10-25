package ranag

type Service interface {
	QueryRange([]int) (interface{}, error)
}
