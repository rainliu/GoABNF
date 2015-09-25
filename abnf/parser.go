package abnf

import (
	"bytes"
	"container/list"
	"fmt"
	"io"
	"strconv"
)

//    ABNF文法解析器
type Parser struct {
	//prefix string

	//    ABNF文法解析器的输入流，这是一个支持peek和read操作的输入流，
	//    支持peek是因为这是一个预测解析器，即需要向前看1～2个字符，
	//    以决定下一步所需要匹配的ABNF文法产生式（或元素）。
	peeker *Peeker
}

//    构造函数，设置规则名的前缀和输入源，并将普通的输入源转化为支持peek操作的输入源。
func NewParser(reader io.Reader) *Parser {
	this := &Parser{}
	//this.prefix = prefix
	this.peeker = NewPeeker(reader)
	return this
}

//        调用parse函数开始对输入源进行解析，返回输入源中定义的ABNF规则列表
func (this *Parser) Parse() (*list.List, error) {
	return this.rulelist()
}

//    match函数用来判断两个字符是否相同
//    （例如判断输入的字符是否与期望的字符相同）
func (this *Parser) MatchExpected(value, expected int) bool {
	return value == expected
}

//    match函数用来判断字符是否在某个范围之内
//    （例如判断输入的字符是否是字母、或数字字符等）
func (this *Parser) MatchRange(value, lower, upper int) bool {
	return value >= lower && value <= upper
}

//    match函数用来判断字符是否与某个字符相同
//    （忽略大小写）

func (this *Parser) MatchExpectedIgnoreCase(value int, expected byte) bool {
	var lowers [2]byte
	lowers[0] = byte(value)
	lowers[1] = expected
	uppers := bytes.ToUpper(lowers[:])
	return uppers[0] == uppers[1] //Character.toUpperCase(value) == Character.toUpperCase(expected);
}

//    match函数用来判断字符是否与某些字符相同
//    （例如判断输入的字符是否为'-','+',或'%'）
func (this *Parser) MatchExpectedChars(value int, expected []int) bool {
	for index := 0; index < len(expected); index++ {
		if value == expected[index] {
			return true
		}
	}
	return false
}

//    如果不匹配则抛出MatchException异常
//    MatchException中包含了产生匹配异常的符号输入流中的行列位置，以及期待的字符。
func (this *Parser) AssertMatchExpected(value, expected int) {
	if !this.MatchExpected(value, expected) {
		panic(NewMatchException("'"+strconv.Itoa(int(expected))+"' ["+fmt.Sprintf("%02X", expected)+"]", int(value), this.peeker.GetPos(), this.peeker.GetLine()).String())
	}
}

//    如果字符不在某个范围之内则抛出MatchException异常
//    MatchException中包含了产生匹配异常的符号输入流中的行列位置，以及期待的字符。
func (this *Parser) AssertMatchRange(value, lower, upper int) {
	if !this.MatchRange(value, lower, upper) {
		panic(NewMatchException(
			"'"+strconv.Itoa(int(lower))+"'~'"+strconv.Itoa(int(upper))+"' "+
				"["+fmt.Sprintf("%02X", lower)+"~"+fmt.Sprintf("%02X", lower)+"]",
			int(value), this.peeker.GetPos(), this.peeker.GetLine()).String())
	}
}

//    如果不匹配（忽略大小写）则抛出MatchException异常
//    MatchException中包含了产生匹配异常的符号输入流中的行列位置，以及期待的字符。
func (this *Parser) AssertMatchExpectedIgnoreCase(value int, expected byte) {
	if !this.MatchExpectedIgnoreCase(value, expected) {
		panic(NewMatchException("'"+string(expected)+"' ["+fmt.Sprintf("%02X", expected)+"]", int(value), this.peeker.GetPos(), this.peeker.GetLine()).String())
	}
}

/////////////////////////////////////////////////////////////
//Begin Core Rules
/////////////////////////////////////////////////////////////

