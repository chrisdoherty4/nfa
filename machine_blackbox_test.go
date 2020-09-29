package nfa_test

import (
	"testing"

	. "github.com/chrisdoherty4/nfa"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMachine(t *testing.T) {

	const (
		Pending  State = "Pending"
		Running  State = "Running"
		Complete State = "Complete"
		Error    State = "Error"
	)

	const (
		Start  Event = "Start"
		Finish Event = "Finish"
	)

	t.Run("NewMachine", func(t *testing.T) {
		machine := NewMachine(
			Pending,
			Transitions{
				Pending: Events{
					Start: NewTransitionD(Running),
				},
			},
		)

		assert.Equal(t, Pending, machine.State())
	})

	t.Run("DeterministicTransition", func(t *testing.T) {
		machine := NewMachine(
			Pending,
			Transitions{
				Pending: Events{
					Start: NewTransitionD(Running),
				},
			},
		)

		err := machine.Transition(Start)
		assert.Nil(t, err, err)

		assert.Equal(t, Running, machine.State())
	})

	t.Run("NondeterministicTransition", func(t *testing.T) {
		machine := NewMachine(
			Running,
			Transitions{
				Running: Events{
					Finish: NewTransition(func() State {
						return Complete
					}),
				},
			},
		)

		err := machine.Transition(Finish)
		require.Nil(t, err, err)

		assert.Equal(t, Complete, machine.State())
	})

	t.Run("MultipleTransitions", func(t *testing.T) {
		machine := NewMachine(
			Pending,
			Transitions{
				Pending: Events{
					Start: NewTransitionD(Running),
				},

				Running: Events{
					Finish: NewTransitionD(Complete),
				},
			},
		)

		machine.Transition(Start)
		assert.Equal(t, Running, machine.State())

		machine.Transition(Finish)
		assert.Equal(t, Complete, machine.State())
	})

	t.Run("DecisionParams", func(t *testing.T) {
		machine := NewMachine(
			Running,
			Transitions{
				Running: Events{
					Finish: NewTransition(func(state State) State {
						return state
					}),
				},
			},
		)

		err := machine.Transition(Finish, Complete)
		assert.Nil(t, err, err)

		assert.Equal(t, Complete, machine.State())
	})

	t.Run("IncorrectDecisionParams", func(t *testing.T) {
		machine := NewMachine(
			Running,
			Transitions{
				Running: Events{
					Finish: NewTransition(func(i int) State {
						return Complete
					}),
				},
			},
		)

		err := machine.Transition(Finish)
		assert.NotNil(t, err, err)
	})

	// t.Run("MachineStartingState", func(t *testing.T) {
	// 	state := State("Starting")
	// 	machine := NewMachine(state)

	// 	assert.Equal(t, state, machine.State())
	// })

	// t.Run("DeterministicTransition", func(t *testing.T) {
	// 	startState := State("Starting")
	// 	finalState := State("Final")
	// 	event := Event("finalize")

	// 	machine := NewMachine(startState)

	// 	err := machine.TransitionD(startState, event, finalState)
	// 	assert.Nil(t, err, err)

	// 	err = machine.Event(event)
	// 	assert.Nil(t, err, err)

	// 	assert.Equal(t, finalState, machine.State())
	// })

	// t.Run("NonDeterministicTransition", func(t *testing.T) {
	// 	startState := State("Starting")
	// 	finalState := State("Final")
	// 	event := Event("finalize")

	// 	machine := NewMachine(startState)

	// 	err := machine.Transition(startState, event, func() State {
	// 		return finalState
	// 	})

	// 	assert.Nil(t, err, err)

	// 	err = machine.Event(event)
	// 	assert.Nil(t, err, err)

	// 	assert.Equal(t, finalState, machine.State())
	// })
}
