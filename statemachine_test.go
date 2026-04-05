package rule

import (
	"testing"
)

// --- Machine[T] tests ---

func TestNewMachineLinear(t *testing.T) {
	stateA := State{Name: "A"}
	stateB := State{Name: "B"}
	stateC := State{Name: "C"}

	m := NewMachine[string](stateA, []Transition[string]{
		{From: stateA, To: stateB, Guard: AlwaysAllow[string]()},
		{From: stateB, To: stateC, Guard: AlwaysAllow[string]()},
	})

	avail := m.AvailableTransitions(stateA, "x")
	if len(avail) != 1 {
		t.Fatalf("expected 1 available transition from A, got %d", len(avail))
	}
	if !avail[0].To.Equal(stateB) {
		t.Errorf("expected transition to B, got %v", avail[0].To)
	}

	availB := m.AvailableTransitions(stateB, "x")
	if len(availB) != 1 || !availB[0].To.Equal(stateC) {
		t.Errorf("expected 1 transition B→C, got %v", availB)
	}

	// C has no outgoing transitions
	availC := m.AvailableTransitions(stateC, "x")
	if len(availC) != 0 {
		t.Errorf("expected 0 transitions from C, got %d", len(availC))
	}
}

func TestMachineTerminalStates(t *testing.T) {
	stateA := State{Name: "A"}
	stateB := State{Name: "B"}
	stateC := State{Name: "C"}

	m := NewMachine[string](stateA, []Transition[string]{
		{From: stateA, To: stateB, Guard: AlwaysAllow[string]()},
		{From: stateB, To: stateC, Guard: AlwaysAllow[string]()},
	})

	terminals := m.TerminalStates()
	if len(terminals) != 1 {
		t.Fatalf("expected 1 terminal state, got %d", len(terminals))
	}
	if !terminals[0].Equal(stateC) {
		t.Errorf("expected terminal state C, got %v", terminals[0])
	}
}

func TestMachineBranchingGuards(t *testing.T) {
	stateA := State{Name: "A"}
	stateB := State{Name: "B"}
	stateC := State{Name: "C"}

	allowB := New(func(s string) Verdict {
		if s == "go-b" {
			return Pass()
		}
		return FailWith("blocked")
	})
	allowC := New(func(s string) Verdict {
		if s == "go-c" {
			return Pass()
		}
		return FailWith("blocked")
	})

	m := NewMachine[string](stateA, []Transition[string]{
		{From: stateA, To: stateB, Guard: allowB},
		{From: stateA, To: stateC, Guard: allowC},
	})

	// candidate "go-b" only opens A→B
	avail := m.AvailableTransitions(stateA, "go-b")
	if len(avail) != 1 || !avail[0].To.Equal(stateB) {
		t.Errorf("expected only A→B for 'go-b', got %v", avail)
	}

	// candidate "go-c" only opens A→C
	avail2 := m.AvailableTransitions(stateA, "go-c")
	if len(avail2) != 1 || !avail2[0].To.Equal(stateC) {
		t.Errorf("expected only A→C for 'go-c', got %v", avail2)
	}

	// unmatched candidate opens nothing
	avail3 := m.AvailableTransitions(stateA, "nope")
	if len(avail3) != 0 {
		t.Errorf("expected 0 transitions for 'nope', got %d", len(avail3))
	}
}

func TestMachineValidateUnreachable(t *testing.T) {
	stateA := State{Name: "A"}
	stateB := State{Name: "B"}
	stateOrphan := State{Name: "Orphan"}

	// Orphan→B is defined but Orphan is not reachable from initial A
	m := NewMachine[string](stateA, []Transition[string]{
		{From: stateA, To: stateB, Guard: AlwaysAllow[string]()},
		{From: stateOrphan, To: stateB, Guard: AlwaysAllow[string]()},
	})

	err := m.Validate()
	if err == nil {
		t.Error("expected Validate to return an error for unreachable state Orphan")
	}
}