//ALPHA            =  0x41-5A / 0x61-7A
func (this *Parser) ALPHA() string {
	//              向前看一个字符
	if this.peeker.Peek(0) >= 0x41 && this.peeker.Peek(0) <= 0x5A {
		this.AssertMatchRange(this.peeker.Peek(0), 0x41, 0x5A)
	} else {
		this.AssertMatchRange(this.peeker.Peek(0), 0x61, 0x7A)
	}

	value := this.peeker.Read()
	return string(value)
}

//BIT			= "0" / "1"
func (this *Parser) BIT() string {
	this.AssertMatchRange(this.peeker.Peek(0), 0x30, 0x31)
	value := this.peeker.Read()
	//      返回空格的字符串值
	return string(value)
}

//  CHAR          =  %x01-7E
func (this *Parser) CHAR() string {
	this.AssertMatchRange(this.peeker.Peek(0), 0x01, 0x7E)
	value := this.peeker.Read()
	//      返回空格的字符串值
	return string(value)
}

//  CR             =  %x0D
func (this *Parser) CR() string {
	//      断言下一个字符为0x0D（否则抛出MatchException异常）
	this.AssertMatchExpected(this.peeker.Peek(0), 0x0D)
	value := this.peeker.Read()
	//      返回回车的字符串值
	return string(value)
}

//  CRLF           =  CR LF
func (this *Parser) CRLF() string {
	//      回车和换行符号，直接调用相应的CR、LF方法来进行解析就OK。
	return this.CR() + this.LF()
}

//CTL            =  0x00-1F / 0x7F
func (this *Parser) CTL() string {
	//              向前看一个字符
	if this.peeker.Peek(0) == 0x7F {
		this.AssertMatchExpected(this.peeker.Peek(0), 0x7F)
	} else {
		this.AssertMatchRange(this.peeker.Peek(0), 0x00, 0x1F)
	}

	value := this.peeker.Read()
	return string(value)
}

//  DIGIT          =  %x30-39
func (this *Parser) DIGIT() string {
	this.AssertMatchRange(this.peeker.Peek(0), 0x30, 0x39)
	value := this.peeker.Read()
	//      返回空格的字符串值
	return string(value)
}

//  DQUOTE          =  %x22
func (this *Parser) DQUOTE() string {
	this.AssertMatchExpected(this.peeker.Peek(0), 0x22)
	value := this.peeker.Read()
	//      返回空格的字符串值
	return string(value)
}

//HEXDIG            =  DIGIT/"A"/"B"/"C"/"D"/"E"/"F"
func (this *Parser) HEXDIG() string {
	//              向前看一个字符
	if this.peeker.Peek(0) >= 0x30 && this.peeker.Peek(0) <= 0x39 {
		this.AssertMatchRange(this.peeker.Peek(0), 0x30, 0x39)
	} else {
		this.AssertMatchRange(this.peeker.Peek(0), 0x41, 0x46)
	}

	value := this.peeker.Read()
	return string(value)
}

//  HTAB           =  %x09
func (this *Parser) HTAB() string {
	//      断言下一个字符为0x09（否则抛出MatchException异常）
	this.AssertMatchExpected(this.peeker.Peek(0), 0x09)
	value := this.peeker.Read()
	//      返回HTAB的字符串值
	return string(value)
}

//  LF             =  %x0A
func (this *Parser) LF() string {
	//      断言下一个字符为0x0A（否则抛出MatchException异常）
	this.AssertMatchExpected(this.peeker.Peek(0), 0x0A)
	value := this.peeker.Read()
	//      返回换行的字符串值
	return string(value)
}

// LWSP			= *(WSP / CRLF WSP) ?
func (this *Parser) LWSP() string {
	//this.AssertMatchRange(this.peeker.Peek(0), 0x00, 0xFF)
	//value := this.peeker.Read()
	//      返回空格的字符串值
	//return string(value)
	println("LWSP TODO")
	return ""
}

//  OCTET          =  %x00-FF
func (this *Parser) OCTET() string {
	this.AssertMatchRange(this.peeker.Peek(0), 0x00, 0xFF)
	value := this.peeker.Read()
	//      返回空格的字符串值
	return string(value)
}

