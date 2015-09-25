package automata

import ()

type NFA struct {
	//开始状态startState
	startState *NFAState

	//接收状态acceptingStates
	acceptingStates Set_NFAState //= new HashSet<NFAState>();
}

func NewNFA0() *NFA {
	return NewNFA2(NewNFAState(), NewNFAState())
}

func NewNFA1(startState *NFAState) *NFA {
	return NewNFA2(startState, NewNFAState())
}

func NewNFA2(startState, acceptingState *NFAState) *NFA {
	this := &NFA{}
	this.startState = startState
	this.acceptingStates = make(Set_NFAState)
	this.AddAcceptingState(acceptingState)
	return this
}

func (this *NFA) GetStartState() *NFAState { return this.startState }

func (this *NFA) GetAcceptingStates() Set_NFAState { return this.acceptingStates }

func (this *NFA) Accept(state *NFAState) bool {
	_, present := this.acceptingStates[state]
	return present
}

func (this *NFA) AddAcceptingState(state *NFAState) {
	this.acceptingStates[state] = state
}

//在上面的NFAState类实现中，新的状态节点是在添加迁移映射的过程中生成的，
//这个过程中NFA并没有介入，因此NFA类不能直接得到状态集S的成员
//而是需要从状态startState开始，不断迭代找出所有的状态节点
func (this *NFA) GetStateSet2(current *NFAState, states Set_NFAState) {
	if _, present := states[current]; present {
		return
	}

	states[current] = current

	allstates := current.GetNextStates()
	for _, state := range allstates {
		this.GetStateSet2(state, states)
	}
}

func (this *NFA) GetStateSet() Set_NFAState {
	states := make(Set_NFAState) //new HashSet<NFAState>();
	this.GetStateSet2(this.GetStartState(), states)
	return states
}
