package rule

// AllOf returns a Rule that is satisfied only when every inner rule passes.
// All rules are evaluated regardless of failures; violations are collected
// from every failing rule.
func AllOf[T any](rules ...Rule[T]) allOf[T] {
	return allOf[T]{rules: rules}
}

type allOf[T any] struct {
	rules []Rule[T]
}

func (a allOf[T]) Code() Code { return "all-of" }

func (a allOf[T]) Check(candidate T) Verdict {
	for _, r := range a.rules {
		v := r.Check(candidate)
		if !v.OK {
			return Fail()
		}
	}
	return Pass()
}

// Evaluate runs all rules and collects violations from every failing rule.
func (a allOf[T]) Evaluate(candidate T) Violations {
	var items []Violation
	for _, r := range a.rules {
		items = append(items, collect(candidate, r)...)
	}
	return Violations{Items: items}
}

// AnyOf returns a Rule that is satisfied when at least one inner rule passes.
// Short-circuits on the first success.
func AnyOf[T any](rules ...Rule[T]) anyOf[T] {
	return anyOf[T]{rules: rules}
}

type anyOf[T any] struct {
	rules []Rule[T]
}

func (a anyOf[T]) Code() Code { return "any-of" }

func (a anyOf[T]) Check(candidate T) Verdict {
	for _, r := range a.rules {
		v := r.Check(candidate)
		if v.OK {
			return Pass()
		}
	}
	return Fail()
}

// Evaluate returns no violations if any rule passes; otherwise collects all.
func (a anyOf[T]) Evaluate(candidate T) Violations {
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

// Not returns a Rule that inverts another rule. It is satisfied when the
// inner rule fails. Produces a violation with "not:" prefixed to the inner
// rule's code.
func Not[T any](r Rule[T]) notRule[T] {
	return notRule[T]{inner: r}
}

type notRule[T any] struct {
	inner Rule[T]
}

func (n notRule[T]) Code() Code {
	return "not:" + n.inner.Code()
}

func (n notRule[T]) Check(candidate T) Verdict {
	v := n.inner.Check(candidate)
	if v.OK {
		return Fail()
	}
	return Pass()
}

// Evaluate produces a violation when the inner rule passes (i.e. Not fails).
func (n notRule[T]) Evaluate(candidate T) Violations {
	v := n.inner.Check(candidate)
	if !v.OK {
		return Violations{}
	}
	return Violations{Items: []Violation{{Code: n.Code()}}}
}
