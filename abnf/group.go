package abnf

import (
	"GoABNF/automata"
)

//  group          =  "(" *c-wsp alternation *c-wsp ")"

type Group struct { //implements Element {
	alternation *Alternation
}

func NewGroup(alternation *Alternation) *Group {
	this := &Group{}
	this.alternation = alternation
	return this
}

func (this *Group) String() string{
	return "("+this.alternation.String()+")"
}

func (this *Group) GetElementType() ElementType{
	return ELEMENT_GROUP;
}

func (this *Group) GetDependentRuleNames() Set_RuleName {
	return this.alternation.GetDependentRuleNames();
}

func (this *Group) GetNFA(rules map[string]*Rule) *automata.NFA { //throws IllegalAbnfException {
	startState := automata.NewNFAState()
	acceptingState := automata.NewNFAState()
	this.GetNFAStates(startState, acceptingState, rules)
	return automata.NewNFA2(startState, acceptingState)
}

func (this *Group) GetNFAStates(startState, acceptingState *automata.NFAState, rules map[string]*Rule) {
	this.alternation.GetNFAStates(startState, acceptingState, rules);
}