package spec

// Spec defines a single named business rule for a candidate of type T.
type Spec[T any] interface {
	Code() Code
	IsSatisfiedBy(candidate T) PredicateResult
}

// funcSpec implements Spec using a function.
type funcSpec[T any] struct {
	code Code
	fn   func(T) PredicateResult
}

// New creates a Spec from a code and a predicate function.
func New[T any](code Code, fn func(T) PredicateResult) Spec[T] {
	return funcSpec[T]{code: code, fn: fn}
}

func (s funcSpec[T]) Code() Code {
	return s.code
}

func (s funcSpec[T]) IsSatisfiedBy(candidate T) PredicateResult {
	return s.fn(candidate)
}
