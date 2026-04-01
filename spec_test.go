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

func (o order) HasCustomer() spec.Verdict {
	if o.CustomerID != "" {
		return spec.Pass()
	}
	return spec.Fail()
}

func (o order) HasItems() spec.Verdict {
	if len(o.Items) > 0 {
		return spec.Pass()
	}
	return spec.Fail()
}

func (o order) HasPositiveTotal() spec.Verdict {
	if o.Total > 0 {
		return spec.Pass()
	}
	return spec.FailWith(map[string]any{
		"total": o.Total,
	})
}

func TestNewRule_Satisfied(t *testing.T) {
	r := spec.New("has-customer", order.HasCustomer)

	v := r.Check(order{CustomerID: "cust-1"})
	if !v.OK {
		t.Fatal("rule should be satisfied")
	}
	if v.Context != nil {
		t.Errorf("expected nil context, got %v", v.Context)
	}
}

func TestNewRule_NotSatisfied(t *testing.T) {
	r := spec.New("has-customer", order.HasCustomer)

	v := r.Check(order{})
	if v.OK {
		t.Fatal("rule should not be satisfied for empty customer")
	}
}

func TestNewRule_Code(t *testing.T) {
	r := spec.New("has-customer", order.HasCustomer)

	if r.Code() != "has-customer" {
		t.Errorf("got code %q, want %q", r.Code(), "has-customer")
	}
}

func TestNewRule_WithContext(t *testing.T) {
	r := spec.New("has-positive-total", order.HasPositiveTotal)

	v := r.Check(order{Total: -100})
	if v.OK {
		t.Fatal("rule should not be satisfied for negative total")
	}
	if v.Context == nil {
		t.Fatal("expected context, got nil")
	}
	if v.Context["total"] != int64(-100) {
		t.Errorf("got total %v, want -100", v.Context["total"])
	}
}

func TestNewRule_MethodExpression(t *testing.T) {
	hasCustomer := spec.New(spec.MustBePresent, order.HasCustomer)
	hasItems := spec.New(spec.MustNotBeEmpty, order.HasItems)

	valid := order{CustomerID: "cust-1", Items: []string{"item-1"}}

	v := hasCustomer.Check(valid)
	if !v.OK {
		t.Fatal("hasCustomer should be satisfied")
	}

	v = hasItems.Check(valid)
	if !v.OK {
		t.Fatal("hasItems should be satisfied")
	}
}
