# axon-rule

Composable domain specifications for Go using the Specification pattern with generics.

## Architecture

Single-package library (`rule`) with four files:

| File | Contents |
|------|----------|
| `code.go` | `Code` type (typed string), built-in codes (`MustBePresent`, `MustNotBeEmpty`, `MustBePositive`) |
| `spec.go` | `Rule[T]` interface, `New` constructor |
| `combinators.go` | `AllOf`, `AnyOf`, `Not` — all return `Rule[T]` |
| `evaluate.go` | `Evaluate` function, internal `evaluator` interface for composite recursion |
| `result.go` | `Violation` (Code + Context), `Result` ([]Violation + accessors) |

## Key design decisions

- **One method signature**: all predicates return `Verdict`. Simple rules return `Pass()` or `Fail()`, rich rules use `FailWith(typedContext)`.
- **No presentation in violations**: `Violation` has `Code` and `Context` only. Message lookup is a consumer concern.
- **Violation codes are typed constants**: domain packages own their codes, axon-rule provides common ones.
- **No `CompositeRule`**: everything is `Rule[T]`. Composites implement an unexported `evaluator` interface for `Evaluate` to recurse into.
- **Method expressions as primary pattern**: `rule.New(code, DomainType.Method)` — no closure wrappers needed.

## Testing

```bash
go test ./...
```

All tests are table-driven or single-assertion. Domain test type is `order` defined in `spec_test.go`.

## Dependencies

None. Standard library only.
