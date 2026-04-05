# axon-rule

Composable business-rule predicates and a typed state machine for Go.

`axon-rule` provides two complementary systems:

1. **Rule predicates** — typed, composable checks (`Rule[T]`, `AllOf`, `AnyOf`, `Not`)
2. **State machine** — guard-driven transitions built on top of those predicates (`Machine[T]`, `Instance[T]`)

## Prerequisites

- Go 1.24+
- [just](https://github.com/casey/just)

## Build & Test

```bash
just test
just vet
```

---

## Rule Predicates

### Core types

| Type | Description |
|------|-------------|
| `Rule[T]` | Interface: `Check(T) Verdict` |
| `Verdict` | Result of a rule check: `OK bool`, `Context any` |
| `New[T](fn)` | Create a Rule from a predicate function |
| `Pass()` | Return a passing Verdict |
| `FailWith(context)` | Return a failing Verdict with typed context |

### Combinators

```go
AllOf[T](rules ...Rule[T]) Rule[T]   // passes when all rules pass
AnyOf[T](rules ...Rule[T]) Rule[T]   // passes when at least one rule passes
Not[T](rule Rule[T]) Rule[T]         // inverts a rule
AlwaysAllow[T]() Rule[T]             // always passes (useful as a default guard)
```

### Quick example

```go
type Order struct { Total float64; Approved bool }

type NoValue struct{}
type NotApproved struct{}

var (
    hasValue  = rule.New(func(o Order) rule.Verdict {
        if o.Total > 0 { return rule.Pass() }
        return rule.FailWith(NoValue{})
    })
    isApproved = rule.New(func(o Order) rule.Verdict {
        if o.Approved { return rule.Pass() }
        return rule.FailWith(NotApproved{})
    })
    canFulfil = rule.AllOf(hasValue, isApproved)
)

v := canFulfil.Check(Order{Total: 42.0, Approved: true})
// v.OK == true
```

---

## State Machine

### Types

| Type | Description |
|------|-------------|
| `State` | Named state with optional `Metadata map[string]any`. Equality is name-based. |
| `Transition[T]` | Edge from `From State` to `To State`, guarded by a `Rule[T]`. |
| `HistoryEntry` | Records a `From`, `To`, and `Timestamp` for each applied transition. |
| `Machine[T]` | Immutable definition: initial state + transitions. |
| `Instance[T]` | Mutable running instance of a `Machine[T]`. Single-owner (not goroutine-safe). |

### Machine API

```go
// Construct
m := rule.NewMachine(initial, []rule.Transition[T]{...})

// Inspect
m.AvailableTransitions(current State, candidate T) []Transition[T]
m.TerminalStates() []State
m.Validate() error   // returns error if any state is unreachable from initial
```

### Instance API

```go
inst := rule.NewInstance(m)

inst.Current() State
inst.History() []HistoryEntry
inst.IsTerminal() bool

// Advance: apply first passing transition from current state
inst.Advance(candidate T) error

// AdvanceTo: move to a specific target state (guard must pass)
inst.AdvanceTo(target State, candidate T) error
```

---

## Worked example: factory pipeline

The factory-floor build pipeline moves a build through five stages:
`backlog → scaffolding → building → qc → lamina`

Each transition can be guarded by composed predicates.

```go
package main

import (
    "fmt"
    "github.com/benaskins/axon-rule"
)

type Build struct {
    HasSpec     bool
    TestsPassed bool
    QCPassed    bool
}

type MissingSpec struct{}
type TestsNotPassed struct{}
type QCNotPassed struct{}

func main() {
    // States
    backlog     := rule.State{Name: "backlog"}
    scaffolding := rule.State{Name: "scaffolding"}
    building    := rule.State{Name: "building"}
    qc          := rule.State{Name: "qc"}
    lamina      := rule.State{Name: "lamina"}

    // Guards
    hasSpec := rule.New(func(b Build) rule.Verdict {
        if b.HasSpec {
            return rule.Pass()
        }
        return rule.FailWith(MissingSpec{})
    })
    testsPassed := rule.New(func(b Build) rule.Verdict {
        if b.TestsPassed {
            return rule.Pass()
        }
        return rule.FailWith(TestsNotPassed{})
    })
    qcPassed := rule.New(func(b Build) rule.Verdict {
        if b.QCPassed {
            return rule.Pass()
        }
        return rule.FailWith(QCNotPassed{})
    })

    m := rule.NewMachine(backlog, []rule.Transition[Build]{
        {From: backlog,     To: scaffolding, Guard: hasSpec},
        {From: scaffolding, To: building,    Guard: rule.AlwaysAllow[Build]()},
        {From: building,    To: qc,          Guard: testsPassed},
        {From: qc,          To: lamina,      Guard: rule.AllOf(testsPassed, qcPassed)},
    })

    if err := m.Validate(); err != nil {
        panic(err)
    }

    inst := rule.NewInstance(m)
    b := Build{HasSpec: true, TestsPassed: true, QCPassed: true}

    for !inst.IsTerminal() {
        if err := inst.Advance(b); err != nil {
            fmt.Println("blocked:", err)
            break
        }
        fmt.Println("→", inst.Current().Name)
    }
    // Output:
    // → scaffolding
    // → building
    // → qc
    // → lamina
}
```

### Guard composition with combinators

```go
// Require both tests and QC before entering lamina
readyForLamina := rule.AllOf(testsPassed, qcPassed)

// Accept a build if it has a spec OR was manually approved
canScaffold := rule.AnyOf(hasSpec, manuallyApproved)

// Block a build if it is flagged
notFlagged := rule.Not(isFlagged)
```
