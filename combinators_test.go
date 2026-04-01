package rule_test

import (
	"testing"

	"github.com/benaskins/axon-rule"
)

func TestAllOf_AllPass(t *testing.T) {
	s := rule.AllOf(
		rule.New(order.HasCustomer),
		rule.New(order.HasItems),
	)

	v := s.Check(order{CustomerID: "c1", Items: []string{"x"}})
	if !v.OK {
		t.Fatal("AllOf should be satisfied when all rules pass")
	}
}

func TestAllOf_OneFails(t *testing.T) {
	s := rule.AllOf(
		rule.New(order.HasCustomer),
		rule.New(order.HasItems),
	)

	v := s.Check(order{CustomerID: "c1"})
	if v.OK {
		t.Fatal("AllOf should not be satisfied when one rule fails")
	}
}

func TestAllOf_AllFail(t *testing.T) {
	s := rule.AllOf(
		rule.New(order.HasCustomer),
		rule.New(order.HasItems),
	)

	v := s.Check(order{})
	if v.OK {
		t.Fatal("AllOf should not be satisfied when all rules fail")
	}
}

func TestAnyOf_OnePass(t *testing.T) {
	s := rule.AnyOf(
		rule.New(order.HasCustomer),
		rule.New(order.HasItems),
	)

	v := s.Check(order{CustomerID: "c1"})
	if !v.OK {
		t.Fatal("AnyOf should be satisfied when at least one rule passes")
	}
}

func TestAnyOf_NonePass(t *testing.T) {
	s := rule.AnyOf(
		rule.New(order.HasCustomer),
		rule.New(order.HasItems),
	)

	v := s.Check(order{})
	if v.OK {
		t.Fatal("AnyOf should not be satisfied when no rules pass")
	}
}

func TestNot_Inverts_Failure(t *testing.T) {
	s := rule.Not(rule.New(order.HasCustomer))

	v := s.Check(order{})
	if !v.OK {
		t.Fatal("Not should be satisfied when inner rule fails")
	}
}

func TestNot_Inverts_Success(t *testing.T) {
	s := rule.Not(rule.New(order.HasCustomer))

	v := s.Check(order{CustomerID: "c1"})
	if v.OK {
		t.Fatal("Not should not be satisfied when inner rule passes")
	}
}

func TestNot_ProducesNegatedViolation(t *testing.T) {
	s := rule.Not(rule.New(order.HasCustomer))

	result := s.Evaluate(order{CustomerID: "c1"})
	if result.IsValid() {
		t.Fatal("Not should produce a violation when inner rule passes")
	}
	if len(result.Items) != 1 {
		t.Fatalf("got %d violations, want 1", len(result.Items))
	}
	if result.Items[0].Code != "Negated" {
		t.Errorf("got code %q, want %q", result.Items[0].Code, "Negated")
	}
	if _, ok := result.Items[0].Context.(rule.Negated); !ok {
		t.Fatalf("expected Negated context, got %T", result.Items[0].Context)
	}
}

func TestComposition_Nested(t *testing.T) {
	s := rule.AllOf(
		rule.New(order.HasCustomer),
		rule.AnyOf(
			rule.New(order.HasItems),
			rule.New(order.HasPositiveTotal),
		),
	)

	v := s.Check(order{CustomerID: "c1", Total: 100})
	if !v.OK {
		t.Fatal("nested composition should be satisfied")
	}

	v = s.Check(order{CustomerID: "c1"})
	if v.OK {
		t.Fatal("nested composition should fail when AnyOf has no match")
	}
}
