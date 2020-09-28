# NFA: Non-deterministic Finite Automaton / Non-deterministic Finite-state Machine

> In automata theory, a finite-state machine is called a deterministic finite automaton (DFA), if
>   - each of its transitions is uniquely determined by its source state and input symbol,
>   - and reading an input symbol is required for each state transition.
>
> A nondeterministic finite automaton (NFA), or nondeterministic finite-state machine, does not need to obey these restrictions. In particular, every DFA is also an NFA.

- [Wikipedia](https://en.wikipedia.org/wiki/Nondeterministic_finite_automaton)

A simple NFA implementation written in Go.

# Usage

## Deterministic Finite-state Machine Example

> Every DFA is also an NFA

```go
import "github.com/chrisdoherty4/nfa"

var (
    PendingState = nfa.State("Pending")
    RunningState = nfa.State("Running")
    CompleteState = nfa.State("Complete")
)

var (
    StartEvent = nfa.Event("Start")
    FinishEvent = nfa.Event("Finish")
)

func main() {
    machine := nfa.NewMachine(PendingState)

    machine.TransitionD(PendingState, StartEvent, RunningState)
    machine.TransitionD(RunningState, FinishEvent, CompleteState)

    machine.Event(StartEvent)
    fmt.Println(machine.State()) // Running

    machine.Event(FinishEvent)
    fmt.Println(machine.State()) // Complete
}
```

## Non-deterministic Finite-state Machine Example


```go
import "github.com/chrisdoherty4/nfa"

var (
    PendingState = nfa.State("Pending")
    RunningState = nfa.State("Running")
    SuccessState = nfa.State("Complete")
    ErrorState = nfa.State("Error")
)

var (
    StartEvent = nfa.Event("Start")
    CompleteEvent = nfa.Event("Complete")
)

func main() {
    machine := nfa.NewMachine(PendingState)

    machine.TransitionD(PendingState, StartEvent, RunningState)
    machine.Transition(RunningState, FinishEvent, func(result bool) State {
        if result {
            return SuccessState
        }

        return ErrorState
    })

    machine.Event(StartEvent)
    fmt.Println(machine.State()) // Running

    machine.Event(CompleteEvent, false)
    fmt.Println(machine.State()) // Error
}
```