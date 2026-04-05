# axon-rule

## Build & Test

```bash
go test ./...
go vet ./...
just build     # builds to bin/axon-rule
just install   # copies to ~/.local/bin/axon-rule
```

## Module Selections

- **axon-rule**: Composable business-rule predicates (`Rule[T]`, `AllOf`, `AnyOf`, `Not`, `Verdict`) plus a guard-driven state machine (`State`, `Transition[T]`, `Machine[T]`, `Instance[T]`, `HistoryEntry`, `AlwaysAllow`). All state machine types are additive — no existing types were changed. (deterministic)

## Architecture

### Rule predicates

`Rule[T]` is a single-method interface (`Check(T) Verdict`). Rules are created with `New[T]` and composed with `AllOf`, `AnyOf`, and `Not`. `Verdict` carries an `OK` flag and a `Context []string` of violation messages.

### State machine

Built entirely on top of the rule predicate system — no new dependencies.

| Type | Role |
|------|------|
| `State` | Named vertex. Equality is name-based; `Metadata` is informational only. |
| `Transition[T]` | Directed edge `From → To` with a `Guard Rule[T]`. |
| `Machine[T]` | Immutable definition: initial state + transitions slice (copied at construction). |
| `Instance[T]` | Mutable running instance. Tracks `current State` and `[]HistoryEntry`. Single-owner, not goroutine-safe. |
| `HistoryEntry` | Snapshot of a completed transition: `From`, `To`, `Timestamp`. |
| `AlwaysAllow[T]()` | Convenience guard that always returns a passing `Verdict`. |

### Key methods

```
Machine[T].AvailableTransitions(current, candidate) — guard-filtered outgoing transitions
Machine[T].TerminalStates()                         — states with no outgoing transitions
Machine[T].Validate()                               — BFS reachability check from initial state
Instance[T].Advance(candidate)                      — apply first available transition
Instance[T].AdvanceTo(target, candidate)            — move to a named target (guard must pass)
Instance[T].IsTerminal()                            — true when current state is terminal
```

## Deterministic / Non-deterministic Boundary

| From | To | Type |
|------|----|------|
| `Machine[T]` | `Rule[T]` (guard evaluation) | det |
| `Instance[T].Advance` | `Machine[T].AvailableTransitions` | det |
| `Instance[T].AdvanceTo` | `Machine[T]` (transition lookup + guard check) | det |

## Dependency Graph

```
axon-rule (stdlib only)
  └── no external dependencies
```