//  SP             =  %x20
func (this *Parser) SP() string {
	//      断言下一个字符为0x20（否则抛出MatchException异常）
	this.AssertMatchExpected(this.peeker.Peek(0), 0x20)
	value := this.peeker.Read()
	//      返回空格的字符串值
	return string(value)
}

//  VCHAR          =  %x21-7E
func (this *Parser) VCHAR() string {
	this.AssertMatchRange(this.peeker.Peek(0), 0x21, 0x7E)
	value := this.peeker.Read()
	//      返回空格的字符串值
	return string(value)
}

//WSP            =  SP / HTAB
func (this *Parser) WSP() string {
	//              向前看一个字符
	switch this.peeker.Peek(0) {
	//              如果这个字符是0x20，则是SP（空格），调用SP()方法
	case 0x20:
		return this.SP()
		//              如果这个字符是0x09，则是HTAB（跳格），调用HTAB()方法
	case 0x09:
		return this.HTAB()
		//              否则抛出匹配异常MatchException
	default:
		panic(NewMatchException("[0x20, 0x09]", int(this.peeker.Peek(0)),
			this.peeker.GetPos(), this.peeker.GetLine()).String())
	}
}

/////////////////////////////////////////////////////////////
//End Core Rules
/////////////////////////////////////////////////////////////

//              rule           =  rulename defined-as elements c-nl
//      解析rule的方法
func (this *Parser) rule() *Rule {
	//              rule的第一个元素是rulename，首先调用rulename()方法，并记录解析到的规则名
	rulename := this.rulename()
	//println(rulename.String())
	//              rulename后面紧接着defined-as元素，调用相应的方法
	definedAs := this.defined_as()
	//println(definedAs)
	//              defined-as后面接着elements元素，调用elements()
	elements := this.elements()
	//println(elements.String())
	//              elements后面接着c-nl元素，调用之。
	this.c_nl()

	//              返回解析到的规则
	return NewRule(rulename, definedAs, elements)
}

//              c-nl           =  comment / CRLF
func (this *Parser) c_nl() string {
	//              向前看一个字符
	switch this.peeker.Peek(0) {
	//              如果是分号，则是注释，调用comment()方法进行解析
	case ';':
		return this.comment()
		//              如果是0x0D，则是回车，调用CRLF()方法进行解析
	case 0x0D:
		return this.CRLF()
		//              否则抛出异常
	case PEEKER_EOF:
		return ""
	default:
		panic(NewMatchException("[';', 0x0D]", int(this.peeker.Peek(0)), this.peeker.GetPos(), this.peeker.GetLine()).String())
	}
}

//              element        =  rulename / group / option /
//                                char-val / num-val / prose-val
func (this *Parser) element() Element {
	//      向前看一个字符，如果在0x41~0x5A或0x61~0x7A之间（即大小写英文字母），则是规则名，调用rulename()方法进行解析
	if this.MatchRange(this.peeker.Peek(0), 0x41, 0x5A) || this.MatchRange(this.peeker.Peek(0), 0x61, 0x7A) {
		return this.rulename()
	}

	//      否则再检查这个字符
	switch this.peeker.Peek(0) {
	//          如果是左括号，则是group，调用group()
	case '(':
		return this.group()
		//          如果是左方括号，则调用option()
	case '[':
		return this.option()
		//          如果是双引号，则调用char_var()
	case 0x22:
		return this.char_val()
		//          如果是百分号，则调用num_val()
	case '%':
		return this.num_val()
		//          如果是左尖括号（小于号），则调用prose_val()
	case '<':
		return this.prose_val()
		//          否则抛出匹配异常
	default:
		panic(NewMatchException("['(', '[', 0x22, '%', '<']", int(this.peeker.Peek(0)), this.peeker.GetPos(), this.peeker.GetLine()).String())
	}
}

