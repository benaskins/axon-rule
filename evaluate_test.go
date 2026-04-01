package spec_test

import (
	"testing"

	spec "github.com/benaskins/axon-spec"
)

func TestEvaluate_AllPass(t *testing.T) {
	r := spec.Evaluate(
		order{CustomerID: "c1", Items: []string{"x"}, Total: 100},
		spec.New("has-customer", order.HasCustomer),
		spec.New("has-items", order.HasItems),
	)
	if !r.IsValid() {
		t.Fatal("all specs pass, result should be valid")
	}
}

func TestEvaluate_CollectsAllViolations(t *testing.T) {
	r := spec.Evaluate(
		order{},
		spec.New("has-customer", order.HasCustomer),
		spec.New("has-items", order.HasItems),
		spec.New("has-positive-total", order.HasPositiveTotal),
	)

	if r.IsValid() {
		t.Fatal("result should not be valid")
	}
	if len(r.Violations) != 3 {
		t.Fatalf("got %d violations, want 3", len(r.Violations))
	}

	codes := r.ViolationCodes()
	want := []spec.Code{"has-customer", "has-items", "has-positive-total"}
	for i, c := range codes {
		if c != want[i] {
			t.Errorf("violation %d: got %q, want %q", i, c, want[i])
		}
	}
}

func TestEvaluate_PreservesContext(t *testing.T) {
	r := spec.Evaluate(
		order{Total: -50},
		spec.New("has-positive-total", order.HasPositiveTotal),
	)

	if len(r.Violations) != 1 {
		t.Fatalf("got %d violations, want 1", len(r.Violations))
	}
	v := r.Violations[0]
	if v.Context["total"] != int64(-50) {
		t.Errorf("got total %v, want -50", v.Context["total"])
	}
}

func TestEvaluate_AllOf_CollectsChildViolations(t *testing.T) {
	r := spec.Evaluate(
		order{},
		spec.AllOf(
			spec.New("has-customer", order.HasCustomer),
			spec.New("has-items", order.HasItems),
		),
	)

	if len(r.Violations) != 2 {
		t.Fatalf("got %d violations, want 2", len(r.Violations))
	}
	codes := r.ViolationCodes()
	if codes[0] != "has-customer" || codes[1] != "has-items" {
		t.Errorf("got codes %v", codes)
	}
}

func TestEvaluate_AnyOf_NoViolationsOnPass(t *testing.T) {
	r := spec.Evaluate(
		order{CustomerID: "c1"},
		spec.AnyOf(
			spec.New("has-customer", order.HasCustomer),
			spec.New("has-items", order.HasItems),
		),
	)

	if !r.IsValid() {
		t.Fatal("AnyOf with one passing spec should produce no violations")
	}
}

func TestEvaluate_AnyOf_CollectsAllOnFailure(t *testing.T) {
	r := spec.Evaluate(
		order{},
		spec.AnyOf(
			spec.New("has-customer", order.HasCustomer),
			spec.New("has-items", order.HasItems),
		),
	)

	if len(r.Violations) != 2 {
		t.Fatalf("got %d violations, want 2", len(r.Violations))
	}
}

func TestEvaluate_Not_ViolationOnSatisfied(t *testing.T) {
	r := spec.Evaluate(
		order{CustomerID: "c1"},
		spec.Not(spec.New("has-customer", order.HasCustomer)),
	)

	if len(r.Violations) != 1 {
		t.Fatalf("got %d violations, want 1", len(r.Violations))
	}
	if r.Violations[0].Code != "not:has-customer" {
		t.Errorf("got code %q, want %q", r.Violations[0].Code, "not:has-customer")
	}
}

func TestEvaluate_Not_NoViolationOnFailed(t *testing.T) {
	r := spec.Evaluate(
		order{},
		spec.Not(spec.New("has-customer", order.HasCustomer)),
	)

	if !r.IsValid() {
		t.Fatal("Not should produce no violations when inner spec fails")
	}
}

func TestEvaluate_NestedComposition(t *testing.T) {
	isValid := spec.AllOf(
		spec.New("has-customer", order.HasCustomer),
		spec.AnyOf(
			spec.New("has-items", order.HasItems),
			spec.New("has-positive-total", order.HasPositiveTotal),
		),
		spec.Not(spec.New("has-customer", order.HasCustomer)),
	)

	// Customer present but no items/total, and Not(has-customer) fails
	r := spec.Evaluate(order{CustomerID: "c1"}, isValid)

	codes := r.ViolationCodes()
	// AnyOf fails (has-items + has-positive-total), Not fails (not:has-customer)
	if len(codes) != 3 {
		t.Fatalf("got %d violations %v, want 3", len(codes), codes)
	}
}

func TestEvaluate_NoSpecs(t *testing.T) {
	r := spec.Evaluate(order{})
	if !r.IsValid() {
		t.Fatal("no specs means no violations")
	}
}
