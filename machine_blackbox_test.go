package nfa_test

import (
	"testing"

	. "git.alteryx.com/eaas/engine-service/nfa"

	"github.com/stretchr/testify/assert"
)

func TestMachine(t *testing.T) {
	t.Run("NewMachine", func(t *testing.T) {
		_ = NewMachine(State("StartingSate"))
	})

	t.Run("MachineStartingState", func(t *testing.T) {
		state := State("Starting")
		machine := NewMachine(state)

		assert.Equal(t, state, machine.State())
	})

	t.Run("DeterministicTransition", func(t *testing.T) {
		startState := State("Starting")
		finalState := State("Final")
		event := Event("finalize")

		machine := NewMachine(startState)

		err := machine.TransitionD(startState, event, finalState)
		assert.Nil(t, err, err)

		err = machine.Event(event)
		assert.Nil(t, err, err)

		assert.Equal(t, finalState, machine.State())
	})

	t.Run("NonDeterministicTransition", func(t *testing.T) {
		startState := State("Starting")
		finalState := State("Final")
		event := Event("finalize")

		machine := NewMachine(startState)

		err := machine.Transition(startState, event, func() State {
			return finalState
		})

		assert.Nil(t, err, err)

		err = machine.Event(event)
		assert.Nil(t, err, err)

		assert.Equal(t, finalState, machine.State())
	})
}