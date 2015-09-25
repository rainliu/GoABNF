package abnf

import (
	"GoABNF/automata"
)

type RuleName struct { //implements Element {
	rulename string
}

type Set_RuleName map[string]*RuleName

func NewRuleName(rulename string) *RuleName {
	this := &RuleName{}
	this.rulename = rulename
	return this
}

func (this *RuleName) String() string {
	return this.rulename
}

func (this *RuleName) GetElementType() ElementType {
	return ELEMENT_RULENAME
}

func (this *RuleName) GetDependentRuleNames() Set_RuleName {
	ruleNames := make(Set_RuleName)
	ruleNames[this.rulename] = this
	return ruleNames
}

func (this *RuleName) GetNFA(rules map[string]*Rule) *automata.NFA { //throws IllegalAbnfException {
	startState := automata.NewNFAState()
	acceptingState := automata.NewNFAState()
	this.GetNFAStates(startState, acceptingState, rules)
	return automata.NewNFA2(startState, acceptingState)
}

func (this *RuleName) GetNFAStates(startState, acceptingState *automata.NFAState, rules map[string]*Rule) {
	if _, present := rules[this.String()]; !present {
		panic("Fail to find the definition of " + this.String())
	}

	rule := rules[this.String()]
	if rule.GetDefinedAs() == "=/" {
		panic("Can not handle incremental definition while generating NFA.")
	}

	rule.GetElements().GetNFAStates(startState, acceptingState, rules)
}