func TestMachineValidateReachable(t *testing.T) {
	stateA := State{Name: "A"}
	stateB := State{Name: "B"}
	stateC := State{Name: "C"}

	m := NewMachine[string](stateA, []Transition[string]{
		{From: stateA, To: stateB, Guard: AlwaysAllow[string]()},
		{From: stateB, To: stateC, Guard: AlwaysAllow[string]()},
	})

	if err := m.Validate(); err != nil {
		t.Errorf("expected valid machine, got error: %v", err)
	}
}

// --- Instance[T] tests ---

func TestInstanceLifecycle(t *testing.T) {
	stateA := State{Name: "A"}
	stateB := State{Name: "B"}
	stateC := State{Name: "C"}

	m := NewMachine[string](stateA, []Transition[string]{
		{From: stateA, To: stateB, Guard: AlwaysAllow[string]()},
		{From: stateB, To: stateC, Guard: AlwaysAllow[string]()},
	})

	inst := NewInstance(m)
	if !inst.Current().Equal(stateA) {
		t.Fatalf("expected initial state A, got %v", inst.Current())
	}
	if len(inst.History()) != 0 {
		t.Fatalf("expected empty history, got %d entries", len(inst.History()))
	}

	if err := inst.Advance("x"); err != nil {
		t.Fatalf("unexpected error advancing A→B: %v", err)
	}
	if !inst.Current().Equal(stateB) {
		t.Errorf("expected state B after advance, got %v", inst.Current())
	}

	if err := inst.Advance("x"); err != nil {
		t.Fatalf("unexpected error advancing B→C: %v", err)
	}
	if !inst.Current().Equal(stateC) {
		t.Errorf("expected state C after advance, got %v", inst.Current())
	}

	h := inst.History()
	if len(h) != 2 {
		t.Fatalf("expected 2 history entries, got %d", len(h))
	}
	if !h[0].From.Equal(stateA) || !h[0].To.Equal(stateB) {
		t.Errorf("expected history[0] A→B, got %v→%v", h[0].From, h[0].To)
	}
	if !h[1].From.Equal(stateB) || !h[1].To.Equal(stateC) {
		t.Errorf("expected history[1] B→C, got %v→%v", h[1].From, h[1].To)
	}
	if h[0].Timestamp.IsZero() {
		t.Error("expected non-zero timestamp in history entry")
	}
}

func TestInstanceAdvanceNoTransition(t *testing.T) {
	stateA := State{Name: "A"}
	stateB := State{Name: "B"}

	blocked := New(func(_ string) Verdict { return FailWith("blocked") })
	m := NewMachine[string](stateA, []Transition[string]{
		{From: stateA, To: stateB, Guard: blocked},
	})

	inst := NewInstance(m)
	err := inst.Advance("x")
	if err == nil {
		t.Error("expected error when no transition is available")
	}
}

func TestInstanceHistoryIsCopy(t *testing.T) {
	stateA := State{Name: "A"}
	stateB := State{Name: "B"}

	m := NewMachine[string](stateA, []Transition[string]{
		{From: stateA, To: stateB, Guard: AlwaysAllow[string]()},
	})

	inst := NewInstance(m)
	_ = inst.Advance("x")

	h1 := inst.History()
	h1[0].From = State{Name: "mutated"}

	h2 := inst.History()
	if h2[0].From.Name == "mutated" {
		t.Error("History() must return a copy; external mutation affected internal state")
	}
}

func TestInstanceAdvanceTo(t *testing.T) {
	stateA := State{Name: "A"}
	stateB := State{Name: "B"}
	stateC := State{Name: "C"}

	m := NewMachine[string](stateA, []Transition[string]{
		{From: stateA, To: stateB, Guard: AlwaysAllow[string]()},
		{From: stateA, To: stateC, Guard: AlwaysAllow[string]()},
	})

	inst := NewInstance(m)
	if err := inst.AdvanceTo(stateC, "x"); err != nil {
		t.Fatalf("unexpected error in AdvanceTo C: %v", err)
	}
	if !inst.Current().Equal(stateC) {
		t.Errorf("expected state C, got %v", inst.Current())
	}
}

