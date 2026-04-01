# axon-rule

Composable business rules for Go.

axon-rule implements the Specification pattern using Go generics. Rules produce typed violations — the type name becomes the violation code, giving you compile-time safety and stable keys for localisation.

## Install

```bash
go get github.com/benaskins/axon-rule
```

## Quick start

Define violation types and predicates on your domain type:

```go
package ledger

import (
    "github.com/benaskins/axon-rule"
)

type TooFewLines struct{}
type MissingDescription struct{}

type BalanceMismatch struct {
    TotalDebits  string
    TotalCredits string
}

func (e JournalEntry) HasDescription() rule.Verdict {
    if e.Description != "" {
        return rule.Pass()
    }
    return rule.FailWith(MissingDescription{})
}

func (e JournalEntry) HasAtLeastTwoLines() rule.Verdict {
    if len(e.Lines) >= 2 {
        return rule.Pass()
    }
    return rule.FailWith(TooFewLines{})
}

func (e JournalEntry) DebitsEqualCredits() rule.Verdict {
    var d, c int64
    for _, l := range e.Lines {
        d += l.Debit
        c += l.Credit
    }
    if d == c {
        return rule.Pass()
    }
    return rule.FailWith(BalanceMismatch{
        TotalDebits:  fmt.Sprint(d),
        TotalCredits: fmt.Sprint(c),
    })
}
```

Compose rules using method expressions:

```go
var IsValid = rule.AllOf(
    rule.New(JournalEntry.HasDescription),
    rule.New(JournalEntry.HasAtLeastTwoLines),
    rule.New(JournalEntry.DebitsEqualCredits),
)
```

Evaluate and consume:

```go
violations := ledger.IsValid.Evaluate(entry)

if !violations.IsValid() {
    for _, v := range violations.Items {
        fmt.Println(v.Code) // "BalanceMismatch" — derived from type name
    }
}
```

Match on types for compile-time safety:

```go
for _, v := range violations.Items {
    switch ctx := v.Context.(type) {
    case ledger.BalanceMismatch:
        fmt.Printf("debits=%s credits=%s\n", ctx.TotalDebits, ctx.TotalCredits)
    case ledger.TooFewLines:
        fmt.Println("need at least two lines")
    }
}
```

Use the code string for i18n lookups:

```go
for _, v := range violations.Items {
    msg := localiser.Lookup(v.Code) // "BalanceMismatch" → localised string
    fmt.Println(msg)
}
```

## Combinators

Combine rules to express complex eligibility:

```go
var CanPlaceOrder = rule.AllOf(
    rule.New(Customer.IsActive),
    rule.AnyOf(
        rule.New(Customer.HasVerifiedEmail),
        rule.New(Customer.HasVerifiedPhone),
    ),
    rule.Not(rule.New(Customer.IsSuspended)),
)
```

| Combinator | Behaviour |
|------------|-----------|
| `AllOf` | All rules must pass. Evaluates every rule, collects all violations. |
| `AnyOf` | At least one rule must pass. If none pass, collects all violations. |
| `Not` | Inverts a rule. Produces a `Negated` violation when the inner rule passes. |

## Design principles

- **Type is identity** — violation codes are derived from `reflect.TypeOf(context).Name()`. No manual string constants.
- **One interface** — `Rule[T]` has a single method: `Check(T) Verdict`.
- **No presentation** — `Violation` carries a code and context. Messages and translations live elsewhere.
- **Typed context** — consumers use type switches for compile-time safety. Marker types (`struct{}`) for context-free violations.
- **Zero dependencies** — standard library only.
