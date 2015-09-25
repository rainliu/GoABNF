package abnf

import (
	"GoABNF/automata"
)

type ElementType	int

const (
	ELEMENT_RULENAME ElementType= iota
	ELEMENT_GROUP
	ELEMENT_OPTION
	ELEMENT_CHARVAL
	ELEMENT_NUMVAL
	ELEMENT_PROSEVAL
)

type Element interface{
	String() string;
	GetElementType() ElementType;
	GetDependentRuleNames() Set_RuleName
	GetNFA(rules map[string]*Rule) *automata.NFA
    GetNFAStates(startState, acceptingState *automata.NFAState, rules map[string]*Rule)
};