func TestInstanceAdvanceToNotFound(t *testing.T) {
	stateA := State{Name: "A"}
	stateB := State{Name: "B"}
	stateC := State{Name: "C"}

	m := NewMachine[string](stateA, []Transition[string]{
		{From: stateA, To: stateB, Guard: AlwaysAllow[string]()},
	})

	inst := NewInstance(m)
	err := inst.AdvanceTo(stateC, "x")
	if err == nil {
		t.Error("expected error when target transition not found")
	}
}

func TestInstanceAdvanceToGuardFails(t *testing.T) {
	stateA := State{Name: "A"}
	stateB := State{Name: "B"}

	blocked := New(func(_ string) Verdict { return FailWith("blocked") })
	m := NewMachine[string](stateA, []Transition[string]{
		{From: stateA, To: stateB, Guard: blocked},
	})

	inst := NewInstance(m)
	err := inst.AdvanceTo(stateB, "x")
	if err == nil {
		t.Error("expected error when guard fails in AdvanceTo")
	}
	if !inst.Current().Equal(stateA) {
		t.Error("state should not change when AdvanceTo fails")
	}
}

func TestInstanceIsTerminal(t *testing.T) {
	stateA := State{Name: "A"}
	stateB := State{Name: "B"}

	m := NewMachine[string](stateA, []Transition[string]{
		{From: stateA, To: stateB, Guard: AlwaysAllow[string]()},
	})

	inst := NewInstance(m)
	if inst.IsTerminal() {
		t.Error("A should not be terminal")
	}

	_ = inst.Advance("x")
	if !inst.IsTerminal() {
		t.Error("B should be terminal")
	}
}

// --- Step 1 tests (unchanged) ---

func TestStateEqualityByName(t *testing.T) {
	a := State{Name: "active", Metadata: map[string]any{"foo": 1}}
	b := State{Name: "active", Metadata: map[string]any{"bar": 2}}
	c := State{Name: "inactive"}

	if !a.Equal(b) {
		t.Error("states with same name should be equal regardless of metadata")
	}
	if a.Equal(c) {
		t.Error("states with different names should not be equal")
	}
}

func TestTransitionFieldAccess(t *testing.T) {
	from := State{Name: "pending"}
	to := State{Name: "approved"}
	guard := AlwaysAllow[string]()

	tr := Transition[string]{From: from, To: to, Guard: guard}

	if !tr.From.Equal(from) {
		t.Errorf("expected From %v, got %v", from, tr.From)
	}
	if !tr.To.Equal(to) {
		t.Errorf("expected To %v, got %v", to, tr.To)
	}
	if tr.Guard == nil {
		t.Error("expected Guard to be set")
	}
}

func TestAlwaysAllowAlwaysPasses(t *testing.T) {
	rule := AlwaysAllow[string]()

	v := rule.Check("anything")
	if !v.OK {
		t.Error("AlwaysAllow should return a passing verdict")
	}

	v2 := rule.Check("")
	if !v2.OK {
		t.Error("AlwaysAllow should return a passing verdict for empty string")
	}
}

// --- Combinator integration tests ---

// artifact is a build artifact candidate used by combinator integration tests.
type artifact struct {
	HasCode     bool
	TestsPassed bool
	Reviewed    bool
	Approved    bool
	IsBlocked   bool
}

func artifactRule(name string, fn func(artifact) bool) Rule[artifact] {
	return New(func(a artifact) Verdict {
		if fn(a) {
			return Pass()
		}
		return FailWith(name)
	})
}

