# axon-rule — Initial Build Plan
# 2026-04-03

Each step is commit-sized. Execute via `/iterate`.

## Step 1 — Define State and Transition types

Add `State` type (name string, metadata map[string]any) with name-based equality, and `Transition[T]` type (From State, To State, Guard Rule[T]) to a new file `statemachine.go` in the axon-rule package. Add `AlwaysAllow[T]() Rule[T]` convenience constructor that returns a Rule whose Check always returns a passing Verdict. Write unit tests in `statemachine_test.go` verifying State equality by name, Transition field access, and AlwaysAllow always passes.

Commit: `feat: add State, Transition, and AlwaysAllow types to axon-rule`

## Step 2 — Implement Machine[T] definition type

Add `Machine[T]` struct storing the initial State and transitions slice (immutable after construction). Implement: `NewMachine[T](initial State, transitions []Transition[T]) *Machine[T]`, `AvailableTransitions(current State, candidate T) []Transition[T]` (filters by From state then evaluates each guard, collecting those whose Verdict passes), `TerminalStates() []State` (states with no outgoing transitions, derived from all mentioned states), and `Validate() error` (detects unreachable states — states not reachable from initial — and non-terminal dead-ends with zero outgoing transitions). Test cases: machine with linear A→B→C, branching guards, validate catches unreachable state, terminal states identified correctly.

Commit: `feat: add Machine[T] with constructor, AvailableTransitions, TerminalStates, and Validate`

## Step 3 — Implement Instance[T] running instance type

Add `HistoryEntry` struct (From State, To State, Timestamp time.Time) and `Instance[T]` struct (machine *Machine[T], current State, history []HistoryEntry). Implement: `NewInstance[T](machine *Machine[T]) *Instance[T]` (starts at machine's initial state, empty history), `Current() State`, `History() []HistoryEntry` (returns a copy to prevent external mutation), `Advance(candidate T) error` (calls AvailableTransitions, applies first passing transition, appends HistoryEntry with time.Now(), returns error if no transition available), `AdvanceTo(target State, candidate T) error` (finds the specific From=current To=target transition, checks guard, returns error if not found or guard fails, applies and records if passing), `IsTerminal() bool` (delegates to Machine.TerminalStates). Test: full lifecycle advance through A→B→C, AdvanceTo success and failure paths, error on no available transition, history length and entries, IsTerminal at terminal state.

Commit: `feat: add Instance[T] with Advance, AdvanceTo, Current, History, and IsTerminal`

## Step 4 — Add combinator-guard integration tests

Write integration-style tests (still pure unit tests, no I/O) that demonstrate state machine guards composed from existing axon-rule combinators: AllOf, AnyOf, Not. Scenarios: a pipeline machine (backlog→scaffolding→building→qc→lamina) using AllOf guards on a struct candidate, a branching machine where AnyOf allows multiple conditions to permit a transition, a Not combinator blocking a transition when a negative condition holds. Verify AvailableTransitions returns correct subsets and Instance.Advance follows guard logic. All tests in existing package, no external deps, no writes outside t.TempDir().

Commit: `feat: add integration tests composing state machine guards with AllOf, AnyOf, and Not combinators`

## Step 5 — Document state machine API in README and AGENTS.md

Update axon-rule README.md to document the new state machine API: State, Transition, Machine, Instance, HistoryEntry, AlwaysAllow. Add a worked example showing the factory pipeline (backlog→scaffolding→building→qc→lamina) with guard composition. Update AGENTS.md module entry for axon-rule to reference the state machine enhancement. No code changes — documentation only.

Commit: `docs: update README and AGENTS.md with state machine API and usage examples`

