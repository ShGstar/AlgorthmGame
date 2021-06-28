package main

import graph "./stategraph"

const (
	STATE_NO_1 = "STATE1"
	STATE_NO_2 = "STATE2"
	STATE_NO_3 = "STATE3"
	STATE_NO_4 = "STATE4"
	STATE_NO_5 = "STATE5"
)

type StateTest1 struct {
	graph.StateBase
}

func NewStateTest1(fsm *graph.FSM, value string) {

}

func main() {

}
