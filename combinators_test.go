package spec_test

import (
	"testing"

	spec "github.com/benaskins/axon-spec"
)

func TestAllOf_AllPass(t *testing.T) {
	s := spec.AllOf(
		spec.New("has-customer", order.HasCustomer),
		spec.New("has-items", order.HasItems),
	)

	v := s.Check(order{CustomerID: "c1", Items: []string{"x"}})
	if !v.OK {
		t.Fatal("AllOf should be satisfied when all rules pass")
	}
}

func TestAllOf_OneFails(t *testing.T) {
	s := spec.AllOf(
		spec.New("has-customer", order.HasCustomer),
		spec.New("has-items", order.HasItems),
	)

	v := s.Check(order{CustomerID: "c1"})
	if v.OK {
		t.Fatal("AllOf should not be satisfied when one rule fails")
	}
}

func TestAllOf_AllFail(t *testing.T) {
	s := spec.AllOf(
		spec.New("has-customer", order.HasCustomer),
		spec.New("has-items", order.HasItems),
	)

	v := s.Check(order{})
	if v.OK {
		t.Fatal("AllOf should not be satisfied when all rules fail")
	}
}

func TestAllOf_Code(t *testing.T) {
	s := spec.AllOf(
		spec.New("has-customer", order.HasCustomer),
	)

	if s.Code() != "all-of" {
		t.Errorf("got code %q, want %q", s.Code(), "all-of")
	}
}

func TestAnyOf_OnePass(t *testing.T) {
	s := spec.AnyOf(
		spec.New("has-customer", order.HasCustomer),
		spec.New("has-items", order.HasItems),
	)

	v := s.Check(order{CustomerID: "c1"})
	if !v.OK {
		t.Fatal("AnyOf should be satisfied when at least one rule passes")
	}
}

func TestAnyOf_NonePass(t *testing.T) {
	s := spec.AnyOf(
		spec.New("has-customer", order.HasCustomer),
		spec.New("has-items", order.HasItems),
	)

	v := s.Check(order{})
	if v.OK {
		t.Fatal("AnyOf should not be satisfied when no rules pass")
	}
}

func TestAnyOf_Code(t *testing.T) {
	s := spec.AnyOf(
		spec.New("has-customer", order.HasCustomer),
	)

	if s.Code() != "any-of" {
		t.Errorf("got code %q, want %q", s.Code(), "any-of")
	}
}

func TestNot_Inverts_Failure(t *testing.T) {
	s := spec.Not(spec.New("has-customer", order.HasCustomer))

	v := s.Check(order{})
	if !v.OK {
		t.Fatal("Not should be satisfied when inner rule fails")
	}
}

func TestNot_Inverts_Success(t *testing.T) {
	s := spec.Not(spec.New("has-customer", order.HasCustomer))

	v := s.Check(order{CustomerID: "c1"})
	if v.OK {
		t.Fatal("Not should not be satisfied when inner rule passes")
	}
}

func TestNot_Code(t *testing.T) {
	s := spec.Not(spec.New("suspended", order.HasCustomer))

	if s.Code() != "not:suspended" {
		t.Errorf("got code %q, want %q", s.Code(), "not:suspended")
	}
}

func TestComposition_Nested(t *testing.T) {
	s := spec.AllOf(
		spec.New("has-customer", order.HasCustomer),
		spec.AnyOf(
			spec.New("has-items", order.HasItems),
			spec.New("has-positive-total", order.HasPositiveTotal),
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
