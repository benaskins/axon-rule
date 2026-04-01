package spec

// evaluator is implemented by composite specs that can produce their own
// Result by recursing into child specs.
type evaluator[T any] interface {
	evaluate(candidate T) Result
}

// Evaluate runs all specs against a candidate and returns a Result containing
// violations from every failing spec. It never short-circuits.
func Evaluate[T any](candidate T, specs ...Spec[T]) Result {
	var violations []Violation
	for _, s := range specs {
		violations = append(violations, collect(candidate, s)...)
	}
	return Result{Violations: violations}
}

func collect[T any](candidate T, s Spec[T]) []Violation {
	if e, ok := s.(evaluator[T]); ok {
		r := e.evaluate(candidate)
		return r.Violations
	}

	pr := s.IsSatisfiedBy(candidate)
	if pr.OK {
		return nil
	}
	return []Violation{{Code: s.Code(), Context: pr.Context}}
}
