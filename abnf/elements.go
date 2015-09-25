package abnf

import (
	"GoABNF/automata"
)

// elements       =  alternation *c-wsp
type Elements struct {
	alternation *Alternation
}

func NewElements(alternation *Alternation) *Elements {
	this := &Elements{}
	this.alternation = alternation
	return this
}

func (this *Elements) GetAlternation() *Alternation {
	return this.alternation
}

func (this *Elements) String() string{
	return this.alternation.String();
}

func (this *Elements) GetDependentRuleNames() Set_RuleName {
    return this.alternation.GetDependentRuleNames();
}

func (this *Elements) GetNFA(rules map[string]*Rule) *automata.NFA { //throws IllegalAbnfException {
	startState := automata.NewNFAState()
	acceptingState := automata.NewNFAState()
	this.GetNFAStates(startState, acceptingState, rules)
	return automata.NewNFA2(startState, acceptingState)
}

func (this *Elements) GetNFAStates(startState, acceptingState *automata.NFAState, rules map[string]*Rule) {
	this.alternation.GetNFAStates(startState, acceptingState, rules);
}