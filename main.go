package fsm

import (
	"sync"
	"sync/atomic"
)

type StateID =int64

type FSM struct {
	states map[int64]*State
	state int64
	sync.Mutex
}

// State() returns existing state struct.
func (f *FSM) State(s int64) *State {
	if state, ok :=  f.states[s] ; ok {
		return state
	} else {
		f.states[s] =&State{
			fsm: f,
			id: s,
			next: make(map[int64]func() bool),
		}
		return f.states[s]
	}
}


// Go moves from stateA to stateB if transition is valid and passes any condition(if present), returns false if not possible or condition failed

func (f *FSM) Go(from, to StateID) bool {
	if verifyFunc, ok := f.states[from].next[to] ; ok {
		if verifyFunc == nil {
			return atomic.CompareAndSwapInt64(&f.state,from, to)
			return true
		} else {
			isOk := verifyFunc()
			if isOk {
				return atomic.CompareAndSwapInt64(&f.state,from, to)
			}
			return isOk
		}
	} else {
		return false
	}
}

// ToSerial() moves state machine to another state, returns true if success, false if
//
// * state is invalid
// * condition function returned false
//
// ToSerial() does not lock FSM. do not use it in concurrent use cases
// use To() for slower but synchronized state
func (f *FSM) ToSerial(s int64) bool {
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
// To() moves machine to specified state, returns true if success, false if
//
// * state is invalid
// * condition function returned false
//
// for slightly faster but concurrency-unsafe function, look at ToSerial().
//
// Note that it *doesn't care* about current state, as long as transition from that state to desired one
// is valid. If you want to be sure that desired action will be
// stateA -> stateB (and not anyValid -> stateB), use Go()
func (f *FSM) To(s int64) bool {
	f.Lock()
	// defer is hideously slow compared to alternative (adds extra ~60ns)
	z := f.ToSerial(s)
	f.Unlock()
	return z
}


type State struct {
	fsm *FSM
	next map[int64]func() bool
	onEntry []func(from, to StateID)
	onExit []func(from, to StateID)
	onEntryFrom map[int64][]func( to StateID)
	onExitTo map[int64][]func(from StateID)
	id int64
}

// Return list of next possible states
func (st *State) Next() []int64 {
	state := make([]int64, len(st.next))
	i := 0
	for k := range st.next {
		state[i] = k
		i++
	}
	return state
}

// Add function that will be ran when entering this state
func (st *State) AddEntryFunc(f func (from, to StateID)) *State {
	st.onEntry = append(st.onEntry, f)
	return st
}
// Add function that will be ran when exiting this state
func (st *State) AddExitFunc(f func (from, to StateID)) *State {
	st.onExit = append(st.onExit, f)
	return st
}
// Add function that will be ran if the state is entered from a certain other state
func (st *State) AddEntryFromFunc(state StateID, f func(to StateID)) *State {
	if _, ok := st.onEntryFrom[state] ; !ok {
		st.onEntryFrom[state] = make([]func(StateID),0)
	}
	st.onEntryFrom[state] = append(st.onEntryFrom[state],f)
	return st
}

// Add function that will be ran if the state is entered from a certain other state
func (st *State) AddExitToFunc(state StateID, f func(to StateID)) *State {
	if _, ok := st.onExitTo[state] ; !ok {
		st.onExitTo[state] = make([]func(StateID),0)
	}
	st.onExitTo[state] = append(st.onExitTo[state],f)
	return st
}



type Transitions struct{
	From int64
	To []int64
	Condition func() bool
}




func New(startingState int64, stateTable []Transitions) (*FSM, error) {
	f :=  FSM{
		states: make(map[int64]*State),
		state: startingState,
	}
	for _, transition := range stateTable {
		for _, trTo := range transition.To {
			f.State(transition.From).next[trTo] = transition.Condition
		}
	}
	return &f, nil
}
