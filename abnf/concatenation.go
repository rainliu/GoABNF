package abnf

import (
	"bytes"
	"container/list"
	"GoABNF/automata"
)

// concatenation  =  repetition *(1*c-wsp repetition)
type Concatenation struct {
	repetitions *list.List
}

func NewConcatenation() *Concatenation {
	this := &Concatenation{}
	this.repetitions = list.New()
	return this
}

func (this *Concatenation) AddRepetition(repetition *Repetition) {
	this.repetitions.PushBack(repetition)
}

func (this *Concatenation) GetRepetitions() *list.List {
	return this.repetitions
}

func (this *Concatenation) String() string {
	var s bytes.Buffer

	for e := this.repetitions.Front(); e != nil; e = e.Next() {
		v := e.Value.(*Repetition)
		s.WriteString(v.String())
		if e.Next() != nil {
			s.WriteString(" ")
		}
	}

	return s.String()
}

func (this *Concatenation) GetDependentRuleNames() Set_RuleName {
	ruleNames := make(Set_RuleName)
	for e := this.repetitions.Front(); e != nil; e = e.Next() {
		v := e.Value.(*Repetition)
		s := v.GetDependentRuleNames()
		for _, r := range s {
			ruleNames[r.rulename] = r
		}
	}
	return ruleNames
}

func (this *Concatenation) GetNFA(rules map[string]*Rule) *automata.NFA { //throws IllegalAbnfException {
	startState := automata.NewNFAState()
	acceptingState := automata.NewNFAState()
	this.GetNFAStates(startState, acceptingState, rules)
	return automata.NewNFA2(startState, acceptingState)
}

func (this *Concatenation) GetNFAStates(startState, acceptingState *automata.NFAState, rules map[string]*Rule) {
	current := startState
	e := this.repetitions.Front()
	var next *automata.NFAState
	for index := 0; index < this.repetitions.Len()-1; index++ {
		next = automata.NewNFAState()
		e.Value.(*Repetition).GetNFAStates(current, next, rules)
		current = next
		e = e.Next()
	}
	e.Value.(*Repetition).GetNFAStates(current, acceptingState, rules)
}
