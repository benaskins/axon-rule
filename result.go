package spec

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
