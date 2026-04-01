package rule

// Verdict is the return value of a predicate function.
type Verdict struct {
	OK      bool
	Context map[string]any
}

// Pass returns a successful verdict.
func Pass() Verdict {
	return Verdict{OK: true}
}

// Fail returns a failed verdict with no context.
func Fail() Verdict {
	return Verdict{OK: false}
}

// FailWith returns a failed verdict with context.
func FailWith(context map[string]any) Verdict {
	return Verdict{OK: false, Context: context}
}

// Violation records a rule that was not satisfied.
type Violation struct {
	Code    Code
	Context map[string]any
}

// Violations holds the outcome of evaluating rules against a candidate.
type Violations struct {
	Items []Violation
}

// IsValid returns true when no violations were recorded.
func (v Violations) IsValid() bool {
	return len(v.Items) == 0
}

// Codes returns the codes of all recorded violations.
func (v Violations) Codes() []Code {
	codes := make([]Code, len(v.Items))
	for i, item := range v.Items {
		codes[i] = item.Code
	}
	return codes
}
