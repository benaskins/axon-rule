package rule

// AllOf returns a Rule that is satisfied only when every inner rule passes.
// All rules are evaluated regardless of failures; violations are collected
// from every failing rule.
func AllOf[T any](rules ...Rule[T]) AllOfRule[T] {
	return AllOfRule[T]{rules: rules}
}

type AllOfRule[T any] struct {
	rules []Rule[T]
}

func (a AllOfRule[T]) Check(candidate T) Verdict {
	for _, r := range a.rules {
		v := r.Check(candidate)
		if !v.OK {
			return v
		}
	}
	return Pass()
}

// Evaluate runs all rules and collects violations from every failing rule.
func (a AllOfRule[T]) Evaluate(candidate T) Violations {
	var items []Violation
	for _, r := range a.rules {
		items = append(items, collect(candidate, r)...)
	}
	return Violations{Items: items}
}

// AnyOf returns a Rule that is satisfied when at least one inner rule passes.
// Short-circuits on the first success.
func AnyOf[T any](rules ...Rule[T]) AnyOfRule[T] {
	return AnyOfRule[T]{rules: rules}
}

type AnyOfRule[T any] struct {
	rules []Rule[T]
}

func (a AnyOfRule[T]) Check(candidate T) Verdict {
	for _, r := range a.rules {
		v := r.Check(candidate)
		if v.OK {
			return Pass()
		}
	}
	return FailWith(nil)
}

// Evaluate returns no violations if any rule passes; otherwise collects all.
func (a AnyOfRule[T]) Evaluate(candidate T) Violations {
	var items []Violation
	for _, r := range a.rules {
		vs := collect(candidate, r)
		if len(vs) == 0 {
			return Violations{}
		}
		items = append(items, vs...)
	}
	return Violations{Items: items}
}

// Negated is the violation context produced when a Not rule fails.
type Negated struct{}

// Not returns a Rule that inverts another rule. It is satisfied when the
// inner rule fails. Produces a Negated violation when the inner rule passes.
func Not[T any](r Rule[T]) NotRule[T] {
	return NotRule[T]{inner: r}
}

type NotRule[T any] struct {
	inner Rule[T]
}

func (n NotRule[T]) Check(candidate T) Verdict {
	v := n.inner.Check(candidate)
	if v.OK {
		return FailWith(Negated{})
	}
	return Pass()
}

// Evaluate produces a violation when the inner rule passes (i.e. Not fails).
func (n NotRule[T]) Evaluate(candidate T) Violations {
	v := n.inner.Check(candidate)
	if !v.OK {
		return Violations{}
	}
	neg := Negated{}
	return Violations{Items: []Violation{{Code: codeName(neg), Context: neg}}}
}
