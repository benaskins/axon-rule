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

func (o order) HasCustomer() (bool, map[string]any) {
	return o.CustomerID != "", nil
}

func (o order) HasItems() (bool, map[string]any) {
	return len(o.Items) > 0, nil
}

func (o order) HasPositiveTotal() (bool, map[string]any) {
	return o.Total > 0, map[string]any{
		"total": o.Total,
	}
}

func TestNewSpec_Satisfied(t *testing.T) {
	s := spec.New("has-customer", order.HasCustomer)

	ok, ctx := s.IsSatisfiedBy(order{CustomerID: "cust-1"})
	if !ok {
		t.Fatal("spec should be satisfied")
	}
	if ctx != nil {
		t.Errorf("expected nil context, got %v", ctx)
	}
}

func TestNewSpec_NotSatisfied(t *testing.T) {
	s := spec.New("has-customer", order.HasCustomer)

	ok, _ := s.IsSatisfiedBy(order{})
	if ok {
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

	ok, ctx := s.IsSatisfiedBy(order{Total: -100})
	if ok {
		t.Fatal("spec should not be satisfied for negative total")
	}
	if ctx == nil {
		t.Fatal("expected context, got nil")
	}
	if ctx["total"] != int64(-100) {
		t.Errorf("got total %v, want -100", ctx["total"])
	}
}

func TestNewSpec_MethodExpression(t *testing.T) {
	// Verify method expressions work as the primary usage pattern.
	hasCustomer := spec.New(spec.MustBePresent, order.HasCustomer)
	hasItems := spec.New(spec.MustNotBeEmpty, order.HasItems)

	valid := order{CustomerID: "cust-1", Items: []string{"item-1"}}

	ok, _ := hasCustomer.IsSatisfiedBy(valid)
	if !ok {
		t.Fatal("hasCustomer should be satisfied")
	}

	ok, _ = hasItems.IsSatisfiedBy(valid)
	if !ok {
		t.Fatal("hasItems should be satisfied")
	}
}
