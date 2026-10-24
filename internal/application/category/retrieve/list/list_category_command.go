package list

import "github.com/danyukod/codeflix-catalog-admin/internal/domain/pagination"

type CategoryCommand struct {
	page      int
	perPage   int
	terms     string
	sort      string
	direction string
}

func (c CategoryCommand) With(page, perPage int, terms, sort, direction string) *CategoryCommand {
	return &CategoryCommand{page, perPage, terms, sort, direction}
}

func (c *CategoryCommand) ToSearchQuery() pagination.SearchQuery {
	return pagination.NewSearchQuery(c.page, c.perPage, c.terms, c.sort, c.direction)
}
