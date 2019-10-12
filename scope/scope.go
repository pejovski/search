package scope

type Scope struct {
	SearchQuery string
	Pagination  *Pagination
	Filters     Filters
	Sorting     *Sorting
}

func New(sq string, p *Pagination, f Filters, s *Sorting) *Scope {
	return &Scope{
		SearchQuery: sq,
		Pagination:  p,
		Filters:     f,
		Sorting:     s,
	}
}
