package state_machine

import (
	"errors"
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/multi"
	"log"
)

type Event string
type Operator string

var NodeIDCntr int64 = 0
var LineIDCntr int64 = 0

type StateMachine struct {
	PresentState  State
	previousState State
	graph         *multi.DirectedGraph
}

func New() *StateMachine {
	sm := &StateMachine{}
	sm.graph = multi.NewDirectedGraph()
	return sm
}

func (sm *StateMachine) Init(initStateValue interface{}) State {
	sm.PresentState = State{
		Id:    NodeIDCntr,
		Value: initStateValue,
	}
	sm.graph.AddNode(sm.PresentState)
	sm.previousState = sm.PresentState
	NodeIDCntr++
	return sm.PresentState
}

func (sm *StateMachine) PreviousState() State {
	return sm.previousState
}

func (sm *StateMachine) NewState(stateValue interface{}) State {
	state := State{
		Id:    NodeIDCntr,
		Value: stateValue,
	}
	sm.graph.AddNode(state)
	NodeIDCntr++
	return state
}

func (sm *StateMachine) LinkStates(s1, s2 State, rules map[Operator]Event) {
	sm.graph.SetLine(Link{from: s1, to: s2, id: LineIDCntr, rules: rules})
	LineIDCntr++
}

func NewRule(op Operator, e Event) map[Operator]Event {
	return map[Operator]Event{op: e}
}

func (sm *StateMachine) TriggerEvent(e Event) error {
	presentNode := sm.PresentState

	it := sm.graph.From(presentNode.Id)

	for it.Next() {
		node := sm.graph.Node(it.Node().ID()).(State)
		link := graph.LinesOf(sm.graph.Lines(presentNode.Id, node.Id))[0].(Link)

		for key, val := range link.rules {
			switch key {
			case eq:
				if val == e {
					sm.previousState = sm.PresentState
					sm.PresentState = node
					return nil
				}
			default:
				//TODO: implement other operator
				log.Printf("Operator %s is not supported\n", key)
				return errors.New("UNSUPPORTED_OPERATOR")
			}
		}
	}
	return nil
}

func (sm *StateMachine) RestoreState() State {
	sm.PresentState = sm.previousState
	return sm.PresentState
}

const (
	eq  Operator = "eq"
	neq Operator = "neq"
	le  Operator = "le"
	gr  Operator = "gr"
	lee Operator = "lee"
	gre Operator = "gre"
)

type State struct {
	Id    int64
	Value interface{}
}

type Link struct {
	id       int64
	from, to graph.Node
	rules    map[Operator]Event
}

func (st State) ID() int64 {
	return st.Id
}

func (l Link) From() graph.Node {
	return l.from
}

func (l Link) To() graph.Node {
	return l.to
}

func (l Link) ID() int64 {
	return l.id
}

func (l Link) ReversedLine() graph.Line {
	return Link{from: l.to, to: l.from}
}
