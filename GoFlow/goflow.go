package GoFlow

import (
	"runtime/debug"
	"sync"
)




type ICallable = func(_data *_Data)*_Result

type IBoolFunc = func(_data *_Data)bool

type IPrepareFunc = func(input _PrepareInput)*_Data

type INodeBeginLogger = func(msg string,_data *_Data)

type INodeEndLogger = func(msg string, _data *_Data, _result *_Result)

type NodeType int64

const(
	NotSpecific NodeType = iota
	NormalNodeType
	IfNodeType
	ElseNodeType
	ForNodeType
	ParallelNodeType
	ElseIfNodeType
)

type IBasicCheckNode interface {
	SetParentResult(result *_Result)
	GetParentResult() *_Result
	DoCheck()
	ImplCheck() *_Result
	SetNext(node IBasicCheckNode)
	GetNext()IBasicCheckNode
	GetNodeType() NodeType
	SetShouldSkip(shouldSkip bool)
}

//Checker Implementation

type Checker struct {

}

func (c *Checker)Prepare(prepareFunc IPrepareFunc, input _PrepareInput)*CheckEngine{
	data:= prepareFunc(input)
	return NewCheckEngine(data)
}

//END Checker


//Errors

type ConditionNotFoundError struct {}

func NewConditionNotFoundError()*ConditionNotFoundError{
	return &ConditionNotFoundError{}
}

func (c *ConditionNotFoundError)Error()string{
	return "condition is nil"
}


type PanicHappened struct {
	Msg string
}

func NewPanicHappened(bt string)*PanicHappened{
	return &PanicHappened{Msg: bt}
}

func (c *PanicHappened)Error()string{
	return c.Msg
}
//END Errors

// BasicCheckNode Implementation
type BasicCheckNode struct{
	Functors []ICallable
	NodeType NodeType
	Next IBasicCheckNode
	Data *_Data
	ShouldSkip bool
	parentResult **_Result

}

func NewBasicCheckNode(data *_Data,parentResult **_Result, nodeType NodeType, functors ...ICallable)*BasicCheckNode{
	return &BasicCheckNode{
		Functors:     functors,
		NodeType:     nodeType,
		Data:         data,
		parentResult: parentResult,
	}
}

func (b *BasicCheckNode)SetParentResult(result *_Result){
	*b.parentResult = result
}

func (b *BasicCheckNode)GetParentResult()*_Result{
	return *b.parentResult
}

func (b *BasicCheckNode)DoCheck(){
	if b.ShouldSkip || b.GetParentResult().Err!=nil ||b.GetParentResult().StatusCode!=0 {
		return
	}
	result := b.ImplCheck()
	b.SetParentResult(result)
}

func (b *BasicCheckNode)ImplCheck()*_Result{
	return &_Result{
		Err:        nil,
		StatusCode: 0,
		StatusMsg:  "",
	}
}

func (b *BasicCheckNode)GetNext()IBasicCheckNode{
	return b.Next
}

func (b *BasicCheckNode)SetNext(node IBasicCheckNode){
	b.Next = node
}

func (b *BasicCheckNode)GetNodeType()NodeType{
	return b.NodeType
}

func (b *BasicCheckNode)SetShouldSkip(shouldSkip bool){
	b.ShouldSkip = shouldSkip
}

//END BasicCheckNode


//IfNode Implementation
type IfNode struct{
	*BasicCheckNode
	Condition IBoolFunc
}

func NewIfNode(data *_Data, parentResult **_Result, condition IBoolFunc, functors ...ICallable)*IfNode{
	return &IfNode{
		BasicCheckNode: NewBasicCheckNode(data,parentResult,IfNodeType,functors...),
		Condition:      condition,
	}
}

