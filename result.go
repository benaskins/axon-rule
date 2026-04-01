package spec

// PredicateResult is the return value of a predicate function.
type PredicateResult struct {
	OK      bool
	Context map[string]any
}

// Pass returns a successful predicate result.
func Pass() PredicateResult {
	return PredicateResult{OK: true}
}

// Fail returns a failed predicate result with no context.
func Fail() PredicateResult {
	return PredicateResult{OK: false}
}

// FailWith returns a failed predicate result with context.
func FailWith(context map[string]any) PredicateResult {
	return PredicateResult{OK: false, Context: context}
}

// Violation records a spec that was not satisfied.
type Violation struct {
	Code    Code
	Context map[string]any
}

// Result holds the outcome of evaluating specs against a candidate.
type Result struct {
	Violations []Violation
}

// IsValid returns true when no violations were recorded.
func (r Result) IsValid() bool {
	return len(r.Violations) == 0
}

// ViolationCodes returns the codes of all recorded violations.
func (r Result) ViolationCodes() []Code {
	codes := make([]Code, len(r.Violations))
	for i, v := range r.Violations {
		codes[i] = v.Code
	}
	return codes
}
