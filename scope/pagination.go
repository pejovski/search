package scope

import (
	"net/http"
	"strconv"
)

const (
	PaginationOffsetDefault = 0
	PaginationLimitDefault  = 10
	RequestParameterOffset  = "offset"
	RequestParameterLimit   = "limit"
)

type Pagination struct {
	Offset int
	Limit  int
}

func NewPagination(r *http.Request) *Pagination {
	o, l := PaginationOffsetDefault, PaginationLimitDefault

	if ro := r.URL.Query().Get(RequestParameterOffset); len(ro) > 0 {
		o, _ = strconv.Atoi(ro)
	}
	if rl := r.URL.Query().Get(RequestParameterLimit); len(rl) > 0 {
		l, _ = strconv.Atoi(rl)
	}

	return &Pagination{
		Offset: o,
		Limit:  l,
	}
}
