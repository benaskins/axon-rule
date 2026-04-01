package rule

// Rule defines a single business rule for a candidate of type T.
type Rule[T any] interface {
	Check(candidate T) Verdict
}

// funcRule implements Rule using a function.
type funcRule[T any] struct {
	fn func(T) Verdict
}

// New creates a Rule from a predicate function.
func New[T any](fn func(T) Verdict) Rule[T] {
	return funcRule[T]{fn: fn}
}

func (s funcRule[T]) Check(candidate T) Verdict {
	return s.fn(candidate)
}
