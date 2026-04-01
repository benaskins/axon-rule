package rule_test

import (
	"testing"

	"github.com/benaskins/axon-rule"
)

func TestViolations_IsValid_NoViolations(t *testing.T) {
	v := rule.Violations{}
	if !v.IsValid() {
		t.Fatal("empty violations should be valid")
	}
}

func TestViolations_IsValid_WithViolations(t *testing.T) {
	v := rule.Violations{
		Items: []rule.Violation{
			{Code: "SomethingFailed"},
		},
	}
	if v.IsValid() {
		t.Fatal("violations with items should not be valid")
	}
}

func TestViolations_Codes(t *testing.T) {
	v := rule.Violations{
		Items: []rule.Violation{
			{Code: "First"},
			{Code: "Second"},
		},
	}

	codes := v.Codes()
	if len(codes) != 2 {
		t.Fatalf("got %d codes, want 2", len(codes))
	}
	if codes[0] != "First" || codes[1] != "Second" {
		t.Errorf("got codes %v, want [First Second]", codes)
	}
}

func TestViolations_Codes_Empty(t *testing.T) {
	v := rule.Violations{}
	codes := v.Codes()
	if len(codes) != 0 {
		t.Fatalf("got %d codes, want 0", len(codes))
	}
}

func TestViolation_CodeDerivedFromContextType(t *testing.T) {
	type InsufficientFunds struct {
		Available int64
		Required  int64
	}

	v := rule.FailWith(InsufficientFunds{Available: 3000, Required: 5000})
	if v.Context == nil {
		t.Fatal("expected context")
	}
	ctx, ok := v.Context.(InsufficientFunds)
	if !ok {
		t.Fatalf("expected InsufficientFunds, got %T", v.Context)
	}
	if ctx.Available != 3000 || ctx.Required != 5000 {
		t.Errorf("got %+v", ctx)
	}
}
