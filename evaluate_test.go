package rule_test

import (
	"testing"

	"github.com/benaskins/axon-rule"
)

func TestEvaluate_AllPass(t *testing.T) {
	rules := rule.AllOf(
		rule.New(order.HasCustomer),
		rule.New(order.HasItems),
	)

	result := rules.Evaluate(order{CustomerID: "c1", Items: []string{"x"}, Total: 100})
	if !result.IsValid() {
		t.Fatal("all rules pass, result should be valid")
	}
}

func TestEvaluate_CollectsAllViolations(t *testing.T) {
	rules := rule.AllOf(
		rule.New(order.HasCustomer),
		rule.New(order.HasItems),
		rule.New(order.HasPositiveTotal),
	)

	result := rules.Evaluate(order{})
	if result.IsValid() {
		t.Fatal("result should not be valid")
	}
	if len(result.Items) != 3 {
		t.Fatalf("got %d violations, want 3", len(result.Items))
	}

	codes := result.Codes()
	want := []string{"MissingCustomer", "EmptyItems", "NegativeTotal"}
	for i, c := range codes {
		if c != want[i] {
			t.Errorf("violation %d: got %q, want %q", i, c, want[i])
		}
	}
}

func TestEvaluate_PreservesContext(t *testing.T) {
	rules := rule.AllOf(
		rule.New(order.HasPositiveTotal),
	)

	result := rules.Evaluate(order{Total: -50})
	if len(result.Items) != 1 {
		t.Fatalf("got %d violations, want 1", len(result.Items))
	}
	v := result.Items[0]
	if v.Code != "NegativeTotal" {
		t.Errorf("got code %q, want %q", v.Code, "NegativeTotal")
	}
	ctx, ok := v.Context.(NegativeTotal)
	if !ok {
		t.Fatalf("expected NegativeTotal context, got %T", v.Context)
	}
	if ctx.Total != int64(-50) {
		t.Errorf("got total %v, want -50", ctx.Total)
	}
}

func TestEvaluate_AllOf_CollectsChildViolations(t *testing.T) {
	rules := rule.AllOf(
		rule.New(order.HasCustomer),
		rule.New(order.HasItems),
	)

	result := rules.Evaluate(order{})
	if len(result.Items) != 2 {
		t.Fatalf("got %d violations, want 2", len(result.Items))
	}
	codes := result.Codes()
	if codes[0] != "MissingCustomer" || codes[1] != "EmptyItems" {
		t.Errorf("got codes %v", codes)
	}
}

func TestEvaluate_AnyOf_NoViolationsOnPass(t *testing.T) {
	rules := rule.AllOf(
		rule.AnyOf(
			rule.New(order.HasCustomer),
			rule.New(order.HasItems),
		),
	)

	result := rules.Evaluate(order{CustomerID: "c1"})
	if !result.IsValid() {
		t.Fatal("AnyOf with one passing rule should produce no violations")
	}
}

func TestEvaluate_AnyOf_CollectsAllOnFailure(t *testing.T) {
	rules := rule.AllOf(
		rule.AnyOf(
			rule.New(order.HasCustomer),
			rule.New(order.HasItems),
		),
	)

	result := rules.Evaluate(order{})
	if len(result.Items) != 2 {
		t.Fatalf("got %d violations, want 2", len(result.Items))
	}
}

func TestEvaluate_Not_ViolationOnSatisfied(t *testing.T) {
	rules := rule.AllOf(
		rule.Not(rule.New(order.HasCustomer)),
	)

	result := rules.Evaluate(order{CustomerID: "c1"})
	if len(result.Items) != 1 {
		t.Fatalf("got %d violations, want 1", len(result.Items))
	}
	if result.Items[0].Code != "Negated" {
		t.Errorf("got code %q, want %q", result.Items[0].Code, "Negated")
	}
}

func TestEvaluate_Not_NoViolationOnFailed(t *testing.T) {
	rules := rule.AllOf(
		rule.Not(rule.New(order.HasCustomer)),
	)

	result := rules.Evaluate(order{})
	if !result.IsValid() {
		t.Fatal("Not should produce no violations when inner rule fails")
	}
}

func TestEvaluate_NestedComposition(t *testing.T) {
	rules := rule.AllOf(
		rule.New(order.HasCustomer),
		rule.AnyOf(
			rule.New(order.HasItems),
			rule.New(order.HasPositiveTotal),
		),
		rule.Not(rule.New(order.HasCustomer)),
	)

	// Customer present but no items/total, and Not(has-customer) fails
	result := rules.Evaluate(order{CustomerID: "c1"})
	codes := result.Codes()
	// AnyOf fails (EmptyItems + NegativeTotal), Not fails (Negated)
	if len(codes) != 3 {
		t.Fatalf("got %d violations %v, want 3", len(codes), codes)
	}
}

func TestEvaluate_NoRules(t *testing.T) {
	rules := rule.AllOf[order]()
	result := rules.Evaluate(order{})
	if !result.IsValid() {
		t.Fatal("no rules means no violations")
	}
}
