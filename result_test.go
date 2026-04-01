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
			{Code: "something-failed"},
		},
	}
	if v.IsValid() {
		t.Fatal("violations with items should not be valid")
	}
}

func TestViolations_Codes(t *testing.T) {
	v := rule.Violations{
		Items: []rule.Violation{
			{Code: "first"},
			{Code: "second"},
		},
	}

	codes := v.Codes()
	if len(codes) != 2 {
		t.Fatalf("got %d codes, want 2", len(codes))
	}
	if codes[0] != "first" || codes[1] != "second" {
		t.Errorf("got codes %v, want [first second]", codes)
	}
}

func TestViolations_Codes_Empty(t *testing.T) {
	v := rule.Violations{}
	codes := v.Codes()
	if len(codes) != 0 {
		t.Fatalf("got %d codes, want 0", len(codes))
	}
}

func TestViolationContext(t *testing.T) {
	v := rule.Violation{
		Code: "debits-must-equal-credits",
		Context: map[string]any{
			"total_debits":  int64(5000),
			"total_credits": int64(3000),
			"difference":    int64(2000),
		},
	}

	if v.Code != "debits-must-equal-credits" {
		t.Errorf("got code %q", v.Code)
	}
	ctx, ok := v.Context.(map[string]any)
	if !ok {
		t.Fatalf("expected map[string]any context, got %T", v.Context)
	}
	if ctx["difference"] != int64(2000) {
		t.Errorf("got difference %v", ctx["difference"])
	}
}
