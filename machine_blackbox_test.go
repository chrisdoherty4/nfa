package nfa_test

import (
	"testing"

	"github.com/chrisdoherty4/nfa"
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

	t.Run("GraphFSM", func(t *testing.T) {
		var (
			S1 nfa.State = "S1"
			S2 nfa.State = "S2"
			S3 nfa.State = "S3"
			E1 nfa.Event = "E1"
		)

		machine := NewMachine(
			S1,
			Transitions{
				S1: Events{
					E1: NewTransition(func(result bool) State {
						if result {
							return S2
						}

						return S3
					}),
				},

				S2: Events{
					E1: NewTransitionD(S1),
				},

				S3: Events{
					E1: NewTransitionD(S1),
				},
			},
		)

		assert.Equal(t, S1, machine.State())

		machine.Transition(E1, true)
		assert.Equal(t, S2, machine.State())

		machine.Transition(E1)
		assert.Equal(t, S1, machine.State())

		machine.Transition(E1, false)
		assert.Equal(t, S3, machine.State())

		machine.Transition(E1)
		assert.Equal(t, S1, machine.State())
	})
}
