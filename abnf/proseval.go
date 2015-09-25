package abnf

import (
	"GoABNF/automata"
)

//prose-val      =  "<" *(%x20-3D / %x3F-7E) ">"
//                       ; bracketed string of SP and VCHAR
//                          without angles
//                       ; prose description, to be used as
//                          last resort

type ProseVal struct { //implements Element {
	value string
}

func NewProseVal(value string) *ProseVal {
	this := &ProseVal{}
	this.value = value
	return this
}

func (this *ProseVal) String() string {
	return this.value
}

func (this *ProseVal) GetElementType() ElementType {
	return ELEMENT_PROSEVAL
}

func (this *ProseVal) GetDependentRuleNames() Set_RuleName {
	return make(Set_RuleName)
}

func (this *ProseVal) GetNFA(rules map[string]*Rule) *automata.NFA { //throws IllegalAbnfException {
	startState := automata.NewNFAState()
	acceptingState := automata.NewNFAState()
	this.GetNFAStates(startState, acceptingState, rules)
	return automata.NewNFA2(startState, acceptingState)
}

func (this *ProseVal) GetNFAStates(startState, acceptingState *automata.NFAState, rules map[string]*Rule) {
	if len(this.value) == 0 {
		startState.AddTransitEpsilon(acceptingState)
		return
	}

	current := startState
	buffer := []byte(this.value)
	for j := 0; j < len(buffer); j++ {
		if j < len(buffer)-1 {
			current = current.AddTransitByte1(buffer[j])
		} else {
			current = current.AddTransitByte2(buffer[j], acceptingState)
		}
	}
}
