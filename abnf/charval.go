package abnf

import (
	"GoABNF/automata"
)

//char-val       =  DQUOTE *(%x20-21 / %x23-7E) DQUOTE
//                       ; quoted string of SP and VCHAR
//                          without DQUOTE

type CharVal struct { //implements Element {
	value string
}

func NewCharVal(value string) *CharVal {
	this := &CharVal{}
	this.value = value
	return this
}

func (this *CharVal) String() string {
	return "\"" + this.value + "\""
}

func (this *CharVal) GetElementType() ElementType {
	return ELEMENT_CHARVAL
}

//@Override
func (this *CharVal) GetDependentRuleNames() Set_RuleName {
    return make(Set_RuleName);
}

//@Override
func (this *CharVal) GetNFA(rules map[string]*Rule) *automata.NFA { //throws IllegalAbnfException {
	startState := automata.NewNFAState()
	acceptingState := automata.NewNFAState()
	this.GetNFAStates(startState, acceptingState, rules)
	return automata.NewNFA2(startState, acceptingState)
}

//@Override
func (this *CharVal) GetNFAStates(startState, acceptingState *automata.NFAState, rules map[string]*Rule){
	//若CharVal的内容为空，则创建一条从开始状态到接受状态的epsilon迁移。
	if len(this.value) == 0 {
		startState.AddTransitEpsilon(acceptingState)
		return
	}

	current := startState
	buffer := []byte(this.value)
	for j := 0; j < len(buffer); j++ {
		//最后一个节点使用方法参数中的acceptingState，中间节点自行创建
		if j < len(buffer)-1 {
			current = current.AddTransitByte1(buffer[j])
		} else {
			current = current.AddTransitByte2(buffer[j], acceptingState)
		}
	}
	return
}
