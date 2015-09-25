package abnf

import (
	"GoABNF/automata"
	"bytes"
	"container/list"
	"strconv"
)

/*
   num-val        =  "%" (bin-val / dec-val / hex-val)

   bin-val        =  "b" 1*BIT
                     [ 1*("." 1*BIT) / ("-" 1*BIT) ]
                          ; series of concatenated bit values
                          ; or single ONEOF range

   dec-val        =  "d" 1*DIGIT
                     [ 1*("." 1*DIGIT) / ("-" 1*DIGIT) ]

   hex-val        =  "x" 1*HEXDIG
                     [ 1*("." 1*HEXDIG) / ("-" 1*HEXDIG) ]
*/
type NumVal struct { //implements Element {//, Terminal {
	base   string
	ranged bool
	values *list.List //private List<String> values = new ArrayList<String>();
}

func NewNumVal(base string, ranged bool) *NumVal {
	this := &NumVal{}
	this.base = base
	this.ranged = ranged
	this.values = list.New()
	return this
}

func (this *NumVal) AddValue(value string) {
	this.values.PushBack(value)
}

func (this *NumVal) String() string {
	var s bytes.Buffer
	s.WriteString("%")
	s.WriteString(this.base)
	for e := this.values.Front(); e != nil; e = e.Next() {
		v := e.Value.(string)
		s.WriteString(v)
		if e.Next() != nil {
			if this.ranged {
				s.WriteString("-")
			} else {
				s.WriteString(".")
			}
		}
	}
	return s.String()
}

func (this *NumVal) GetElementType() ElementType {
	return ELEMENT_NUMVAL
}

func (this *NumVal) GetDependentRuleNames() Set_RuleName {
	return make(Set_RuleName)
}

func (this *NumVal) GetNFA(rules map[string]*Rule) *automata.NFA { //throws IllegalAbnfException {
	startState := automata.NewNFAState()
	acceptingState := automata.NewNFAState()
	this.GetNFAStates(startState, acceptingState, rules)
	return automata.NewNFA2(startState, acceptingState)
}

//@Override
func (this *NumVal) GetNFAStates(startState, acceptingState *automata.NFAState, rules map[string]*Rule) {
	if this.values.Len() == 0 {
		startState.AddTransitEpsilon(acceptingState)
		return
	}

	var radix int
	if this.base == "B" || this.base == "b" {
		radix = 2
	} else if this.base == "D" || this.base == "d" {
		radix = 10
	} else if this.base == "X" || this.base == "x" {
		radix = 16
	} else {
		panic("NumVal base can not be handled.")
	}

	current := startState
	e := this.values.Front()
	for j := 0; j < this.values.Len()-1; j++ {
		v := e.Value.(string)
		i, _ := strconv.ParseInt(v, radix, 64)
		current = current.AddTransitInt1(int(i))
		e = e.Next()
	}
	v := e.Value.(string)
	i, _ := strconv.ParseInt(v, radix, 64)
	current.AddTransitInt2(int(i), acceptingState)
}

type Matcher interface {
	Match(value int) bool
	Expected() string
}

type BinVal NumVal

func (this *BinVal) Match(value int) bool {
	//              如果符号是0或1就匹配
	return value == '0' || value == '1'
}
func (this *BinVal) Expected() string {
	//              提示符号不在符号集内（仅用于异常情况）
	return "['0', '1']"
}

type DecVal NumVal

func (this *DecVal) Match(value int) bool {
	return (value >= 0x30 && value <= 0x39)
}
func (this *DecVal) Expected() string {
	return "['0'-'9']"

}

type HexVal NumVal

func (this *HexVal) Match(value int) bool {
	return (value >= 0x30 && value <= 0x39) || (value >= 'A' && value <= 'F') || (value >= 'a' && value <= 'f')
}
func (this *HexVal) Expected() string {
	return "['0'-'9', 'A'-'F', 'a'-'f']"
}

func (this *NumVal) GetValues() *list.List {
	return this.values
}