func (i *IfNode)ImplCheck()*_Result{
	if i.Condition==nil{
		return &_Result{
			Err:  NewConditionNotFoundError(),
			StatusCode: 0,
			StatusMsg:  "",
		}
	}

	if i.Condition(i.Data){
		for _, functor := range i.Functors{
			result := functor(i.Data)
			if result.Err!=nil || result.StatusCode!=0{
				return result
			}
		}
	}

	current := i.Next
	for current!=nil && (current.GetNodeType() == ElseIfNodeType || current.GetNodeType() == ElseNodeType){
		current.SetShouldSkip(true)
		current = current.GetNext()
	}

	return i.GetParentResult()
}
//END IfNode

//ElseNode Implementation
type ElseNode struct{
	*BasicCheckNode
}

func NewElseNode(data *_Data, parentResult **_Result, functors ...ICallable)*ElseNode{
	return &ElseNode{NewBasicCheckNode(data,parentResult,ElseNodeType,functors...)}
}

func (e *ElseNode)ImplCheck()*_Result{
	for _, functor := range e.Functors{
		result := functor(e.Data)
		if result.Err != nil || result.StatusCode != 0{
			return result
		}
	}
	return e.GetParentResult()
}

//END ElseNode


// ElseIfNode Implementation
type ElseIfNode struct{
	*BasicCheckNode
	Condition IBoolFunc
}

func NewElseIfNode(data *_Data, parentResult **_Result, condition IBoolFunc, functors ...ICallable)*ElseIfNode{
	return &ElseIfNode{
		BasicCheckNode: NewBasicCheckNode(data,parentResult,ElseIfNodeType,functors...),
		Condition:      condition,
	}
}

func (i *ElseIfNode)ImplCheck()*_Result{
	if i.Condition==nil{
		return &_Result{
			Err:  NewConditionNotFoundError(),
			StatusCode: 0,
			StatusMsg:  "",
		}
	}

	if i.Condition(i.Data){
		for _, functor := range i.Functors{
			result := functor(i.Data)
			if result.Err!=nil || result.StatusCode!=0{
				return result
			}
		}
	}

	current := i.Next
	for current!=nil && (current.GetNodeType() == ElseIfNodeType || current.GetNodeType() == ElseNodeType){
		current.SetShouldSkip(true)
		current = current.GetNext()
	}

	return i.GetParentResult()
}
//END ElseIfNode

//NormalNode Implementation
type NormalNode struct{
	*BasicCheckNode
}

func NewNormalNode(data *_Data, parentResult **_Result, functors ...ICallable)*ElseNode{
	return &ElseNode{NewBasicCheckNode(data,parentResult,NormalNodeType,functors...)}
}

func (e *NormalNode)ImplCheck()*_Result{
	for _, functor := range e.Functors{
		result := functor(e.Data)
		if result.Err != nil || result.StatusCode != 0{
			return result
		}
	}
	return e.GetParentResult()
}

//END NormalNode

//ForNode Implementation
type ForNode struct{
	*BasicCheckNode
	Times int
}

func NewForNode(times int,data *_Data, parentResult **_Result, functors ...ICallable)*ForNode{
	return &ForNode{
		BasicCheckNode: NewBasicCheckNode(data,parentResult,ForNodeType,functors...),
		Times: times,
	}
}

func (e *ForNode)ImplCheck()*_Result{
	for i:=0;i<e.Times;i++{
		for _, functor := range e.Functors{
			result := functor(e.Data)
			if result.Err != nil || result.StatusCode != 0{
				return result
			}
		}
	}
	return e.GetParentResult()
}

//END NormalNode

//ParallelNode Implementation
type ParallelNode struct{
	*BasicCheckNode
	Times int
}

func NewParallelNode(data *_Data, parentResult **_Result, functors ...ICallable)*ParallelNode{
	return &ParallelNode{
		BasicCheckNode: NewBasicCheckNode(data,parentResult,ParallelNodeType,functors...),
	}
}

