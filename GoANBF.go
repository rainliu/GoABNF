package main

import (
	"GoABNF/abnf"
	"GoABNF/automata"
	"container/list"
	"fmt"
	"os"
)

func checkRegularExpression(ruleList *list.List) bool {
	analyzer := abnf.NewRegularAnalyzer(ruleList)
	println("=====================Regular Expressions Begin=====================")
	for e := analyzer.GetRegularRules().Front(); e != nil; e = e.Next() {
		println(e.Value.(*abnf.Rule).String())
	}
	println("=====================Regular Expressions End=======================")
	println("=====================Nonregular Expressions Begin==================")
	for e := analyzer.GetNonRegularRules().Front(); e != nil; e = e.Next() {
		println(e.Value.(*abnf.Rule).String())
	}
	println("=====================Nonregular Expressions End====================")
	println("=====================Undefined Expressions Begin===================")
	for e := analyzer.GetUndefinedRules().Front(); e != nil; e = e.Next() {
		println(e.Value.(*abnf.Rule).String())
	}
	println("=====================Undefined Expressions End=====================")
	return analyzer.GetNonRegularRules().Len() == 0 && analyzer.GetUndefinedRules().Len() == 0
}

func GenerateNFA(ruleName string, regularRuleList *list.List) *automata.NFA {
	rules := make(map[string]*abnf.Rule)
	for e := regularRuleList.Front(); e != nil; e = e.Next() {
		v := e.Value.(*abnf.Rule)
		rules[v.GetRuleName().String()] = v
	}
 	startState := automata.NewNFAState();
    acceptingState := automata.NewNFAState();
    rules[ruleName].GetElements().GetNFAStates(startState, acceptingState, rules)
	return automata.NewNFA2(startState, acceptingState)
}

func main() {
	if len(os.Args) < 2 {
		println("Too few augments. Usage: GoABNF abnf.txt")
		return
	}
	f, err := os.Open(os.Args[1])
	if err != nil {
		println(err.Error())
		return
	}
	defer f.Close()

	p := abnf.NewParser(f)
	ruleList, err := p.Parse()
	if err != nil {
		println(err.Error())
		println("ruleList==nil")
		return
	}
	for e := ruleList.Front(); e != nil; e = e.Next() {
		v := e.Value.(*abnf.Rule)
		fmt.Printf("%s\n", v.String())
	}

	if !checkRegularExpression(ruleList) {
		println("Error: There are non-regular expressions.")
	}

	regularAnalyzer := abnf.NewRegularAnalyzer(ruleList)
	regularRuleList := regularAnalyzer.GetRegularRules()

	nfa := GenerateNFA("RFC3261-SIP-message", regularRuleList)
	//nfa.GetStartState().printToDot();
	fmt.Printf("Total states = %d\n", len(nfa.GetStateSet()))
	//nfa.getStartState().printToDot();
	//println("NFA print completed.");
}