//              char-val       =  DQUOTE *(%x20-21 / %x23-7E) DQUOTE
//  DQUOTE         =  %x22
func (this *Parser) char_val() *CharVal {
	var char_val bytes.Buffer
	//      char-val是双引号开始的
	this.AssertMatchExpected(this.peeker.Peek(0), 0x22)
	//      把这个双引号消化掉 :)
	this.peeker.Read()
	//      双引号后面跟着的0x20-21、0x23-7E都属于合法的字符，读入之
	for this.MatchRange(this.peeker.Peek(0), 0x20, 0x21) || this.MatchRange(this.peeker.Peek(0), 0x23, 0x7E) {
		char_val.WriteByte(byte(this.peeker.Read()))
	}
	//      如果不是跟着0x20-21、0x23-7E，则必须是双引号，否则异常
	this.AssertMatchExpected(this.peeker.Peek(0), 0x22)
	this.peeker.Read()
	//      返回这个字符串
	return NewCharVal(char_val.String())
}

//              prose-val      =  "<" *(%x20-3D / %x3F-7E) ">"
//      prose_val()方法和char_val()方法很类似，请自行阅读
func (this *Parser) prose_val() *ProseVal {
	var proseval bytes.Buffer
	this.AssertMatchExpected(this.peeker.Peek(0), '<')
	this.peeker.Read()
	for this.MatchRange(this.peeker.Peek(0), 0x20, 0x3D) || this.MatchRange(this.peeker.Peek(0), 0x3F, 0x7E) {
		proseval.WriteByte(byte(this.peeker.Read()))
	}
	this.AssertMatchExpected(this.peeker.Peek(0), '>')
	this.peeker.Read()
	return NewProseVal(proseval.String())
}

//    c-wsp          =  WSP / (c-nl WSP)
func (this *Parser) c_wsp() string {
	//              由于c-wsp可以派生为WSP或者c-nl WSP，当我们向前看字符输入流时，
	//              需要看这两个产生式的第一个字母，在龙书中也就时求FIRST(WSP)和
	//              FIRST(c-nl WSP)两个函数，其中：
	//              FIRST(WSP) = {0x20, 0x09};
	//              FIRST(c-nl WSP) = {';', 0x0D};
	//              我们开心的看到，FIRST(WSP)和FIRST(c-nl WSP)没有交集，
	//              因此只需要向前看一个字符就足够了。
	switch this.peeker.Peek(0) {
	case 0x20:
		fallthrough
	case 0x09: //println("i'm only wsp");
		return this.WSP()
	case ';':
		fallthrough
	case 0x0D: //println("i'm only c_nl+wsp");
		return this.c_nl() + this.WSP()
	default:
		panic(NewMatchException("[0x20, ';']", int(this.peeker.Peek(0)), this.peeker.GetPos(), this.peeker.GetLine()).String())
	}
}

//    comment        =  ";" *(WSP / VCHAR) CRLF
func (this *Parser) comment() string {
	var comment bytes.Buffer
	//      注释是分号开头
	this.AssertMatchExpected(this.peeker.Peek(0), ';')
	comment.WriteByte(byte(this.peeker.Read()))

	//      如果下一个字符是0x20, 0x09，或者0x21-0x7E之间，则进入循环
	//      直至输入字符不再这个范围之内。
	//      循环内是一个解析WSP或者VCHAR的过程，WSP和VCHAR的FIRST交集为空，
	//      因而是可以通过向前看一个字符区分开来的。
	//      另外WSP/VCHAR允许0个或任意多个，所以此处使用while是可以的。
	for this.MatchExpected(this.peeker.Peek(0), 0x20) || this.MatchExpected(this.peeker.Peek(0), 0x09) || this.MatchRange(this.peeker.Peek(0), 0x21, 0x7E) {
		if this.MatchExpected(this.peeker.Peek(0), 0x20) || this.MatchExpected(this.peeker.Peek(0), 0x09) {
			comment.WriteString(this.WSP())
		} else {
			comment.WriteString(this.VCHAR())
		}
		//          if (peekMatch ==0x20 || peekMatch == 0x09) WSP();
		//          else if (peekMatch >= 0x21 && peekMatch <= 0x7E) VCHAR();
	}
	//      结束之前要匹配回车换行字符
	comment.WriteString(this.CRLF())
	return comment.String()
}

