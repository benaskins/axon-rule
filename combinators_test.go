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

	r := s.IsSatisfiedBy(order{CustomerID: "c1", Items: []string{"x"}})
	if !r.OK {
		t.Fatal("AllOf should be satisfied when all specs pass")
	}
}

func TestAllOf_OneFails(t *testing.T) {
	s := spec.AllOf(
		spec.New("has-customer", order.HasCustomer),
		spec.New("has-items", order.HasItems),
	)

	r := s.IsSatisfiedBy(order{CustomerID: "c1"})
	if r.OK {
		t.Fatal("AllOf should not be satisfied when one spec fails")
	}
}

func TestAllOf_AllFail(t *testing.T) {
	s := spec.AllOf(
		spec.New("has-customer", order.HasCustomer),
		spec.New("has-items", order.HasItems),
	)

	r := s.IsSatisfiedBy(order{})
	if r.OK {
		t.Fatal("AllOf should not be satisfied when all specs fail")
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

	r := s.IsSatisfiedBy(order{CustomerID: "c1"})
	if !r.OK {
		t.Fatal("AnyOf should be satisfied when at least one spec passes")
	}
}

func TestAnyOf_NonePass(t *testing.T) {
	s := spec.AnyOf(
		spec.New("has-customer", order.HasCustomer),
		spec.New("has-items", order.HasItems),
	)

	r := s.IsSatisfiedBy(order{})
	if r.OK {
		t.Fatal("AnyOf should not be satisfied when no specs pass")
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

	r := s.IsSatisfiedBy(order{})
	if !r.OK {
		t.Fatal("Not should be satisfied when inner spec fails")
	}
}

func TestNot_Inverts_Success(t *testing.T) {
	s := spec.Not(spec.New("has-customer", order.HasCustomer))

	r := s.IsSatisfiedBy(order{CustomerID: "c1"})
	if r.OK {
		t.Fatal("Not should not be satisfied when inner spec passes")
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

	r := s.IsSatisfiedBy(order{CustomerID: "c1", Total: 100})
	if !r.OK {
		t.Fatal("nested composition should be satisfied")
	}

	r = s.IsSatisfiedBy(order{CustomerID: "c1"})
	if r.OK {
		t.Fatal("nested composition should fail when AnyOf has no match")
	}
}
