package rule

// Rule defines a single named business rule for a candidate of type T.
type Rule[T any] interface {
	Code() Code
	Check(candidate T) Verdict
}

// funcRule implements Rule using a function.
type funcRule[T any] struct {
	code Code
	fn   func(T) Verdict
}

// New creates a Rule from a code and a predicate function.
func New[T any](code Code, fn func(T) Verdict) Rule[T] {
	return funcRule[T]{code: code, fn: fn}
}

func (s funcRule[T]) Code() Code {
	return s.code
}

func (s funcRule[T]) Check(candidate T) Verdict {
	return s.fn(candidate)
}
