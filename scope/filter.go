package scope

import (
	"net/http"
	"strconv"
	"strings"
)

const (
	FilterKeyBrand = "brand"
	FilterKeyPrice = "price"
)

const (
	FilterOperatorEqual            = "eq"
	FilterOperatorGreaterThanEqual = "gte"
	FilterOperatorLowerThanEqual   = "lte"
)

var AvailableFilterKeys = []string{
	FilterKeyBrand,
	FilterKeyPrice,
}

type Filters []*Filter

type Filter struct {
	Key      string
	Value    interface{}
	Operator string
}

func NewFilters(r *http.Request) Filters {
	var f []*Filter

	for _, key := range AvailableFilterKeys {

		value := r.URL.Query().Get(key)
		if value == "" {
			continue
		}

		switch key {
		case FilterKeyBrand:
			f = append(f, &Filter{
				Key:      key,
				Value:    value,
				Operator: FilterOperatorEqual,
			})
		case FilterKeyPrice:
			i := strings.Index(value, "-")
			if i == -1 {
				continue
			}

			gte, err := strconv.Atoi(value[:i])
			if err == nil {
				f = append(f, &Filter{
					Key:      key,
					Value:    gte,
					Operator: FilterOperatorGreaterThanEqual,
				})
			}

			lte, err := strconv.Atoi(value[i+1:])
			if err == nil {
				f = append(f, &Filter{
					Key:      key,
					Value:    lte,
					Operator: FilterOperatorLowerThanEqual,
				})
			}
		}
	}

	if len(f) == 0 {
		return nil
	}

	return f
}