//  解析各个进制
func (this *Parser) val(base byte, matcher Matcher) Element {
	//      检查进制符号
	this.AssertMatchExpectedIgnoreCase(this.peeker.Peek(0), base)
	baseValue := this.peeker.Read()
	var from bytes.Buffer  //= "";
	//var value bytes.Buffer //= "";

	//      进制符号之后的第一个字符，必须在Matcher定义的字符集内，否则异常
	if matcher.Match(this.peeker.Peek(0)) {
		//          连续读入符合字符的字符，构成NumVal的第一个数值。
		for matcher.Match(this.peeker.Peek(0)) {
			from.WriteByte(byte(this.peeker.Read()))
		}
		//println(from.String());
		//          第一个数值后面如果是跟着点号，则是一个数列NumVal，如果是－破折号，则是一个范围型数值RangedNumVal，如果都不是，则是单一个数值
		if this.MatchExpected(this.peeker.Peek(0), '.') {
			numval := NewNumVal(string(baseValue), false)
			//              将刚才匹配到的数值作为第一个数值加到将要返回的NumVal中
			numval.AddValue(from.String())
			//              如果后面跟着点号，则继续加入新的数值到NumVal中
			for this.MatchExpected(this.peeker.Peek(0), '.') {
				next := this.peeker.Peek(1)
				if !(matcher.Match(next)) {
					break
				}
				this.peeker.Read()
				//value = "";
				var value bytes.Buffer
				for matcher.Match(this.peeker.Peek(0)) {
					value.WriteByte(byte(this.peeker.Read()))
				}
				numval.AddValue(value.String())
				//println(value.String());
			}
			//              直到不能匹配到点号，数列结束，返回
			return numval
		} else if this.MatchExpected(this.peeker.Peek(0), '-') {
			//              这里向前读取两个字符，因此即使破折号后面跟着的不是数字，也能返回单一个数字而且将破折号留给后面的分析程序
			//              这是本程序里为数不多的能够具备回溯的代码段之一，嘿嘿。
			numval := NewNumVal(string(baseValue), true)
			numval.AddValue(from.String())

			next := this.peeker.Peek(1)
			if !(matcher.Match(next)) {
				//                  如果破折号后面跟的不是数字，则破折号不读入，返回单一数值
				numval := NewNumVal(string(baseValue), false)
				numval.AddValue(from.String())
				return numval
			} else {
				//              否则，破折号后面是数值，读取之，并返回RangedNumVal类型
				this.peeker.Read()
				var value bytes.Buffer
				value.WriteByte(byte(this.peeker.Read()))
				for matcher.Match(this.peeker.Peek(0)) {
					value.WriteByte(byte(this.peeker.Read()))
				}
				numval.AddValue(value.String())
				return numval
			}
		} else {
			//println("i am other");
			//              第一个数值之后跟的既不是点号，也不是破折号，则返回单一数值
			numval := NewNumVal(string(baseValue), false)
			numval.AddValue(from.String())
			return numval
		}
	} else {
		panic(NewMatchException(matcher.Expected(), int(this.peeker.Peek(0)), this.peeker.GetPos(), this.peeker.GetLine()).String())
	}

}

//              num-val        =  "%" (bin-val / dec-val / hex-val)
//      解析num-val
func (this *Parser) num_val() Element {
	//String base = "", from ="", val ="";
	//              百分号开头
	this.AssertMatchExpected(this.peeker.Peek(0), '%')
	this.peeker.Read()
	//              根据进制符号选择相应的解析方法（函数）
	switch this.peeker.Peek(0) {
	case 'b':
		fallthrough
	case 'B':
		var bin BinVal
		return this.val('b', &bin)
	case 'd':
		fallthrough
	case 'D':
		var dec DecVal
		return this.val('d', &dec)
	case 'x':
		fallthrough
	case 'X':
		var hex HexVal
		return this.val('x', &hex)
	default:
		panic(NewMatchException("['b', 'd', 'x']", int(this.peeker.Peek(0)), this.peeker.GetPos(), this.peeker.GetLine()).String())
	}
}

