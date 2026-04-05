package rule

import (
	"fmt"
	"time"
)

// State represents a named state in a state machine.
// Equality is based on Name only; Metadata is informational.
type State struct {
	Name     string
	Metadata map[string]any
}

// Equal reports whether two States are equal by name.
func (s State) Equal(other State) bool {
	return s.Name == other.Name
}

// Transition defines a guarded edge between two states.
// Guard is evaluated against the subject to determine if the transition is allowed.
type Transition[T any] struct {
	From  State
	To    State
	Guard Rule[T]
}

// AlwaysAllow returns a Rule whose Check always returns a passing Verdict.
func AlwaysAllow[T any]() Rule[T] {
	return New(func(_ T) Verdict {
		return Pass()
	})
}

// Machine[T] is an immutable state machine definition.
// It stores the initial state and the full set of transitions.
type Machine[T any] struct {
	initial     State
	transitions []Transition[T]
}

// NewMachine creates a Machine with the given initial state and transitions.
// The transitions slice is copied so the machine is immutable after construction.
func NewMachine[T any](initial State, transitions []Transition[T]) *Machine[T] {
	ts := make([]Transition[T], len(transitions))
	copy(ts, transitions)
	return &Machine[T]{initial: initial, transitions: ts}
}

// AvailableTransitions returns all transitions from current whose guard passes for candidate.
func (m *Machine[T]) AvailableTransitions(current State, candidate T) []Transition[T] {
	var result []Transition[T]
	for _, t := range m.transitions {
		if !t.From.Equal(current) {
			continue
		}
		if t.Guard.Check(candidate).OK {
			result = append(result, t)
		}
	}
	return result
}

// TerminalStates returns all states that have no outgoing transitions.
// A state is terminal if it appears in any transition (as From or To) but never as From.
func (m *Machine[T]) TerminalStates() []State {
	mentioned := map[string]State{}
	hasOutgoing := map[string]bool{}

	// include initial state in the graph
	mentioned[m.initial.Name] = m.initial

	for _, t := range m.transitions {
		mentioned[t.From.Name] = t.From
		mentioned[t.To.Name] = t.To
		hasOutgoing[t.From.Name] = true
	}

	var terminals []State
	for name, state := range mentioned {
		if !hasOutgoing[name] {
			terminals = append(terminals, state)
		}
	}
	return terminals
}

// HistoryEntry records a single state transition with its timestamp.
type HistoryEntry struct {
	From      State
	To        State
	Timestamp time.Time
}

// Instance[T] is a running instance of a Machine[T].
// It tracks the current state and transition history.
// Instance is not safe for concurrent use; it follows a single-owner pattern.
type Instance[T any] struct {
	machine *Machine[T]
	current State
	history []HistoryEntry
}

// NewInstance creates an Instance starting at the machine's initial state.
func NewInstance[T any](machine *Machine[T]) *Instance[T] {
	return &Instance[T]{
		machine: machine,
		current: machine.initial,
	}
}

// Current returns the current state.
func (i *Instance[T]) Current() State {
	return i.current
}

// History returns a copy of the transition history.
func (i *Instance[T]) History() []HistoryEntry {
	h := make([]HistoryEntry, len(i.history))
	copy(h, i.history)
	return h
}

// Advance evaluates candidate against all available transitions from the current state
// and applies the first passing one. Returns an error if no transition is available.
func (i *Instance[T]) Advance(candidate T) error {
	available := i.machine.AvailableTransitions(i.current, candidate)
	if len(available) == 0 {
		return fmt.Errorf("no available transition from state %q", i.current.Name)
	}
	t := available[0]
	i.history = append(i.history, HistoryEntry{From: i.current, To: t.To, Timestamp: time.Now()})
	i.current = t.To
	return nil
}

// AdvanceTo finds the transition from the current state to target, checks its guard,
// and applies it. Returns an error if the transition is not found or the guard fails.
func (i *Instance[T]) AdvanceTo(target State, candidate T) error {
	for _, t := range i.machine.transitions {
		if !t.From.Equal(i.current) || !t.To.Equal(target) {
			continue
		}
		v := t.Guard.Check(candidate)
		if !v.OK {
			return fmt.Errorf("transition %q→%q blocked: %v", i.current.Name, target.Name, v.Context)
		}
		i.history = append(i.history, HistoryEntry{From: i.current, To: t.To, Timestamp: time.Now()})
		i.current = t.To
		return nil
	}
	return fmt.Errorf("no transition from %q to %q", i.current.Name, target.Name)
}

// IsTerminal reports whether the current state is a terminal state.
func (i *Instance[T]) IsTerminal() bool {
	for _, s := range i.machine.TerminalStates() {
		if s.Equal(i.current) {
			return true
		}
	}
	return false
}

// Validate checks the machine for structural errors.
// It returns an error if any state is not reachable from the initial state.
func (m *Machine[T]) Validate() error {
	// Build adjacency: From → []To
	adj := map[string][]string{}
	for _, t := range m.transitions {
		adj[t.From.Name] = append(adj[t.From.Name], t.To.Name)
	}

	// Collect all states mentioned in transitions plus initial
	all := map[string]bool{}
	all[m.initial.Name] = true
	for _, t := range m.transitions {
		all[t.From.Name] = true
		all[t.To.Name] = true
	}

	// BFS from initial
	visited := map[string]bool{m.initial.Name: true}
	queue := []string{m.initial.Name}
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		for _, next := range adj[cur] {
			if !visited[next] {
				visited[next] = true
				queue = append(queue, next)
			}
		}
	}

	var unreachable []string
	for name := range all {
		if !visited[name] {
			unreachable = append(unreachable, name)
		}
	}
	if len(unreachable) > 0 {
		return fmt.Errorf("unreachable states: %v", unreachable)
	}
	return nil
}
