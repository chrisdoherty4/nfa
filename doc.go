/*
Package nfa is a Non-deterministic Finite Automaton a.k.a a Non-detrministic
finite state machine

A `Machine` defines is the nfa that holds the current state. Using the
`Transition(...)` and `TransitionD(...)` methods we can register how
the NFA transitions from 1 state to another.

`TransitionD(...)` defines deterministic transitions. `Transition(...)`
defines non-deterministic transitions.

Example

	var (
		PendingState = nfa.State("Pending")
		RunningState = nfa.State("Running")
		SuccessState = nfa.State("Success")
		ErrorState = nfa.State("Error")
	)

	var (
		StartEvent = nfa.Event("Start")
		FinishEvent = nfa.Event("Finish")
	)

	machine := NewMachine(PendingState)
	machine.TransitionD(PendingState, StartEvent, RunningState)
	machine.Transition(RunningState, FinishEvent, func(result bool) nfa.State {
		if result {
			return SuccessState
		}

		return ErrorState
	})

	machine.Event(StartEvent) // machine.State() == RunningState

	machine.Event(FinishEvent, true) // machine.State() == SuccessState

*/
package nfa