//              repetition     =  [repeat] element
//    DIGIT          =  %x30-39
func (this *Parser) repetition() *Repetition {
	var r *Repeat
	//      若以数字或者星号开头，则进入repeat
	if this.MatchRange(this.peeker.Peek(0), 0x30, 0x39) || this.MatchExpected(this.peeker.Peek(0), '*') {
		r = this.repeat()
	}
	//      element是必须的
	e := this.element()
	return NewRepetition(r, e)
}

//              repeat         =  1*DIGIT / (*DIGIT "*" *DIGIT)
func (this *Parser) repeat() *Repeat {
	min := 0
	max := -1 //infinite
	//      如果repeat是以星号开头，则重复的最小次数为0次，即repeat后面的element可以不出现。
	if this.MatchExpected(this.peeker.Peek(0), '*') {
		this.peeker.Read()
		//          如果星号后面有数字，则重复的最大次数是该数字所表示的次数，否则最大次数没有限制
		if this.MatchRange(this.peeker.Peek(0), 0x30, 0x39) {
			max = 0;
			for this.MatchRange(this.peeker.Peek(0), 0x30, 0x39) {
				i, _ := strconv.Atoi(string(this.peeker.Read()))
				max = max*10 + i
			}
		}
		return NewRepeat(min, max, true)
	} else if this.MatchRange(this.peeker.Peek(0), 0x30, 0x39) {
		//      repeat是以数字开头，其值表示重复的最小次数
		for this.MatchRange(this.peeker.Peek(0), 0x30, 0x39) {
			i, _ := strconv.Atoi(string(this.peeker.Read()))
			min = min*10 + i
		}
		//          如果有星号，则表示有范围
		if this.MatchExpected(this.peeker.Peek(0), '*') {
			this.peeker.Read()
			//              星号后面接着数字，表示重复的最大次数，否则最大次数没有限制
			if this.MatchRange(this.peeker.Peek(0), 0x30, 0x39) {
				max = 0;
				for this.MatchRange(this.peeker.Peek(0), 0x30, 0x39) {
					i, _ := strconv.Atoi(string(this.peeker.Read()))
					max = max*10 + i
				}
			}
			return NewRepeat(min, max, true)
		} else {
			//          没有星号，表示固定的重复次数
			return NewRepeat(min, max, false)
		}
	} else {
		panic(NewMatchException("['0'-'9', '*']", int(this.peeker.Peek(0)), this.peeker.GetPos(), this.peeker.GetLine()).String())
	}
}

//              alternation    =  concatenation
//                                *(*c-wsp "/" *c-wsp concatenation)
func (this *Parser) alternation() *Alternation {
	alternation := NewAlternation()
	//              每个alternation至少有一个候选项，这个候选项的类型是concatenation（连结项）
	alternation.AddConcatenation(this.concatenation())
	//      从第二个候选项开始，每个候选项都是都是以空格（可选）以及“/”引导的，
	//      因此，只要遇到空格或者/号，就认为接下来的又是一个候选项
	//      当然，如果遇到空格但后面跟的不是/号，又或者如果/号之后跟的不是候选项，
	//      那就只能异常了，因为这个算法不能回溯到空格或者/号之前
	var chars = []int{0x20, ';', '/'}
	for this.MatchExpectedChars(this.peeker.Peek(0), chars) {
		//          如遇到空格或者分号，则进入c_wsp()
		for this.MatchExpected(this.peeker.Peek(0), 0x20) || this.MatchExpected(this.peeker.Peek(0), ';') {
			this.c_wsp()
		}
		//          此处必须是/号了，否则异常，没有办法回溯
		this.AssertMatchExpected(this.peeker.Peek(0), '/')
		this.peeker.Read()
		//          /号后面可以跟若干空格或注释
		for this.MatchExpected(this.peeker.Peek(0), 0x20) || this.MatchExpected(this.peeker.Peek(0), ';') {
			this.c_wsp()
		}
		//          空格之后的新的候选项，候选项本身是concatenation，所以进入相应的函数。
		alternation.AddConcatenation(this.concatenation())
	}
	return alternation
}

