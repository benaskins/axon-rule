# axon-rule

Composable business rules for Go using the Specification pattern with generics.

## Architecture

Single-package library (`rule`) with four files:

| File | Contents |
|------|----------|
| `spec.go` | `Rule[T]` interface (single method: `Check`), `New` constructor |
| `combinators.go` | `AllOf`, `AnyOf`, `Not` — all return exported combinator types with `Evaluate` |
| `evaluate.go` | Internal `collect` function, `evaluator` interface for composite recursion |
| `result.go` | `Verdict`, `Violation`, `Violations`, `codeName` (reflect-based code derivation) |

## Key design decisions

- **Type-driven violation codes**: `Violation.Code` is derived from `reflect.TypeOf(context).Name()`. No manual string codes — the type is the identity.
- **Single-method interface**: `Rule[T]` has only `Check(T) Verdict`. No `Code()` method.
- **All failures are typed**: predicates return `FailWith(TypedContext{})`. No bare `Fail()`. Marker types (e.g. `TooFewLines struct{}`) for context-free violations.
- **No presentation in violations**: `Violation` has `Code` (derived) and `Context` (typed). Messages and translations are a consumer concern.
- **Method expressions as primary pattern**: `rule.New(DomainType.Method)` — no closure wrappers needed.
- **Exported combinator types**: `AllOfRule[T]`, `AnyOfRule[T]`, `NotRule[T]` expose `Evaluate` for violation collection.

## Testing

```bash
go test ./...
```

All tests are table-driven or single-assertion. Domain test type is `order` defined in `spec_test.go`.

## Dependencies

Standard library only (`reflect`).
