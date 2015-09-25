package abnf

import ()

type Rule struct {
	ruleName  *RuleName
	definedAs string
	elements  *Elements
}

func NewRule(ruleName *RuleName, definedAs string, elements *Elements) *Rule {
	this := &Rule{}
	this.ruleName = ruleName
	this.definedAs = definedAs
	this.elements = elements
	return this
}

func (this *Rule) GetRuleName() *RuleName {
	return this.ruleName
}

func (this *Rule) GetDefinedAs() string {
	return this.definedAs
}

func (this *Rule) SetDefinedAs(definedAs string) {
	this.definedAs = definedAs
}

func (this *Rule) GetElements() *Elements {
	return this.elements
}

func (this *Rule) String() string{
	return this.ruleName.String()+" "+this.definedAs+" "+this.elements.String();
}
