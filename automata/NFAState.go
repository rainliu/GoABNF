package automata

import (
	"bytes"
	//"container/list"
)

var NFAState_COUNT = 0

type Set_NFAState map[*NFAState]*NFAState

type NFAState struct { //implements Comparable<NFAState> {
	//状态标识，每个NFA状态节点都有唯一的数值标识
	id                int
	transitions       map[int]Set_NFAState
	epsilonTransition Set_NFAState
}

//在创建NFA状态对象的时候，通过静态变量生成唯一标识
func NewNFAState() *NFAState {
	this := &NFAState{}
	this.transitions = make(map[int]Set_NFAState)
	this.epsilonTransition = make(Set_NFAState)
	this.id = NFAState_COUNT
	NFAState_COUNT++
	return this
}

func (this *NFAState) GetId() int { return this.id }

//迁移函数，由于迁移函数需要两个输入：当前状态和输入符号，因此在一个状态对象内部，
//迁移函数都是针对本对象的，只需要输入符号就可以了，这里通过Map接口实现迁移函数
//protected Map<Integer, Set<NFAState>> transition = new HashMap<Integer, Set<NFAState>>();
func (this *NFAState) GetTransitions() map[int]Set_NFAState { return this.transitions }

//空字符迁移函数，即从当前节点经过空字符输入所能够到达的下一个状态节点
//protected Set<NFAState> epsilonTransition = new HashSet<NFAState>();
func (this *NFAState) GetEpsilonTransition() Set_NFAState { return this.epsilonTransition }

//向迁移函数添加一个映射，不给定下一个状态节点
func (this *NFAState) AddTransitInt1(input int) *NFAState {
	return this.AddTransitInt2(input, NewNFAState())
}

//向迁移函数添加一个映射，给定下一个状态节点
func (this *NFAState) AddTransitInt2(input int, next *NFAState) *NFAState {
	states, present := this.transitions[input]
	if !present {
		states = make(Set_NFAState) //new HashSet<NFAState>();
		this.transitions[input] = states
	}
	states[next] = next
	return next
}

//向迁移函数添加一个映射，不给定下一个状态节点
func (this *NFAState) AddTransitByte1(input byte) *NFAState {
	return this.AddTransitInt2(int(input), NewNFAState())
}

//向迁移函数添加一个映射，给定下一个状态节点
//假定我们的上下文无关文法是大小写不敏感的，当输入字符是char类型并且是字母时，
//生成大写字母和小写字母两个映射
func (this *NFAState) AddTransitByte2(input byte, next *NFAState) *NFAState {
	if input >= 'a' && input <= 'z' && input >= 'A' && input <= 'Z' {
		var b [1]byte
		b[0] = input
		this.AddTransitInt2(int(bytes.ToUpper(b[:])[0]), next)
		this.AddTransitInt2(int(bytes.ToLower(b[:])[0]), next)
		return next
	} else {
		this.AddTransitInt2(int(input), next)
		return next
	}
}

//添加一个空字符的映射
func (this *NFAState) AddTransitEpsilon(next *NFAState) *NFAState {
	this.epsilonTransition[next] = next
	return next
}

//返回迁移函数
func (this *NFAState) GetTransition(input int) Set_NFAState {
	return this.transitions[input]
}

func (this *NFAState) GetNextStates() Set_NFAState {
	allstates := make(Set_NFAState) // new HashSet<NFAState>();

	for _, states := range this.transitions {
		for _, state := range states {
			allstates[state] = state
		}
	}

	for _, state := range this.epsilonTransition {
		allstates[state] = state
	}

	return allstates
}
