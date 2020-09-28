package nfa

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecision(t *testing.T) {
	state := State("Success")
	d, _ := NewDecision(func() State {
		return state
	})

	received, _ := executeDecision(d)

	assert.Equal(t, state, received)
}