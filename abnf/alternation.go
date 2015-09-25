package abnf

import (
	"GoABNF/automata"
	"bytes"
	"container/list"
)

//  alternation    =  concatenation
//                          *(*c-wsp "/" *c-wsp concatenation)
type Alternation struct {
	concatenations *list.List
}

func NewAlternation() *Alternation {
	this := &Alternation{}
	this.concatenations = list.New()
	return this
}

func (this *Alternation) AddConcatenation(concatenation *Concatenation) {
	this.concatenations.PushBack(concatenation)
}

func (this *Alternation) GetConcatenations() *list.List {
	return this.concatenations
}

func (this *Alternation) String() string {
	var s bytes.Buffer

	for e := this.concatenations.Front(); e != nil; e = e.Next() {
		v := e.Value.(*Concatenation)
		s.WriteString(v.String())
		if e.Next() != nil {
			s.WriteString("/")
		}
	}

	return s.String()
}

func (this *Alternation) GetDependentRuleNames() Set_RuleName {
	ruleNames := make(Set_RuleName)
	for e := this.concatenations.Front(); e != nil; e = e.Next() {
		v := e.Value.(*Concatenation)
		s := v.GetDependentRuleNames()
		for _, r := range s {
			ruleNames[r.rulename] = r
		}
	}
	return ruleNames
}

func (this *Alternation) GetNFA(rules map[string]*Rule) *automata.NFA { //throws IllegalAbnfException {
	startState := automata.NewNFAState()
	acceptingState := automata.NewNFAState()
	this.GetNFAStates(startState, acceptingState, rules)
	return automata.NewNFA2(startState, acceptingState)
}

func (this *Alternation) GetNFAStates(startState, acceptingState *automata.NFAState, rules map[string]*Rule) {
	if this.concatenations.Len() == 0 {
		panic("Alternation is empty.")
	}

	for e := this.concatenations.Front(); e != nil; e = e.Next() {
		e.Value.(*Concatenation).GetNFAStates(startState, acceptingState, rules)
	}
}
