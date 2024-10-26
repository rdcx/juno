package ranag

import "juno/pkg/ranag/dto"

type Service interface {
	QueryRange(offset, total int, req dto.RangeAggregatorRequest) (interface{}, error)
}