//              concatenation  =  repetition *(1*c-wsp repetition)
func (this *Parser) concatenation() *Concatenation {
	concatenation := NewConcatenation()
	//              一个concatenation是由至少一个repetition组成的，
	//              这些repetition有先后顺序之分，用若干空格隔开
	concatenation.AddRepetition(this.repetition())
	//      后面有空格或分号，则认为会接着一个repetition
	//      其实这样是不严谨的，因为空格后面其实不必然是repetition，
	//      也可能是其他文法单位，但作为一个手工编写的解析器
	//      暂时接受它诸多的缺陷吧。
	for this.MatchExpected(this.peeker.Peek(0), 0x20) || this.MatchExpected(this.peeker.Peek(0), ';') {
		for this.MatchExpected(this.peeker.Peek(0), 0x20) || this.MatchExpected(this.peeker.Peek(0), ';') {
			this.c_wsp()
		}
		concatenation.AddRepetition(this.repetition())
	}
	return concatenation
}

//          group          =  "(" *c-wsp alternation *c-wsp ")"
func (this *Parser) group() *Group {
	//      一个group以左圆括号引导
	this.AssertMatchExpected(this.peeker.Peek(0), '(')
	this.peeker.Read()
	//      括号后面的若干空格
	var chars = []int{0x20, ';', 0x0D}
	for this.MatchExpectedChars(this.peeker.Peek(0), chars) {
		this.c_wsp()
	}
	//      一个group包含一个alternation
	alternation := this.alternation()
	for this.MatchExpectedChars(this.peeker.Peek(0), chars) {
		this.c_wsp()
	}
	//      以右圆括号结束
	this.AssertMatchExpected(this.peeker.Peek(0), ')')
	this.peeker.Read()
	return NewGroup(alternation)
}

//              option         =  "[" *c-wsp alternation *c-wsp "]"
//      option与group类似，差别在于是方括号而不是圆括号。
func (this *Parser) option() *Option {
	this.AssertMatchExpected(this.peeker.Peek(0), '[')
	this.peeker.Read()
	var chars = []int{0x20, ';', 0x0D}
	for this.MatchExpectedChars(this.peeker.Peek(0), chars) {
		this.c_wsp()
	}
	alternation := this.alternation()
	for this.MatchExpectedChars(this.peeker.Peek(0), chars) {
		this.c_wsp()
	}
	this.AssertMatchExpected(this.peeker.Peek(0), ']')
	this.peeker.Read()
	return NewOption(alternation)
}

//     rulelist       =  1*( rule / (*c-wsp c-nl) )
func (this *Parser) rulelist() (rl *list.List, err error) {
	defer func() {
		if r := recover(); r != nil {
			var ok bool
			err, ok = r.(error)
			if !ok {
				err = fmt.Errorf("%v", r)
			}
		}
	}()

	ruleMap := make(map[*RuleName]*Rule) //Map<RuleName, Rule> ruleMap = new HashMap<RuleName, Rule>();
	ruleList := list.New()               //List<Rule> ruleList = new ArrayList<Rule>();
	//      如果前向字符是字母、空格、分号、回车，则认为是rule、c-wsp或者c-nl
	for this.MatchRange(this.peeker.Peek(0), 0x41, 0x5A) || this.MatchRange(this.peeker.Peek(0), 0x61, 0x7A) || this.MatchExpected(this.peeker.Peek(0), 0x20) || this.MatchExpected(this.peeker.Peek(0), ';') || this.MatchExpected(this.peeker.Peek(0), 0x0D) {
		//          如果是字母开头，则认为是rule，否则是c-wsp或者c-nl
		if this.MatchRange(this.peeker.Peek(0), 0x41, 0x5A) || this.MatchRange(this.peeker.Peek(0), 0x61, 0x7A) {
			//              解析一条规则
			rule := this.rule()
			//              判断该条规则是否已经有有定义
			if defined, present := ruleMap[rule.GetRuleName()]; !present {
				//                  如果没有定义则放入规则列表
				ruleMap[rule.GetRuleName()] = rule
				ruleList.PushBack(rule)
			} else {
				//                  已有定义，则检查定义方式是否为增量定义
				//Rule defined := ruleMap.get(rule.getRuleName());
				if rule.GetDefinedAs() == "=" && defined.GetDefinedAs() == "=" {
					//                      如果不是增量定义，则抛出重复定义异常
					panic(NewCollisionException(rule.GetRuleName().String()+" is redefined.", this.peeker.GetPos(), this.peeker.GetLine()).String())
				}
				//                  如果是增量定义则合并两条规则
				if rule.GetDefinedAs() == "=" {
					defined.SetDefinedAs("=")
				}
				defined.GetElements().GetAlternation().GetConcatenations().PushBackList(rule.GetElements().GetAlternation().GetConcatenations())
			}
			//println(rule.String())
		} else {
			//println("begin *c-wsp c-nl")
			//              空格、分号、回车，则是c_wsp
			//for this.MatchExpected(this.peeker.Peek(0), 0x20) || this.MatchExpected(this.peeker.Peek(0), ';') || this.MatchExpected(this.peeker.Peek(0), 0x0D) {
			//	this.c_wsp()
			//}
			this.c_nl()
			//println("")
		}
	}
	return ruleList, nil
}

