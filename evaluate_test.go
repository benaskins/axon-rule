package rule_test

import (
	"testing"

	"github.com/benaskins/axon-rule"
)

func TestEvaluate_AllPass(t *testing.T) {
	rules := rule.AllOf(
		rule.New("has-customer", order.HasCustomer),
		rule.New("has-items", order.HasItems),
	)

	result := rules.Evaluate(order{CustomerID: "c1", Items: []string{"x"}, Total: 100})
	if !result.IsValid() {
		t.Fatal("all rules pass, result should be valid")
	}
}

func TestEvaluate_CollectsAllViolations(t *testing.T) {
	rules := rule.AllOf(
		rule.New("has-customer", order.HasCustomer),
		rule.New("has-items", order.HasItems),
		rule.New("has-positive-total", order.HasPositiveTotal),
	)

	result := rules.Evaluate(order{})
	if result.IsValid() {
		t.Fatal("result should not be valid")
	}
	if len(result.Items) != 3 {
		t.Fatalf("got %d violations, want 3", len(result.Items))
	}

	codes := result.Codes()
	want := []rule.Code{"has-customer", "has-items", "has-positive-total"}
	for i, c := range codes {
		if c != want[i] {
			t.Errorf("violation %d: got %q, want %q", i, c, want[i])
		}
	}
}

func TestEvaluate_PreservesContext(t *testing.T) {
	rules := rule.AllOf(
		rule.New("has-positive-total", order.HasPositiveTotal),
	)

	result := rules.Evaluate(order{Total: -50})
	if len(result.Items) != 1 {
		t.Fatalf("got %d violations, want 1", len(result.Items))
	}
	v := result.Items[0]
	if v.Context["total"] != int64(-50) {
		t.Errorf("got total %v, want -50", v.Context["total"])
	}
}

func TestEvaluate_AllOf_CollectsChildViolations(t *testing.T) {
	rules := rule.AllOf(
		rule.New("has-customer", order.HasCustomer),
		rule.New("has-items", order.HasItems),
	)

	result := rules.Evaluate(order{})
	if len(result.Items) != 2 {
		t.Fatalf("got %d violations, want 2", len(result.Items))
	}
	codes := result.Codes()
	if codes[0] != "has-customer" || codes[1] != "has-items" {
		t.Errorf("got codes %v", codes)
	}
}

func TestEvaluate_AnyOf_NoViolationsOnPass(t *testing.T) {
	rules := rule.AllOf(
		rule.AnyOf(
			rule.New("has-customer", order.HasCustomer),
			rule.New("has-items", order.HasItems),
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
			rule.New("has-customer", order.HasCustomer),
			rule.New("has-items", order.HasItems),
		),
	)

	result := rules.Evaluate(order{})
	if len(result.Items) != 2 {
		t.Fatalf("got %d violations, want 2", len(result.Items))
	}
}

func TestEvaluate_Not_ViolationOnSatisfied(t *testing.T) {
	rules := rule.AllOf(
		rule.Not(rule.New("has-customer", order.HasCustomer)),
	)

	result := rules.Evaluate(order{CustomerID: "c1"})
	if len(result.Items) != 1 {
		t.Fatalf("got %d violations, want 1", len(result.Items))
	}
	if result.Items[0].Code != "not:has-customer" {
		t.Errorf("got code %q, want %q", result.Items[0].Code, "not:has-customer")
	}
}

func TestEvaluate_Not_NoViolationOnFailed(t *testing.T) {
	rules := rule.AllOf(
		rule.Not(rule.New("has-customer", order.HasCustomer)),
	)

	result := rules.Evaluate(order{})
	if !result.IsValid() {
		t.Fatal("Not should produce no violations when inner rule fails")
	}
}

func TestEvaluate_NestedComposition(t *testing.T) {
	rules := rule.AllOf(
		rule.New("has-customer", order.HasCustomer),
		rule.AnyOf(
			rule.New("has-items", order.HasItems),
			rule.New("has-positive-total", order.HasPositiveTotal),
		),
		rule.Not(rule.New("has-customer", order.HasCustomer)),
	)

	// Customer present but no items/total, and Not(has-customer) fails
	result := rules.Evaluate(order{CustomerID: "c1"})
	codes := result.Codes()
	// AnyOf fails (has-items + has-positive-total), Not fails (not:has-customer)
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
