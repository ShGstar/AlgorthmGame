package stategraph

import "fmt"

//状态接口
type IState interface {
	Start()
	Update()
	Stop()
	Value() string
}

type FSM struct {
	state      IState
	registers  map[string]IState
	def_state  string
	input_hour uint32
}

func NewStateMachine() *FSM {
	fsm := &FSM{}
	fsm.registers = make(map[string]IState)
	return fsm
}

//
func (fsm *FSM) Process(hour uint32) error {
	fsm.input_hour = hour
	if fsm.state == nil {

	} else {
		fsm.state.Update()
	}
	return nil
}

//切换状态
func (fsm *FSM) ChangeState(value string) {
	if fsm.state == nil {
		fsm.state = fsm.registers[value]
	} else if fsm.state.Value() != value {
		fsm.state.Stop()
		fsm.state = fsm.registers[value]
	} else {
		fmt.Printf("same state\n")
		return
	}

	fsm.state.Start()
}

//
type StateBase struct {
	fsm   *FSM
	value string
}

func (b *StateBase) Start() {
	fmt.Println("-----> StateBase Start() <----------")
}

func (b *StateBase) Stop() {
	fmt.Println("-----> StateBase Stop() <----------")
}

func (b *StateBase) Update() {
	fmt.Println("-----> StateBase Update() <----------")
}

func (b *StateBase) Value() string {
	return b.value
}
