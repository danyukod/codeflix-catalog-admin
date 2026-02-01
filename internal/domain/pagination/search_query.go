package pagination

// SearchQuery is the Go equivalent of the Java record SearchQuery.
type SearchQuery struct {
	Page      int
	PerPage   int
	Terms     string
	Sort      string
	Direction string
}

// NewSearchQuery is an optional helper constructor.
func NewSearchQuery(page, perPage int, terms, sort, direction string) SearchQuery {
	return SearchQuery{
		Page:      page,
		PerPage:   perPage,
		Terms:     terms,
		Sort:      sort,
		Direction: direction,
	}
}
