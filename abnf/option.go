package abnf

import (
	"GoABNF/automata"
)

//  option         =  "[" *c-wsp alternation *c-wsp "]"

type Option struct { //implements Element {
	alternation *Alternation
}

func NewOption(alternation *Alternation) *Option {
	this := &Option{}
	this.alternation = alternation
	return this
}

func (this *Option) GetAlternation() *Alternation {
	return this.alternation
}

func (this *Option) String() string {
	return "[" + this.alternation.String() + "]"
}

func (this *Option) GetElementType() ElementType {
	return ELEMENT_OPTION
}

func (this *Option) GetDependentRuleNames() Set_RuleName {
	return this.alternation.GetDependentRuleNames()
}

func (this *Option) GetNFA(rules map[string]*Rule) *automata.NFA { //throws IllegalAbnfException {
	startState := automata.NewNFAState()
	acceptingState := automata.NewNFAState()
	this.GetNFAStates(startState, acceptingState, rules)
	return automata.NewNFA2(startState, acceptingState)
}

func (this *Option) GetNFAStates(startState, acceptingState *automata.NFAState, rules map[string]*Rule) {
	startState.AddTransitEpsilon(acceptingState)
	this.alternation.GetNFAStates(startState, acceptingState, rules)
}
