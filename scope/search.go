package scope

import (
	"net/http"
	"strings"
)

const RequestParameterQuery = "q"

const SearchTitle = "title"

var SearchFields = []string{SearchTitle}

func NewSearch(r *http.Request) string {
	return strings.Trim(r.URL.Query().Get(RequestParameterQuery), " ")
}