//              rulename       =  ALPHA *(ALPHA / DIGIT / "-")
func (this *Parser) rulename() *RuleName {
	//       ALPHA          =  %x41-5A / %x61-7A   ; A-Z / a-z
	//       DIGIT          =  %x30-39
	//      规则名的第一个字符必须是字母
	if !(this.MatchRange(this.peeker.Peek(0), 0x41, 0x5A) || this.MatchRange(this.peeker.Peek(0), 0x61, 0x7A)) {
		panic(NewMatchException("'A'-'Z'/'a'-'z'", int(this.peeker.Peek(0)), this.peeker.GetPos(), this.peeker.GetLine()).String())
	}
	var rulename bytes.Buffer //= "";
	rulename.WriteByte(byte(this.peeker.Read()))
	//      规则名的后续字符可以是字母、数字、破折号
	for this.MatchRange(this.peeker.Peek(0), 0x41, 0x5A) || this.MatchRange(this.peeker.Peek(0), 0x61, 0x7A) || this.MatchRange(this.peeker.Peek(0), 0x30, 0x39) || this.MatchExpected(this.peeker.Peek(0), '-') {
		rulename.WriteByte(byte(this.peeker.Read()))
	}
	return NewRuleName(rulename.String())
}

//              defined-as     =  *c-wsp ("=" / "=/") *c-wsp
func (this *Parser) defined_as() string {
	var value bytes.Buffer //= "";
	//      等号前面的空格
	for this.MatchExpected(this.peeker.Peek(0), 0x20) || this.MatchExpected(this.peeker.Peek(0), 0x09) || this.MatchExpected(this.peeker.Peek(0), ';') || this.MatchExpected(this.peeker.Peek(0), 0x0D) {
		this.c_wsp()
	}
	//      等号
	this.AssertMatchExpected(this.peeker.Peek(0), '=')
	value.WriteByte(byte(this.peeker.Read()))
	//      是否增量定义
	if this.MatchExpected(this.peeker.Peek(0), '/') {
		value.WriteByte(byte(this.peeker.Read()))
	}
	//      等号后面的空格
	for this.MatchExpected(this.peeker.Peek(0), 0x20) || this.MatchExpected(this.peeker.Peek(0), 0x09) || this.MatchExpected(this.peeker.Peek(0), ';') || this.MatchExpected(this.peeker.Peek(0), 0x0D) {
		this.c_wsp()
	}
	return value.String()
}

//              elements       =  alternation *c-wsp
func (this *Parser) elements() *Elements {
	//              元素elements其实就是alternation再接着若干空格
	alternation := this.alternation()
	for this.MatchExpected(this.peeker.Peek(0), 0x20) || this.MatchExpected(this.peeker.Peek(0), 0x09) || this.MatchExpected(this.peeker.Peek(0), ';') {
		this.c_wsp()
	}
	return NewElements(alternation)
}
