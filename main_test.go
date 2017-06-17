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
	states := []Transitions{
		{From: StStart, To: []int{StPr1Step1, StPr2Step1}},
		{From: StPr1Step1, To: []int{StPr1Step1, StPr2Step2}},
		{From: StPr1Step2, To: []int{StPr1Step3}},
		{From: StPr1Step3, To: []int{StEnd}},
		{From: StStart, To: []int{StPr1Step1, StPr2Step1}},
		{From: StPr2Step1, To: []int{StPr2Step2}},
		{From: StPr2Step2, To: []int{StPr1Step3}},
		{From: StPr2Step3, To: []int{StEnd}},
	}
	sm, _ := New(StStart,states)
	Convey("GetNextState",t,func() {
		So(	sm.State(StStart).Next(),ShouldContain,StPr1Step1)
		So(	sm.State(StStart).Next(),ShouldContain,StPr2Step1)
	})

}
