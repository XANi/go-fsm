package tsm

import (
	"testing"
	. "github.com/smartystreets/goconvey/convey"
)

const (
	StStart = iota
	StPr1Step1
	StPr1Step2
	StPr1Step3
	StPr2Step1
	StPr2Step2
	StPr2Step3
	StEnd
)


func TestState(t *testing.T) {
	sm := prepareTestObject()
	Convey("GetNextState",t,func() {
		So(	sm.State(StStart).Next(),ShouldContain,StPr1Step1)
		So(	sm.State(StStart).Next(),ShouldContain,StPr2Step1)
	})
	Convey("Transition to valid state",t,func() {
		So (sm.Go(StPr1Step1),ShouldBeTrue)
		So (sm.Go(StPr1Step2),ShouldBeTrue)
	})
	Convey("Transition to invalid state",t,func() {
		So (sm.Go(StStart),ShouldBeFalse)
	})
	sm = prepareTestObject()
	Convey("Validation Function",t,func() {
		So(sm.Go(StPr2Step1),ShouldBeTrue)
		So(sm.Go(StPr2Step2),ShouldBeFalse)


	})





}


func prepareTestObject() *FSM {
	states := []Transitions{
		{From: StStart, To: []int{StPr1Step1}, },
		{From: StPr1Step1, To: []int{StPr1Step1, StPr1Step2, StPr1Step2}},
		{From: StPr1Step2, To: []int{StPr1Step3}},
		{From: StPr1Step3, To: []int{StEnd}},
		{From: StStart, To: []int{StPr2Step1},Condition: func() bool {return true}},
		{From: StPr2Step1, To: []int{StPr2Step2,StPr2Step1}, Condition: func() bool {return false}},
		{From: StPr2Step2, To: []int{StPr1Step3}},
		{From: StPr2Step3, To: []int{StEnd}},
	}
	sm, _ := New(StStart,states)
	return sm
}


func BenchmarkNextState(b *testing.B) {
	sm := prepareTestObject()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		sm.State(StStart).Next()
	}
}
func BenchmarkStateTransition(b *testing.B) {
	sm := prepareTestObject()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		sm.Go(StPr1Step1)
	}
}

func BenchmarkStateTransitionWithCondition(b *testing.B) {
	sm := prepareTestObject()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		sm.Go(StPr2Step1)
	}
}
