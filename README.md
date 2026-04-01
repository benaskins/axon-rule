# axon-spec

Composable domain specifications for Go.

axon-spec implements the Specification pattern using Go generics. It provides a type-safe way to express business rules as named, testable predicates that produce structured violation reports.

## Install

```bash
go get github.com/benaskins/axon-spec
```

## Quick start

Define predicates on your domain type. Every predicate returns `(bool, map[string]any)` — the bool is the pass/fail signal, the map carries optional context about the failure.

```go
package ledger

import "time"

type JournalEntry struct {
    Description string
    Lines       []Line
    PostedAt    *time.Time
}

func (e JournalEntry) HasDescription() (bool, map[string]any) {
    return e.Description != "", nil
}

func (e JournalEntry) HasAtLeastTwoLines() (bool, map[string]any) {
    return len(e.Lines) >= 2, nil
}

func (e JournalEntry) IsNotPosted() (bool, map[string]any) {
    return e.PostedAt == nil, nil
}

func (e JournalEntry) DebitsEqualCredits() (bool, map[string]any) {
    var d, c int64
    for _, l := range e.Lines {
        d += l.Debit
        c += l.Credit
    }
    return d == c, map[string]any{
        "total_debits":  d,
        "total_credits": c,
        "difference":    d - c,
    }
}
```

Define violation codes as typed constants:

```go
package ledger

import spec "github.com/benaskins/axon-spec"

const (
    MustHaveDescription     spec.Code = "must-have-description"
    MustHaveAtLeastTwoLines spec.Code = "must-have-at-least-two-lines"
    MustNotBePosted         spec.Code = "must-not-be-posted"
    DebitsMustEqualCredits  spec.Code = "debits-must-equal-credits"
)
```

Compose specs using method expressions:

```go
package ledger

import spec "github.com/benaskins/axon-spec"

var IsValid = spec.AllOf(
    spec.New(MustHaveDescription,     JournalEntry.HasDescription),
    spec.New(MustHaveAtLeastTwoLines, JournalEntry.HasAtLeastTwoLines),
    spec.New(MustNotBePosted,         JournalEntry.IsNotPosted),
    spec.New(DebitsMustEqualCredits,  JournalEntry.DebitsEqualCredits),
)
```

Evaluate:

```go
result := spec.Evaluate(entry, ledger.IsValid)

if !result.IsValid() {
    for _, v := range result.Violations {
        fmt.Println(v.Code)    // "debits-must-equal-credits"
        fmt.Println(v.Context) // map[total_debits:5000 total_credits:3000 difference:2000]
    }
}
```

## Combinators

Combine specs to express complex eligibility rules:

```go
var CanPlaceOrder = spec.AllOf(
    spec.New(MustBeActive, Customer.IsActive),
    spec.AnyOf(
        spec.New(HasVerifiedEmail, Customer.HasVerifiedEmail),
        spec.New(HasVerifiedPhone, Customer.HasVerifiedPhone),
    ),
    spec.Not(spec.New(IsSuspended, Customer.IsSuspended)),
)
```

| Combinator | Behaviour |
|------------|-----------|
| `AllOf` | All specs must pass. Evaluates every spec, collects all violations. |
| `AnyOf` | At least one spec must pass. If none pass, collects all violations. |
| `Not` | Inverts a spec. Produces a violation with `not:` prefix on the code. |

## Event-sourced command handlers

Spec violations map directly to rejection event payloads:

```go
func (l *Ledger) Handle(cmd RecordJournalCommand) []Event {
    entry := l.buildEntry(cmd)
    result := spec.Evaluate(entry, ledger.IsValid)

    if !result.IsValid() {
        return []Event{JournalRejected{
            EntryID:    cmd.EntryID,
            Violations: result.Violations,
            RejectedAt: time.Now(),
        }}
    }

    return []Event{JournalRecorded{...}}
}
```

## Built-in codes

axon-spec provides codes for common business rules:

```go
spec.MustBePresent   // non-zero value
spec.MustNotBeEmpty  // len > 0
spec.MustBePositive  // > 0
```

These are `spec.Code` constants. Use them with your own predicates:

```go
spec.New(spec.MustBePresent, Order.HasCustomer)
spec.New(spec.MustNotBeEmpty, Order.HasLineItems)
```

## Design principles

- **One interface** — `Spec[T]` is the only abstraction. Combinators return `Spec[T]`.
- **No presentation** — `Violation` carries a code and context. Messages, translations, and resolution instructions live elsewhere.
- **Domain-owned codes** — each domain defines its violation codes as typed constants.
- **Zero dependencies** — standard library only.
