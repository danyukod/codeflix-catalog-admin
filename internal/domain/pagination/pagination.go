package pagination

// Pagination is the Go equivalent of the Java record Pagination<T>.
type Pagination[T any] struct {
	CurrentPage int
	PerPage     int
	Total       int64
	Items       []T
}

// Map converts Pagination[T] into Pagination[R] by mapping each item using mapper.
// Note: In Go, this must be a function (not a method) because Go doesn't allow
// methods with additional type parameters (like [R any]) beyond the receiver's.
func Map[T any, R any](p Pagination[T], mapper func(T) R) Pagination[R] {
	newItems := make([]R, len(p.Items))
	for i, item := range p.Items {
		newItems[i] = mapper(item)
	}

	return Pagination[R]{
		CurrentPage: p.CurrentPage,
		PerPage:     p.PerPage,
		Total:       p.Total,
		Items:       newItems,
	}
}

// New is an optional helper constructor.
func New[T any](currentPage, perPage int, total int64, items []T) Pagination[T] {
	return Pagination[T]{
		CurrentPage: currentPage,
		PerPage:     perPage,
		Total:       total,
		Items:       items,
	}
}
