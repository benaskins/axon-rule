package rule_test

import (
	"testing"

	"github.com/benaskins/axon-rule"
)

type order struct {
	CustomerID string
	Items      []string
	Total      int64
}

func (o order) HasCustomer() rule.Verdict {
	if o.CustomerID != "" {
		return rule.Pass()
	}
	return rule.Fail()
}

func (o order) HasItems() rule.Verdict {
	if len(o.Items) > 0 {
		return rule.Pass()
	}
	return rule.Fail()
}

func (o order) HasPositiveTotal() rule.Verdict {
	if o.Total > 0 {
		return rule.Pass()
	}
	return rule.FailWith(map[string]any{
		"total": o.Total,
	})
}

func TestNewRule_Satisfied(t *testing.T) {
	r := rule.New("has-customer", order.HasCustomer)

	v := r.Check(order{CustomerID: "cust-1"})
	if !v.OK {
		t.Fatal("rule should be satisfied")
	}
	if v.Context != nil {
		t.Errorf("expected nil context, got %v", v.Context)
	}
}

func TestNewRule_NotSatisfied(t *testing.T) {
	r := rule.New("has-customer", order.HasCustomer)

	v := r.Check(order{})
	if v.OK {
		t.Fatal("rule should not be satisfied for empty customer")
	}
}

func TestNewRule_Code(t *testing.T) {
	r := rule.New("has-customer", order.HasCustomer)

	if r.Code() != "has-customer" {
		t.Errorf("got code %q, want %q", r.Code(), "has-customer")
	}
}

func TestNewRule_WithContext(t *testing.T) {
	r := rule.New("has-positive-total", order.HasPositiveTotal)

	v := r.Check(order{Total: -100})
	if v.OK {
		t.Fatal("rule should not be satisfied for negative total")
	}
	if v.Context == nil {
		t.Fatal("expected context, got nil")
	}
	ctx, ok := v.Context.(map[string]any)
	if !ok {
		t.Fatalf("expected map[string]any context, got %T", v.Context)
	}
	if ctx["total"] != int64(-100) {
		t.Errorf("got total %v, want -100", ctx["total"])
	}
}

func TestNewRule_MethodExpression(t *testing.T) {
	hasCustomer := rule.New(rule.MustBePresent, order.HasCustomer)
	hasItems := rule.New(rule.MustNotBeEmpty, order.HasItems)

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
