package scope

import (
	"net/http"
	"strings"
)

const RequestParameterSorting = "sort"

const (
	SortingFieldDefault = "price"
	SortingOrderDefault = "asc"
)

type Sorting struct {
	Field string
	Order string
}

func NewSorting(r *http.Request) *Sorting {
	soring := &Sorting{
		Field: SortingFieldDefault,
		Order: SortingOrderDefault,
	}

	sp := r.URL.Query().Get(RequestParameterSorting)
	if sp == "" {
		return nil
	}

	i := strings.Index(sp, "-")
	if i == -1 {
		return nil
	}

	soring.Field, soring.Order = strings.ToLower(sp[:i]), strings.ToLower(sp[i+1:])

	return soring
}
