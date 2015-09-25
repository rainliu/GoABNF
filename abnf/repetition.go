package abnf

import (
	"GoABNF/automata"
)

// repetition     =  [repeat] element
type Repetition struct {
	repeat  *Repeat
	element Element
}

func NewRepetition(repeat *Repeat, element Element) *Repetition {
	this := &Repetition{}
	this.repeat = repeat
	this.element = element
	return this
}

func (this *Repetition) String() string {
	if this.repeat != nil {
		return this.repeat.String() + this.element.String()
	} else {
		return this.element.String()
	}
}

func (this *Repetition) GetDependentRuleNames() Set_RuleName {
	return this.element.GetDependentRuleNames()
}

func (this *Repetition) GetNFA(rules map[string]*Rule) *automata.NFA { //throws IllegalAbnfException {
	startState := automata.NewNFAState()
	acceptingState := automata.NewNFAState()
	this.GetNFAStates(startState, acceptingState, rules)
	return automata.NewNFA2(startState, acceptingState)
}

func (this *Repetition) GetNFAStates(startState, acceptingState *automata.NFAState, rules map[string]*Rule) {
	if this.repeat == nil {
		this.element.GetNFAStates(startState, acceptingState, rules)
		return
	}

	min := this.repeat.GetMin()
	max := this.repeat.GetMax()

	if min < 0 /*|| max < 0*/ {
		panic("Min value of a repeat element can not be less than zero.")
	}

	if max == -1 {
		if min == 0 {
			//              min == 0 && max == -1
			startState.AddTransitEpsilon(acceptingState)
			this.element.GetNFAStates(acceptingState, acceptingState, rules)
			return
		} else {
			//              min > 0 && max == -1
			current := startState
			for j := 0; j < min-1; j++ {
				next := automata.NewNFAState()
				this.element.GetNFAStates(current, next, rules)
				current = next
			}
			this.element.GetNFAStates(current, acceptingState, rules)
			this.element.GetNFAStates(acceptingState, acceptingState, rules)
			return
		}
	} else {
		if min == 0 {
			//              min == 0 && max > 0
			current := startState
			for j := 0; j < max-1; j++ {
				current.AddTransitEpsilon(acceptingState)
				next := automata.NewNFAState()
				this.element.GetNFAStates(current, next, rules)
				current = next
			}
			current.AddTransitEpsilon(acceptingState)
			this.element.GetNFAStates(current, acceptingState, rules)
			return
		} else if min == max {
			//              0 < min == max
			current := startState
			for j := 0; j < max-1; j++ {
				next := automata.NewNFAState()
				this.element.GetNFAStates(current, next, rules)
				current = next
			}
			this.element.GetNFAStates(current, acceptingState, rules)
			return
		} else if min < max {
			//              0 < min < max
			current := startState
			for j := 0; j < min; j++ {
				next := automata.NewNFAState()
				this.element.GetNFAStates(current, next, rules)
				current = next
			}
			for j := 0; j < max-min-1; j++ {
				current.AddTransitEpsilon(acceptingState)
				next := automata.NewNFAState()
				this.element.GetNFAStates(current, next, rules)
				current = next
			}
			current.AddTransitEpsilon(acceptingState)
			this.element.GetNFAStates(current, acceptingState, rules)
			return
		} else {
			panic("Max can not less than min")
		}
	}
}
