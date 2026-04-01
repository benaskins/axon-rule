package rule

func collect[T any](candidate T, r Rule[T]) []Violation {
	type evaluator[U any] interface {
		Evaluate(candidate U) Violations
	}

	if e, ok := Rule[T](r).(evaluator[T]); ok {
		return e.Evaluate(candidate).Items
	}

	v := r.Check(candidate)
	if v.OK {
		return nil
	}
	return []Violation{{Code: r.Code(), Context: v.Context}}
}
