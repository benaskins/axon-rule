package spec

// AllOf returns a Spec that is satisfied only when every inner spec passes.
// All specs are evaluated regardless of failures; violations are collected
// from every failing spec.
func AllOf[T any](specs ...Spec[T]) Spec[T] {
	return allOf[T]{specs: specs}
}

type allOf[T any] struct {
	specs []Spec[T]
}

func (a allOf[T]) Code() Code { return "all-of" }

func (a allOf[T]) IsSatisfiedBy(candidate T) PredicateResult {
	satisfied := true
	for _, s := range a.specs {
		r := s.IsSatisfiedBy(candidate)
		if !r.OK {
			satisfied = false
		}
	}
	if satisfied {
		return Pass()
	}
	return Fail()
}

func (a allOf[T]) evaluate(candidate T) Result {
	var violations []Violation
	for _, s := range a.specs {
		violations = append(violations, collect(candidate, s)...)
	}
	return Result{Violations: violations}
}

// AnyOf returns a Spec that is satisfied when at least one inner spec passes.
// Short-circuits on the first success.
func AnyOf[T any](specs ...Spec[T]) Spec[T] {
	return anyOf[T]{specs: specs}
}

type anyOf[T any] struct {
	specs []Spec[T]
}

func (a anyOf[T]) Code() Code { return "any-of" }

func (a anyOf[T]) IsSatisfiedBy(candidate T) PredicateResult {
	for _, s := range a.specs {
		r := s.IsSatisfiedBy(candidate)
		if r.OK {
			return Pass()
		}
	}
	return Fail()
}

func (a anyOf[T]) evaluate(candidate T) Result {
	var violations []Violation
	for _, s := range a.specs {
		vs := collect(candidate, s)
		if len(vs) == 0 {
			return Result{}
		}
		violations = append(violations, vs...)
	}
	return Result{Violations: violations}
}

// Not returns a Spec that inverts another spec. It is satisfied when the
// inner spec fails. Produces a violation with "not:" prefixed to the inner
// spec's code.
func Not[T any](s Spec[T]) Spec[T] {
	return notSpec[T]{inner: s}
}

type notSpec[T any] struct {
	inner Spec[T]
}

func (n notSpec[T]) Code() Code {
	return "not:" + n.inner.Code()
}

func (n notSpec[T]) IsSatisfiedBy(candidate T) PredicateResult {
	r := n.inner.IsSatisfiedBy(candidate)
	if r.OK {
		return Fail()
	}
	return Pass()
}

func (n notSpec[T]) evaluate(candidate T) Result {
	r := n.inner.IsSatisfiedBy(candidate)
	if !r.OK {
		return Result{}
	}
	return Result{Violations: []Violation{{Code: n.Code()}}}
}
