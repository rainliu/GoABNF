package abnf

import (
	"container/list"
)

type RegularAnalyzer struct {
	nonRegularRules *list.List // new ArrayList<Rule>();
	regularRules    *list.List //= new ArrayList<Rule>();
	undefinedRules  *list.List //= new ArrayList<Rule>();
}

func NewRegularAnalyzer(rules *list.List) *RegularAnalyzer {
	this := &RegularAnalyzer{}

	this.nonRegularRules = list.New() // new ArrayList<Rule>();
	this.regularRules = list.New()    //= new ArrayList<Rule>();
	this.undefinedRules = list.New()  //= new ArrayList<Rule>();

	definedRuleNames := make(Set_RuleName)      // new HashSet<RuleName>();
	observedRules := make([]*Rule, rules.Len()) //new ArrayList<Rule>();
	for i, e := 0, rules.Front(); e != nil; i, e = i+1, e.Next() {
		observedRules[i] = e.Value.(*Rule)
		//println(observedRules[i].GetRuleName().String());
	}

	foundRegular := true
	for foundRegular {
		foundRegular = false
		var index int
		//println("=======");
		for index = len(observedRules) - 1; index >= 0; index-- {
			if observedRules[index] == nil {
				continue
			}
			//println(observedRules[index].GetElements().String());
			dependent := observedRules[index].GetElements().GetDependentRuleNames()

			if this.ContainsAll(definedRuleNames, dependent) {
				definedRuleNames[observedRules[index].GetRuleName().String()] = observedRules[index].GetRuleName()
				this.regularRules.PushBack(observedRules[index])
				observedRules[index] = nil //.remove(index);
				foundRegular = true
				continue
			}
			
			if _, present := dependent[observedRules[index].GetRuleName().String()]; !present {
				continue
			}

			delete(dependent, observedRules[index].GetRuleName().String())
			if this.ContainsAll(definedRuleNames, dependent) {
				definedRuleNames[observedRules[index].GetRuleName().String()] = observedRules[index].GetRuleName()
				this.nonRegularRules.PushBack(observedRules[index])
				observedRules[index] = nil //.remove(index);
				foundRegular = true
			}
		}
	}
	for index := 0; index < len(observedRules); index++ {
		if observedRules[index] != nil {
			this.undefinedRules.PushBack(observedRules[index])
		}
	}
	//observedRules.clear();

	return this
}

func (this *RegularAnalyzer) ContainsAll(definedRuleNames Set_RuleName, dependent Set_RuleName) bool {
	for _, r := range dependent {
		_, present := definedRuleNames[r.String()]
		if !present {
			return false
		}
	}
	return true
}

func (this *RegularAnalyzer) GetNonRegularRules() *list.List { return this.nonRegularRules }

func (this *RegularAnalyzer) GetRegularRules() *list.List { return this.regularRules }

func (this *RegularAnalyzer) GetUndefinedRules() *list.List { return this.undefinedRules }
