# Nondeterministic Finite Automaton (NFA) / Nondeterministic Finite-state Machine (NFM)

> In automata theory, a finite-state machine is called a deterministic finite automaton (DFA), if
>   - each of its transitions is uniquely determined by its source state and input symbol,
>   - and reading an input symbol is required for each state transition.
>
> A nondeterministic finite automaton (NFA), or nondeterministic finite-state machine, does not need to obey these restrictions. In particular, every DFA is also an NFA.

\- [Wikipedia](https://en.wikipedia.org/wiki/Nondeterministic_finite_automaton)

A simple NFA implementation written in Go.

## Usage

```go
import "github.com/chrisdoherty4/nfa"

const (
    PendingState  nfa.State = "Pending"
    RunningState  nfa.State = "Running"
    CompleteState nfa.State = "Complete"
    ErrorState    nfa.State = "Error"
)

const (
    StartEvent  nfa.Event = "Start"
    FinishEvent nfa.Event = "Finish"
)

func main() {
    machine := nfa.NewMachine(
        PendingState,
        nfa.Transitions{
            PendingState: nfa.Events{
                StartEvent: nfa.NewTransitionD(RunningState),
            },

            RunningState: nfa.Events{
                FinishEvent: nfa.NewTransition(func(result bool) {
                    if result {
                        return CompleteState
                    }

                    return ErrorState
                }),
            },
        },
    )

    machine.Transition(StartEvent)
    fmt.Println(machine.State()) // Running

    machine.Transition(FinishEvent, true)
    fmt.Println(machine.State()) // Complete
}
```
