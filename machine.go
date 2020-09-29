package nfa

import (
	"errors"
	"fmt"
	"reflect"
)

// Event defines an event that triggers a transition in the FSM.
type Event string

// State defines a state that a `Machine` can be in.
type State string

type Transitions map[State]Events

type Events map[Event]Transition

type Transition struct {
	state    State
	decision Decision
}

type Decision interface{}

// TransitionD registers a deterministic transition from `currentState` to
// `newState` when `event` occurs.
func NewTransitionD(state State) Transition {
	return Transition{state: state}
}

func NewTransition(decision Decision) Transition {
	return Transition{decision: decision}
}

// Machine is a non-deterministic finite automaton.
type Machine struct {
	// The current state of the machine.
	currentState State

	// A map of states to events they can transition away from.
	transitions Transitions
}

// NewMachine creates a new FSM with the `startingState`. If no transitions are
// registerred against the `startingState`, the FSM will never leave the
// `startingState`.
func NewMachine(startingState State, transitions Transitions) Machine {
	return Machine{
		currentState: startingState,
		transitions:  transitions,
	}
}

// Event triggers a transition. Transitions are only triggered if the current
// state has the `even` registered as a path to a new state. If the `event` is
// registered, the FSM will transition to the new state.
func (m *Machine) Transition(event Event, params ...interface{}) error {
	// Retrieve the events
	// If there are no events we've reached an end state so return an error.
	events := m.events(m.currentState)
	if events == nil {
		return fmt.Errorf("End state (%v)", m.currentState)
	}

	// If there are events check to see if this state can handle that event
	// and if not return an error
	if _, ok := events[event]; !ok {
		return errors.New("Invalid transition")
	}

	// If the transition is a basic transition, do it.
	if events[event].decision == nil {
		m.currentState = events[event].state
	}

	var newState State
	var err error

	switch {
	case events[event].decision == nil:
		newState = events[event].state
	default:
		newState, err = executeDecision(events[event].decision, params...)
	}

	if err != nil {
		return err
	}

	m.currentState = newState

	return nil
}

// State retrieves the current state of the FSM.
func (m *Machine) State() State {
	return m.currentState
}

func (m *Machine) events(state State) Events {
	if e, ok := m.transitions[state]; ok {
		return e
	}

	return nil
}

// executeDecision executes a decision using the supplied `params`. If the
// number of params or the type of param at a given index does not match the
// `Decision` supplied with `m.Transition(...)` an error will occur.
func executeDecision(d Decision, params ...interface{}) (State, error) {
	decisionValue := reflect.ValueOf(d)
	decisionType := decisionValue.Type()

	if decisionType.NumIn() != len(params) {
		return "", fmt.Errorf(
			"Invalid number of params, expected %v",
			decisionType.NumIn(),
		)
	}

	// Iterate over the decision parameters and ensure the kind is the same as
	// the supplied params. Build slice of correct interface types.
	args := make([]reflect.Value, decisionType.NumIn())

	for i := 0; i < decisionType.NumIn(); i++ {
		expectedKind := decisionType.In(i).Kind()
		paramValue := reflect.ValueOf(params[i])

		// Ensure the parameter kind is the same. For types that may be nil
		// such as pointers or channels we do not check the base type
		// used because there could be slices of slices and we wouldn't know
		// where to stop.
		if expectedKind != paramValue.Kind() {
			return "", fmt.Errorf("Expected '%v' for param %v", expectedKind, i)
		}

		args[i] = paramValue
	}

	result := decisionValue.Call((args))

	if state, ok := result[0].Interface().(State); ok {
		return state, nil
	}

	return "", errors.New("Decision func returned wrong type")
}
