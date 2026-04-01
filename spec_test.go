package spec_test

import (
	"testing"

	spec "github.com/benaskins/axon-spec"
)

type order struct {
	CustomerID string
	Items      []string
	Total      int64
}

func (o order) HasCustomer() spec.PredicateResult {
	if o.CustomerID != "" {
		return spec.Pass()
	}
	return spec.Fail()
}

func (o order) HasItems() spec.PredicateResult {
	if len(o.Items) > 0 {
		return spec.Pass()
	}
	return spec.Fail()
}

func (o order) HasPositiveTotal() spec.PredicateResult {
	if o.Total > 0 {
		return spec.Pass()
	}
	return spec.FailWith(map[string]any{
		"total": o.Total,
	})
}

func TestNewSpec_Satisfied(t *testing.T) {
	s := spec.New("has-customer", order.HasCustomer)

	r := s.IsSatisfiedBy(order{CustomerID: "cust-1"})
	if !r.OK {
		t.Fatal("spec should be satisfied")
	}
	if r.Context != nil {
		t.Errorf("expected nil context, got %v", r.Context)
	}
}

func TestNewSpec_NotSatisfied(t *testing.T) {
	s := spec.New("has-customer", order.HasCustomer)

	r := s.IsSatisfiedBy(order{})
	if r.OK {
		t.Fatal("spec should not be satisfied for empty customer")
	}
}

func TestNewSpec_Code(t *testing.T) {
	s := spec.New("has-customer", order.HasCustomer)

	if s.Code() != "has-customer" {
		t.Errorf("got code %q, want %q", s.Code(), "has-customer")
	}
}

func TestNewSpec_WithContext(t *testing.T) {
	s := spec.New("has-positive-total", order.HasPositiveTotal)

	r := s.IsSatisfiedBy(order{Total: -100})
	if r.OK {
		t.Fatal("spec should not be satisfied for negative total")
	}
	if r.Context == nil {
		t.Fatal("expected context, got nil")
	}
	if r.Context["total"] != int64(-100) {
		t.Errorf("got total %v, want -100", r.Context["total"])
	}
}

func TestNewSpec_MethodExpression(t *testing.T) {
	hasCustomer := spec.New(spec.MustBePresent, order.HasCustomer)
	hasItems := spec.New(spec.MustNotBeEmpty, order.HasItems)

	valid := order{CustomerID: "cust-1", Items: []string{"item-1"}}

	r := hasCustomer.IsSatisfiedBy(valid)
	if !r.OK {
		t.Fatal("hasCustomer should be satisfied")
	}

	r = hasItems.IsSatisfiedBy(valid)
	if !r.OK {
		t.Fatal("hasItems should be satisfied")
	}
}