// TestAllOfGuardPipeline tests a 5-stage pipeline where each transition is
// guarded by AllOf with multiple conditions on a struct candidate.
func TestAllOfGuardPipeline(t *testing.T) {
	backlog := State{Name: "backlog"}
	scaffolding := State{Name: "scaffolding"}
	building := State{Name: "building"}
	qc := State{Name: "qc"}
	lamina := State{Name: "lamina"}

	hasCode := artifactRule("has_code", func(a artifact) bool { return a.HasCode })
	testsPassed := artifactRule("tests_passed", func(a artifact) bool { return a.TestsPassed })
	reviewed := artifactRule("reviewed", func(a artifact) bool { return a.Reviewed })
	approved := artifactRule("approved", func(a artifact) bool { return a.Approved })

	m := NewMachine[artifact](backlog, []Transition[artifact]{
		{From: backlog, To: scaffolding, Guard: AllOf[artifact](hasCode)},
		{From: scaffolding, To: building, Guard: AllOf[artifact](hasCode, testsPassed)},
		{From: building, To: qc, Guard: AllOf[artifact](testsPassed, reviewed)},
		{From: qc, To: lamina, Guard: AllOf[artifact](reviewed, approved)},
	})

	if err := m.Validate(); err != nil {
		t.Fatalf("pipeline machine should be valid: %v", err)
	}

	// No conditions met: backlog has no available transitions.
	empty := artifact{}
	if avail := m.AvailableTransitions(backlog, empty); len(avail) != 0 {
		t.Errorf("expected 0 transitions from backlog with no conditions, got %d", len(avail))
	}

	// HasCode=true: backlog→scaffolding unlocks.
	withCode := artifact{HasCode: true}
	avail := m.AvailableTransitions(backlog, withCode)
	if len(avail) != 1 || !avail[0].To.Equal(scaffolding) {
		t.Errorf("expected backlog→scaffolding with HasCode, got %v", avail)
	}

	// TestsPassed alone is not enough for scaffolding→building (also needs HasCode).
	testsOnly := artifact{TestsPassed: true}
	if avail := m.AvailableTransitions(scaffolding, testsOnly); len(avail) != 0 {
		t.Errorf("expected 0 transitions from scaffolding without HasCode, got %d", len(avail))
	}

	// Drive an instance through the full pipeline.
	inst := NewInstance(m)

	steps := []struct {
		cand artifact
		want State
	}{
		{artifact{HasCode: true}, scaffolding},
		{artifact{HasCode: true, TestsPassed: true}, building},
		{artifact{TestsPassed: true, Reviewed: true}, qc},
		{artifact{Reviewed: true, Approved: true}, lamina},
	}

	for _, step := range steps {
		if err := inst.Advance(step.cand); err != nil {
			t.Fatalf("Advance to %q failed: %v", step.want.Name, err)
		}
		if !inst.Current().Equal(step.want) {
			t.Errorf("expected state %q, got %q", step.want.Name, inst.Current().Name)
		}
	}

	if !inst.IsTerminal() {
		t.Error("lamina should be a terminal state")
	}
	if len(inst.History()) != 4 {
		t.Errorf("expected 4 history entries, got %d", len(inst.History()))
	}
}

