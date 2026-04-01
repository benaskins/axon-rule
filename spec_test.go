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

// Violation context types
type MissingCustomer struct{}
type EmptyItems struct{}
type NegativeTotal struct {
	Total int64
}

func (o order) HasCustomer() rule.Verdict {
	if o.CustomerID != "" {
		return rule.Pass()
	}
	return rule.FailWith(MissingCustomer{})
}

func (o order) HasItems() rule.Verdict {
	if len(o.Items) > 0 {
		return rule.Pass()
	}
	return rule.FailWith(EmptyItems{})
}

func (o order) HasPositiveTotal() rule.Verdict {
	if o.Total > 0 {
		return rule.Pass()
	}
	return rule.FailWith(NegativeTotal{Total: o.Total})
}

func TestNewRule_Satisfied(t *testing.T) {
	r := rule.New(order.HasCustomer)

	v := r.Check(order{CustomerID: "cust-1"})
	if !v.OK {
		t.Fatal("rule should be satisfied")
	}
}

func TestNewRule_NotSatisfied(t *testing.T) {
	r := rule.New(order.HasCustomer)

	v := r.Check(order{})
	if v.OK {
		t.Fatal("rule should not be satisfied for empty customer")
	}
}

func TestNewRule_ContextType(t *testing.T) {
	r := rule.New(order.HasPositiveTotal)

	v := r.Check(order{Total: -100})
	if v.OK {
		t.Fatal("rule should not be satisfied for negative total")
	}
	if v.Context == nil {
		t.Fatal("expected context, got nil")
	}
	ctx, ok := v.Context.(NegativeTotal)
	if !ok {
		t.Fatalf("expected NegativeTotal context, got %T", v.Context)
	}
	if ctx.Total != int64(-100) {
		t.Errorf("got total %v, want -100", ctx.Total)
	}
}

func TestNewRule_MethodExpression(t *testing.T) {
	hasCustomer := rule.New(order.HasCustomer)
	hasItems := rule.New(order.HasItems)

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
