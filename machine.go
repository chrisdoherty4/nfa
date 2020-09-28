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

// Machine is a non-deterministic finite automaton.
type Machine struct {
	// A map of states to events they can transition away from.
	transitions  map[State]events

	// The current state of the machine.
	currentState State
}

// NewMachine creates a new FSM with the `startingState`. If no transitions are
// registerred against the `startingState`, the FSM will never leave the
// `startingState`.
func NewMachine(startingState State) Machine {
	return Machine{
		transitions:  make(map[State]events),
		currentState: startingState,
	}
}

// TransitionD registers a deterministic transition from `currentState` to
// `newState` when `event` occurs.
func (m *Machine) TransitionD(currentState State, event Event, newState State) error {
	if currentState == "" {
		return errors.New("Invalid current state")
	}

	if newState == "" {
		return errors.New("Invalid new state")
	}

	return m.addTransition(
		currentState,
		event,
		transition{
			state: newState,
		},
	)
}

// Transition registers a non-deterministic transition from `currentState`
// when `event` occurs.
func (m *Machine) Transition(currentState State, event Event, decision Decision) error {
	if currentState == "" {
		return errors.New("Invalid current state")
	}

	return m.addTransition(
		currentState,
		event,
		transition{
			decision: decision,
		},
	)
}

func (m *Machine) addTransition(state State, event Event, transition transition) error {
	if m.events(state) != nil {
		return fmt.Errorf(
			"Transition already registered for %v",
			string(state),
		)
	}

	m.transitions[state] = make(events)
	m.transitions[state][event] = transition

	return nil
}

// Event triggers a transition. Transitions are only triggered if the current
// state has the `even` registered as a path to a new state. If the `event` is
// registered, the FSM will transition to the new state.
func (m *Machine) Event(event Event, params ...interface{}) error {
	// Retrieve the events
	// If there are no events we've reached an end state so return an error.
	events := m.events(m.currentState)
	if events == nil {
		return fmt.Errorf("End state reached (%v)", m.currentState)
	}

	// If there are events check to see if this state can handle that event
	// and if not return an error
	if _, ok := events[event]; !ok {
		return errors.New("Nothing to do")
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

func (m *Machine) events(state State) events {
	if e, ok := m.transitions[state]; ok {
		return e
	}

	return nil
}

// transition represents the method of transition. Either the transition is
// deterministic and `state` is set, or the transition is non-determinisitic
// and `decision` is set.
type transition struct {
	state    State
	decision Decision
}

// events wraps up the Event -> transition mapping for a given state.
type events map[Event]transition

// Decision represents the decision to be made for a non-deterministic
// transition.
type Decision interface{}

// NewDecision performs type checking on the decision function and returns the
// Decision type.
func NewDecision(f interface{}) (Decision, error) {
	if reflect.TypeOf(f).Kind() != reflect.Func {
		return nil, errors.New("f must be of type reflect.Func")
	}

	return Decision(f), nil
}

// executeDecision executes a decision using the supplied `params`. If the
// number of params or the type of param at a given index does not match the
// `Decision` supplied with `m.Transition(...)` an error will occur.
func executeDecision(d Decision, params ...interface{}) (State, error) {
	decisionValue := reflect.ValueOf(d)
	decisionType := decisionValue.Type()

	if decisionType.NumIn() != len(params) {
		return "", errors.New("Invalid number of params")
	}

	// Iterate over the decision parameters and ensure the kind is the same as
	// the supplied params. Build slice of correct interface types.
	args := make([]reflect.Value, decisionType.NumIn())

	for i := 0; i < decisionType.NumIn(); i++ {
		kind := reflect.TypeOf(decisionType.In(i)).Kind()
		paramValue := reflect.ValueOf(params[i])

		// Ensure the parameter kind is the same. For types that may be nil
		// such as pointers or channels we do not check the underlying type
		// used because there could be slices of slices and we wouldn't know
		// where to stop.
		if kind != paramValue.Kind() {
			return "", fmt.Errorf("Expected %v for param %v", kind, i)
		}

		args = append(args, paramValue)
	}

	result := decisionValue.Call((args))

	if state, ok := result[0].Interface().(State); ok {
		return state, nil
	}

	return "", errors.New("Decision func returned wrong type")
}
