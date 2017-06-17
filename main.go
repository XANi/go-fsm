package tsm

import (
)


type FSM struct {
	states map[int]*State
	state int
}

// State() returns existing state struct.
func (f *FSM) State(s int) *State {
	if state, ok :=  f.states[s] ; ok {
		return state
	} else {
		f.states[s] =&State{
			fsm: f,
			id: s,
			next: make(map[int]func() bool),
		}
		return f.states[s]
	}
}

func (f *FSM) Go(s int) bool {
	if verifyFunc, ok := f.states[f.state].next[s] ; ok {
		if verifyFunc == nil {
			f.state = s
			return true
		} else {
			isOk := verifyFunc()
			if isOk {
				f.state = s
			}
			return isOk
		}
	} else {
		return false
	}
}

type State struct {
	fsm *FSM
	next map[int]func() bool
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
	Condition func() bool
}

func New(startingState int, stateTable []Transitions) (*FSM, error) {
	f :=  FSM{
		states: make(map[int]*State),
		state: startingState,
	}
	for _, transition := range stateTable {
		for _, trTo := range transition.To {
			f.State(transition.From).next[trTo] = transition.Condition
		}
	}
	return &f, nil
}
