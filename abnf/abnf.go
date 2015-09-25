package abnf

import (
	"GoABNF/automata"
)

type Abnf interface{
    GetDependentRuleNames() Set_RuleName;
    GetNFA(rules map[string]*Rule) *automata.NFA
    GetNFAStates(startState, acceptingState *automata.NFAState, rules map[string]*Rule)
}
