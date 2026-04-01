# axon-rule

Composable business rules for Go.

axon-rule implements the Specification pattern using Go generics. It provides a type-safe way to express business rules as named, testable predicates that produce structured violation reports.

## Install

```bash
go get github.com/benaskins/axon-rule
```

## Quick start

Define predicates on your domain type. Every predicate returns a `rule.Verdict` using `Pass()`, `Fail()`, or `FailWith()`.

```go
package ledger

import (
    "time"

    "github.com/benaskins/axon-rule"
)

type JournalEntry struct {
    Description string
    Lines       []Line
    PostedAt    *time.Time
}

func (e JournalEntry) HasDescription() rule.Verdict {
    if e.Description != "" {
        return rule.Pass()
    }
    return rule.Fail()
}

func (e JournalEntry) HasAtLeastTwoLines() rule.Verdict {
    if len(e.Lines) >= 2 {
        return rule.Pass()
    }
    return rule.Fail()
}

func (e JournalEntry) IsNotPosted() rule.Verdict {
    if e.PostedAt == nil {
        return rule.Pass()
    }
    return rule.Fail()
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
        TotalDebits:  d,
        TotalCredits: c,
        Difference:   d - c,
    })
}
```

Define violation codes as typed constants:

```go
package ledger

import "github.com/benaskins/axon-rule"

const (
    MustHaveDescription     rule.Code = "must-have-description"
    MustHaveAtLeastTwoLines rule.Code = "must-have-at-least-two-lines"
    MustNotBePosted         rule.Code = "must-not-be-posted"
    DebitsMustEqualCredits  rule.Code = "debits-must-equal-credits"
)
```

Compose rules using method expressions:

```go
package ledger

import "github.com/benaskins/axon-rule"

var IsValid = rule.AllOf(
    rule.New(MustHaveDescription,     JournalEntry.HasDescription),
    rule.New(MustHaveAtLeastTwoLines, JournalEntry.HasAtLeastTwoLines),
    rule.New(MustNotBePosted,         JournalEntry.IsNotPosted),
    rule.New(DebitsMustEqualCredits,  JournalEntry.DebitsEqualCredits),
)
```

Evaluate:

```go
violations := ledger.IsValid.Evaluate(entry)

if !violations.IsValid() {
    for _, v := range violations.Items {
        fmt.Println(v.Code)    // "debits-must-equal-credits"
        fmt.Println(v.Context) // BalanceMismatch{TotalDebits: 5000, ...}
    }
}
```

## Combinators

Combine rules to express complex eligibility:

```go
var CanPlaceOrder = rule.AllOf(
    rule.New(MustBeActive, Customer.IsActive),
    rule.AnyOf(
        rule.New(HasVerifiedEmail, Customer.HasVerifiedEmail),
        rule.New(HasVerifiedPhone, Customer.HasVerifiedPhone),
    ),
    rule.Not(rule.New(IsSuspended, Customer.IsSuspended)),
)
```

| Combinator | Behaviour |
|------------|-----------|
| `AllOf` | All rules must pass. Evaluates every rule, collects all violations. |
| `AnyOf` | At least one rule must pass. If none pass, collects all violations. |
| `Not` | Inverts a rule. Produces a violation with `not:` prefix on the code. |

## Event-sourced command handlers

Rule violations map directly to rejection event payloads:

```go
func (l *Ledger) Handle(cmd RecordJournalCommand) []Event {
    entry := l.buildEntry(cmd)
    violations := ledger.IsValid.Evaluate(entry)

    if !violations.IsValid() {
        return []Event{JournalRejected{
            EntryID:    cmd.EntryID,
            Violations: violations.Items,
            RejectedAt: time.Now(),
        }}
    }

    return []Event{JournalRecorded{...}}
}
```

## Built-in codes

axon-rule provides codes for common business rules:

```go
rule.MustBePresent   // non-zero value
rule.MustNotBeEmpty  // len > 0
rule.MustBePositive  // > 0
```

These are `rule.Code` constants. Use them with your own predicates:

```go
rule.New(rule.MustBePresent, Order.HasCustomer)
rule.New(rule.MustNotBeEmpty, Order.HasLineItems)
```

## Design principles

- **One interface** — `Rule[T]` is the only abstraction. Combinators return `Rule[T]`.
- **No presentation** — `Violation` carries a code and context. Messages, translations, and resolution instructions live elsewhere.
- **Domain-owned codes** — each domain defines its violation codes as typed constants.
- **Typed context** — `Verdict.Context` is `any`, so domains use their own structs for compile-time safety.
- **Zero dependencies** — standard library only.
