package ranag

import (
	"juno/pkg/api/query/dto"
)

type Service interface {
	QueryRange([]int, dto.Query) (interface{}, error)
}