// TestAnyOfGuardBranching tests a machine where AnyOf allows a transition when
// at least one of several conditions is satisfied.
func TestAnyOfGuardBranching(t *testing.T) {
	pending := State{Name: "pending"}
	accepted := State{Name: "accepted"}
	rejected := State{Name: "rejected"}

	reviewed := artifactRule("reviewed", func(a artifact) bool { return a.Reviewed })
	approved := artifactRule("approved", func(a artifact) bool { return a.Approved })
	blocked := artifactRule("blocked", func(a artifact) bool { return a.IsBlocked })

	m := NewMachine[artifact](pending, []Transition[artifact]{
		// Either reviewed or approved is sufficient to accept.
		{From: pending, To: accepted, Guard: AnyOf[artifact](reviewed, approved)},
		// Only blocked items go to rejected.
		{From: pending, To: rejected, Guard: AnyOf[artifact](blocked)},
	})

	// Neither condition: no transitions available.
	neither := artifact{}
	if avail := m.AvailableTransitions(pending, neither); len(avail) != 0 {
		t.Errorf("expected 0 transitions with no conditions, got %d", len(avail))
	}

	// Reviewed only: pending→accepted available.
	reviewedOnly := artifact{Reviewed: true}
	avail := m.AvailableTransitions(pending, reviewedOnly)
	if len(avail) != 1 || !avail[0].To.Equal(accepted) {
		t.Errorf("expected pending→accepted for Reviewed, got %v", avail)
	}

	// Approved only: pending→accepted available.
	approvedOnly := artifact{Approved: true}
	avail2 := m.AvailableTransitions(pending, approvedOnly)
	if len(avail2) != 1 || !avail2[0].To.Equal(accepted) {
		t.Errorf("expected pending→accepted for Approved, got %v", avail2)
	}

	// Blocked: pending→rejected available.
	blockedArtifact := artifact{IsBlocked: true}
	avail3 := m.AvailableTransitions(pending, blockedArtifact)
	if len(avail3) != 1 || !avail3[0].To.Equal(rejected) {
		t.Errorf("expected pending→rejected for IsBlocked, got %v", avail3)
	}

	// Reviewed AND blocked: both transitions available.
	reviewedAndBlocked := artifact{Reviewed: true, IsBlocked: true}
	avail4 := m.AvailableTransitions(pending, reviewedAndBlocked)
	if len(avail4) != 2 {
		t.Errorf("expected 2 transitions for Reviewed+Blocked, got %d", len(avail4))
	}

	// Instance follows AnyOf: advance with Approved succeeds.
	inst := NewInstance(m)
	if err := inst.Advance(approvedOnly); err != nil {
		t.Fatalf("Advance with Approved should succeed: %v", err)
	}
	if !inst.Current().Equal(accepted) {
		t.Errorf("expected accepted, got %q", inst.Current().Name)
	}
}

// TestNotCombinatorBlocksTransition tests that a Not combinator prevents a
// transition when the negated condition holds.
func TestNotCombinatorBlocksTransition(t *testing.T) {
	ready := State{Name: "ready"}
	deployed := State{Name: "deployed"}

	testsPassed := artifactRule("tests_passed", func(a artifact) bool { return a.TestsPassed })
	isBlocked := artifactRule("is_blocked", func(a artifact) bool { return a.IsBlocked })

	// Guard: tests must pass AND the artifact must NOT be blocked.
	guard := AllOf[artifact](testsPassed, Not[artifact](isBlocked))

	m := NewMachine[artifact](ready, []Transition[artifact]{
		{From: ready, To: deployed, Guard: guard},
	})

	// Tests passed, not blocked: transition available.
	good := artifact{TestsPassed: true, IsBlocked: false}
	avail := m.AvailableTransitions(ready, good)
	if len(avail) != 1 || !avail[0].To.Equal(deployed) {
		t.Errorf("expected ready→deployed for good artifact, got %v", avail)
	}

	// Tests passed but blocked: Not(isBlocked) fails, so no transition.
	blockedPassing := artifact{TestsPassed: true, IsBlocked: true}
	avail2 := m.AvailableTransitions(ready, blockedPassing)
	if len(avail2) != 0 {
		t.Errorf("expected 0 transitions for blocked artifact (Not should block), got %d", len(avail2))
	}

	// Tests not passed, not blocked: AllOf fails on testsPassed.
	noTests := artifact{TestsPassed: false, IsBlocked: false}
	avail3 := m.AvailableTransitions(ready, noTests)
	if len(avail3) != 0 {
		t.Errorf("expected 0 transitions when tests have not passed, got %d", len(avail3))
	}

	// Instance: blocked artifact cannot advance.
	inst := NewInstance(m)
	if err := inst.Advance(blockedPassing); err == nil {
		t.Error("Advance should fail for blocked artifact")
	}
	if !inst.Current().Equal(ready) {
		t.Errorf("state should remain ready after failed Advance, got %q", inst.Current().Name)
	}

	// Instance: good artifact can advance.
	if err := inst.Advance(good); err != nil {
		t.Fatalf("Advance should succeed for good artifact: %v", err)
	}
	if !inst.Current().Equal(deployed) {
		t.Errorf("expected deployed, got %q", inst.Current().Name)
	}
}
