# axon-rule

Composable business rules using the Specification pattern with generics, plus a guard-driven state machine.

## Module

- Module path: `github.com/benaskins/axon-rule`
- Project type: library (no main package)

## Build & Test

```bash
just test    # go test -race ./...
just vet     # go vet ./...
```

## Architecture

Single-package library (`rule`) with two systems:

**Rule predicates:** `Rule[T]` interface (single method: `Check`), `New` constructor, `AllOf`/`AnyOf`/`Not` combinators, `Verdict`/`Violation`/`Violations` results. Type-driven violation codes via `reflect.TypeOf(context).Name()`.

**State machine:** `State` (name-based equality), `Transition[T]` (guarded edges), `Machine[T]` (immutable definition with validation), `Instance[T]` (mutable runtime, single-owner). Guards are `Rule[T]` — full combinator composition.

Read [AGENTS.md](./AGENTS.md) for architecture details.

## Constraints

- Standard library only — zero external dependencies
- All existing types must be preserved — state machine additions are strictly additive
- `Instance[T]` is single-owner — no mutexes or sync primitives
- No third-party assertion libraries — standard `testing` package only