func (e *ParallelNode)ImplCheck()*_Result{
	resultChan := make(chan *_Result,len(e.Functors))

	wg := sync.WaitGroup{}
	wg.Add(len(e.Functors))

	go func(wg *sync.WaitGroup) {
		wg.Wait()
		close(resultChan)
	}(&wg)

	for _, functor := range e.Functors{
		go func(wg *sync.WaitGroup, f ICallable) {
			defer func() {
				wg.Done()
				if a:=recover();a!=nil{
					debug.PrintStack()
					resultChan<-&_Result{
						Err:        NewPanicHappened(""),
						StatusCode: 0,
						StatusMsg:  "",
					}
				}
			}()
			result := f(e.Data)
			resultChan<-result
		}(&wg,functor)
	}

	result := e.GetParentResult()
	for item := range resultChan{
		if result.StatusCode!=0 || result.Err!=nil{
			continue
		}
		result = item
	}

	return result
}

//END NormalNode

//CheckEngine Implementation

type CheckEngine struct{
	data *_Data
	nodes []IBasicCheckNode
	result **_Result
}

func NewCheckEngine(data *_Data)*CheckEngine{
	res := &CheckEngine{
		data:   data,
		nodes:  make([]IBasicCheckNode,0,10),
	}
	tempResult := &_Result{
		Err:        nil,
		StatusCode: 0,
		StatusMsg:  "",
	}
	res.result = &tempResult
	return res
}

func (c *CheckEngine)Check(functors ...ICallable)*CheckEngine{
	node := NewNormalNode(c.data,c.result,functors...)
	if len(c.nodes) != 0{
		c.nodes[len(c.nodes)-1] = node
	}
	c.nodes = append(c.nodes,node)
	return c
}

func (c *CheckEngine)For(times int, functors ...ICallable)*CheckEngine{
	node := NewForNode(times,c.data,c.result,functors...)
	if len(c.nodes) != 0{
		c.nodes[len(c.nodes)-1] = node
	}
	c.nodes = append(c.nodes,node)
	return c
}

func (c *CheckEngine)Parallel(functors ...ICallable)*CheckEngine{
	node := NewParallelNode(c.data,c.result,functors...)
	if len(c.nodes) != 0{
		c.nodes[len(c.nodes)-1] = node
	}
	c.nodes = append(c.nodes,node)
	return c
}

func (c *CheckEngine)If(condition IBoolFunc, functors ...ICallable)*

//END CheckEngine

//ElseCheckEngine implementation

type ElseCheckEngine struct{
	data *_Data
	nodes *[]IBasicCheckNode
	result **_Result
	invoker *CheckEngine
}

func NewElseCheckEngine(data *_Data, invoker *CheckEngine,result **_Result,nodes *[]IBasicCheckNode)*ElseCheckEngine{
	res := &ElseCheckEngine{
		data:   data,
		nodes:  nodes,
		result: result,
		invoker: invoker,
	}
	return res
}

func (c *ElseCheckEngine)Check(functors ...ICallable)*ElseCheckEngine{
	node := NewNormalNode(c.data,c.result,functors...)
	if len(*c.nodes) != 0{
		(*c.nodes)[len(*c.nodes)-1] = node
	}
	*c.nodes = append(*c.nodes,node)
	return c
}

func (c *ElseCheckEngine)For(times int, functors ...ICallable)*ElseCheckEngine{
	node := NewForNode(times,c.data,c.result,functors...)
	if len(*c.nodes) != 0{
		(*c.nodes)[len(*c.nodes)-1] = node
	}
	*c.nodes = append(*c.nodes,node)
	return c
}

func (c *ElseCheckEngine)Parallel(functors ...ICallable)*ElseCheckEngine{
	node := NewParallelNode(c.data,c.result,functors...)
	if len(*c.nodes) != 0{
		(*c.nodes)[len(*c.nodes)-1] = node
	}
	*c.nodes = append(*c.nodes,node)
	return c
}

func (c *ElseCheckEngine)If(condition IBoolFunc, functors ...ICallable)*ElseCheckEngine{
	node := NewIfNode(c.data,c.result,condition,functors...)
	if len(*c.nodes) != 0{
		(*c.nodes)[len(*c.nodes)-1] = node
	}
	*c.nodes = append(*c.nodes,node)
	return c
}
//END ElseCheckEngine