package rule

import "reflect"

// Verdict is the return value of a predicate function.
type Verdict struct {
	OK      bool
	Context any
}

// Pass returns a successful verdict.
func Pass() Verdict {
	return Verdict{OK: true}
}

// FailWith returns a failed verdict with a typed context.
// The context type name becomes the violation code.
func FailWith(context any) Verdict {
	return Verdict{OK: false, Context: context}
}

// Violation records a rule that was not satisfied.
// Code is derived from the type name of Context.
type Violation struct {
	Code    string
	Context any
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
func (v Violations) Codes() []string {
	codes := make([]string, len(v.Items))
	for i, item := range v.Items {
		codes[i] = item.Code
	}
	return codes
}

// codeName returns the type name of a context value for use as a violation code.
func codeName(context any) string {
	if context == nil {
		return "unknown"
	}
	return reflect.TypeOf(context).Name()
}
