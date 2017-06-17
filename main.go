package tsm

import (
)


type FSM struct {
	states map[int]*State
}

// State() returns existing state struct.
func (f *FSM) State(s int) *State {
	if state, ok :=  f.states[s] ; ok {
		return state
	} else {
		f.states[s] =&State{
			fsm: f,
			id: s,
			next: make(map[int]int),
		}
		return f.states[s]
	}
}

type State struct {
	fsm *FSM
	next map[int]int
	id int
}

// Return list of next possible states
func (st *State) Next() []int {
	state := make([]int, len(st.next))
	i := 0
	for k := range st.next {
		state[i] = k
		i++
	}
	return state
}


type Transitions struct{
	From int
	To []int
	Condition *func() bool
}

func New(startingState int, stateTable []Transitions) (*FSM, error) {
	f :=  FSM{
		states: make(map[int]*State),
	}
	return &f, nil
}
