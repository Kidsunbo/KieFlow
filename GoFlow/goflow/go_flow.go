package goflow

import (
	"runtime/debug"
	"sync"
)

type ICallable = func(_data *_Data) *_Result

type IBoolFunc = func(_data *_Data) bool

type IPrepareFunc = func(input _PrepareInput) *_Data

type INodeBeginLogger = func(note string, _data *_Data)

type INodeEndLogger = func(note string, _data *_Data, _result *_Result)

type IOnSuccessFunc = func(_data *_Data, _result *_Result)

type IOnFailFunc = func(_data *_Data, _result *_Result)

type NodeType int64

const (
	NormalNodeType NodeType = iota
	IfNodeType
	ElseNodeType
	ForNodeType
	ParallelNodeType
	ElseIfNodeType
)

type IBasicFlowNode interface {
	SetParentResult(result *_Result)
	GetParentResult() *_Result
	Run()
	ImplTask() *_Result
	SetNext(node IBasicFlowNode)
	GetNext() IBasicFlowNode
	GetNodeType() NodeType
	SetShouldSkip(shouldSkip bool)
	SetNote(note string)
	GetNote() string
	SetBeginLogger(logger INodeBeginLogger)
	GetBeginLogger() INodeBeginLogger
	SetEndLogger(logger INodeEndLogger)
	GetEndLogger() INodeEndLogger
}

//Flow Implementation

type Flow struct {
}

func (c *Flow) Prepare(prepareFunc IPrepareFunc, input _PrepareInput) *FlowEngine {
	data := prepareFunc(input)
	return NewFlowEngine(data)
}

//END Flow

//Errors

type ConditionNotFoundError struct{}

func NewConditionNotFoundError() *ConditionNotFoundError {
	return &ConditionNotFoundError{}
}

func (c *ConditionNotFoundError) Error() string {
	return "condition is nil"
}

type PanicHappened struct {
	Msg string
}

func NewPanicHappened(bt string) *PanicHappened {
	return &PanicHappened{Msg: bt}
}

func (c *PanicHappened) Error() string {
	return c.Msg
}

//END Errors

// BasicFlowNode Implementation
type BasicFlowNode struct {
	Functors     []ICallable
	NodeType     NodeType
	Next         IBasicFlowNode
	Data         *_Data
	ShouldSkip   bool
	parentResult **_Result
	BeginLogger  INodeBeginLogger
	EndLogger    INodeEndLogger
	Note         string
}

func NewBasicFlowNode(data *_Data, parentResult **_Result, nodeType NodeType, functors ...ICallable) *BasicFlowNode {
	return &BasicFlowNode{
		Functors:     functors,
		NodeType:     nodeType,
		Data:         data,
		parentResult: parentResult,
	}
}

func (b *BasicFlowNode) SetParentResult(result *_Result) {
	*b.parentResult = result
}

func (b *BasicFlowNode) GetParentResult() *_Result {
	return *b.parentResult
}

func (b *BasicFlowNode) Run() {
	if b.ShouldSkip || b.GetParentResult().Err != nil || b.GetParentResult().StatusCode != 0 {
		return
	}
	if b.BeginLogger != nil {
		b.BeginLogger(b.Note, b.Data)
	}

	result := b.ImplTask()
	b.SetParentResult(result)

	if b.EndLogger != nil {
		b.EndLogger(b.Note, b.Data, b.GetParentResult())
	}
}

func (b *BasicFlowNode) ImplTask() *_Result {
	return &_Result{
		Err:        nil,
		StatusCode: 0,
		StatusMsg:  "",
	}
}

func (b *BasicFlowNode) GetNext() IBasicFlowNode {
	return b.Next
}

func (b *BasicFlowNode) SetNext(node IBasicFlowNode) {
	b.Next = node
}

func (b *BasicFlowNode) GetNodeType() NodeType {
	return b.NodeType
}

func (b *BasicFlowNode) SetShouldSkip(shouldSkip bool) {
	b.ShouldSkip = shouldSkip
}

func (b *BasicFlowNode) SetNote(note string) {
	b.Note = note
}

func (b *BasicFlowNode) GetNote() string {
	return b.Note
}

func (b *BasicFlowNode) SetBeginLogger(logger INodeBeginLogger) {
	b.BeginLogger = logger
}

func (b *BasicFlowNode) GetBeginLogger() INodeBeginLogger {
	return b.BeginLogger
}

func (b *BasicFlowNode) SetEndLogger(logger INodeEndLogger) {
	b.EndLogger = logger
}

func (b *BasicFlowNode) GetEndLogger() INodeEndLogger {
	return b.EndLogger
}

//END BasicFlowNode

//IfNode Implementation
type IfNode struct {
	*BasicFlowNode
	Condition IBoolFunc
}

func NewIfNode(data *_Data, parentResult **_Result, condition IBoolFunc, functors ...ICallable) *IfNode {
	return &IfNode{
		BasicFlowNode: NewBasicFlowNode(data, parentResult, IfNodeType, functors...),
		Condition:     condition,
	}
}

func (i *IfNode) ImplTask() *_Result {
	if i.Condition == nil {
		return &_Result{
			Err:        NewConditionNotFoundError(),
			StatusCode: 0,
			StatusMsg:  "",
		}
	}

	if i.Condition(i.Data) {
		for _, functor := range i.Functors {
			result := functor(i.Data)
			if result.Err != nil || result.StatusCode != 0 {
				return result
			}
		}
	}

	current := i.Next
	for current != nil && (current.GetNodeType() == ElseIfNodeType || current.GetNodeType() == ElseNodeType) {
		current.SetShouldSkip(true)
		current = current.GetNext()
	}

	return i.GetParentResult()
}

//END IfNode

//ElseNode Implementation
type ElseNode struct {
	*BasicFlowNode
}

func NewElseNode(data *_Data, parentResult **_Result, functors ...ICallable) *ElseNode {
	return &ElseNode{NewBasicFlowNode(data, parentResult, ElseNodeType, functors...)}
}

func (e *ElseNode) ImplTask() *_Result {
	for _, functor := range e.Functors {
		result := functor(e.Data)
		if result.Err != nil || result.StatusCode != 0 {
			return result
		}
	}
	return e.GetParentResult()
}

//END ElseNode

// ElseIfNode Implementation
type ElseIfNode struct {
	*BasicFlowNode
	Condition IBoolFunc
}

func NewElseIfNode(data *_Data, parentResult **_Result, condition IBoolFunc, functors ...ICallable) *ElseIfNode {
	return &ElseIfNode{
		BasicFlowNode: NewBasicFlowNode(data, parentResult, ElseIfNodeType, functors...),
		Condition:     condition,
	}
}

func (i *ElseIfNode) ImplTask() *_Result {
	if i.Condition == nil {
		return &_Result{
			Err:        NewConditionNotFoundError(),
			StatusCode: 0,
			StatusMsg:  "",
		}
	}

	if i.Condition(i.Data) {
		for _, functor := range i.Functors {
			result := functor(i.Data)
			if result.Err != nil || result.StatusCode != 0 {
				return result
			}
		}
	}

	current := i.Next
	for current != nil && (current.GetNodeType() == ElseIfNodeType || current.GetNodeType() == ElseNodeType) {
		current.SetShouldSkip(true)
		current = current.GetNext()
	}

	return i.GetParentResult()
}

//END ElseIfNode

//NormalNode Implementation
type NormalNode struct {
	*BasicFlowNode
}

func NewNormalNode(data *_Data, parentResult **_Result, functors ...ICallable) *ElseNode {
	return &ElseNode{NewBasicFlowNode(data, parentResult, NormalNodeType, functors...)}
}

func (e *NormalNode) ImplTask() *_Result {
	for _, functor := range e.Functors {
		result := functor(e.Data)
		if result.Err != nil || result.StatusCode != 0 {
			return result
		}
	}
	return e.GetParentResult()
}

//END NormalNode

//ForNode Implementation
type ForNode struct {
	*BasicFlowNode
	Times int
}

func NewForNode(times int, data *_Data, parentResult **_Result, functors ...ICallable) *ForNode {
	return &ForNode{
		BasicFlowNode: NewBasicFlowNode(data, parentResult, ForNodeType, functors...),
		Times:         times,
	}
}

func (e *ForNode) ImplTask() *_Result {
	for i := 0; i < e.Times; i++ {
		for _, functor := range e.Functors {
			result := functor(e.Data)
			if result.Err != nil || result.StatusCode != 0 {
				return result
			}
		}
	}
	return e.GetParentResult()
}

//END NormalNode

//ParallelNode Implementation
type ParallelNode struct {
	*BasicFlowNode
	Times int
}

func NewParallelNode(data *_Data, parentResult **_Result, functors ...ICallable) *ParallelNode {
	return &ParallelNode{
		BasicFlowNode: NewBasicFlowNode(data, parentResult, ParallelNodeType, functors...),
	}
}

func (e *ParallelNode) ImplTask() *_Result {
	resultChan := make(chan *_Result, len(e.Functors))

	wg := sync.WaitGroup{}
	wg.Add(len(e.Functors))

	go func(wg *sync.WaitGroup) {
		wg.Wait()
		close(resultChan)
	}(&wg)

	for _, functor := range e.Functors {
		go func(wg *sync.WaitGroup, f ICallable) {
			defer func() {
				wg.Done()
				if a := recover(); a != nil {
					debug.PrintStack()
					resultChan <- &_Result{
						Err:        NewPanicHappened(""),
						StatusCode: 0,
						StatusMsg:  "",
					}
				}
			}()
			result := f(e.Data)
			resultChan <- result
		}(&wg, functor)
	}

	result := e.GetParentResult()
	for item := range resultChan {
		if result.StatusCode != 0 || result.Err != nil {
			continue
		}
		result = item
	}

	return result
}

//END NormalNode

//FlowEngine Implementation

type FlowEngine struct {
	data          *_Data
	nodes         []IBasicFlowNode
	result        **_Result
	onFailFunc    IOnFailFunc
	onSuccessFunc IOnSuccessFunc
}

func NewFlowEngine(data *_Data) *FlowEngine {
	res := &FlowEngine{
		data:  data,
		nodes: make([]IBasicFlowNode, 0, 10),
	}
	tempResult := &_Result{
		Err:        nil,
		StatusCode: 0,
		StatusMsg:  "",
	}
	res.result = &tempResult
	return res
}

func (c *FlowEngine) Do(functors ...ICallable) *FlowEngine {
	node := NewNormalNode(c.data, c.result, functors...)
	if len(c.nodes) != 0 {
		c.nodes[len(c.nodes)-1].SetNext(node)
	}
	c.nodes = append(c.nodes, node)
	return c
}

func (c *FlowEngine) For(times int, functors ...ICallable) *FlowEngine {
	node := NewForNode(times, c.data, c.result, functors...)
	if len(c.nodes) != 0 {
		c.nodes[len(c.nodes)-1].SetNext(node)
	}
	c.nodes = append(c.nodes, node)
	return c
}

func (c *FlowEngine) Parallel(functors ...ICallable) *FlowEngine {
	node := NewParallelNode(c.data, c.result, functors...)
	if len(c.nodes) != 0 {
		c.nodes[len(c.nodes)-1].SetNext(node)
	}
	c.nodes = append(c.nodes, node)
	return c
}

func (c *FlowEngine) If(condition IBoolFunc, functors ...ICallable) *ElseFlowEngine {
	node := NewIfNode(c.data, c.result, condition, functors...)
	if len(c.nodes) != 0 {
		c.nodes[len(c.nodes)-1].SetNext(node)
	}
	c.nodes = append(c.nodes, node)
	return NewElseFlowEngine(c.data, c, c.result, &c.nodes)
}

func (c *FlowEngine) Wait() *_Result {
	for _, node := range c.nodes {
		node.Run()
	}
	if c.onSuccessFunc != nil {
		if (*c.result).Err == nil && (*c.result).StatusCode == 0 {
			c.onSuccessFunc(c.data, *c.result)
		}
	}
	if c.onFailFunc != nil {
		if (*c.result).Err != nil || (*c.result).StatusCode != 0 {
			c.onFailFunc(c.data, *c.result)
		}
	}
	return *c.result
}

func (c *FlowEngine) SetNote(note string) *FlowEngine {
	if len(c.nodes) != 0 {
		c.nodes[len(c.nodes)-1].SetNote(note)
	}
	return c
}

func (c *FlowEngine) SetBeginLogger(logger INodeBeginLogger) *FlowEngine {
	if len(c.nodes) != 0 {
		c.nodes[len(c.nodes)-1].SetBeginLogger(logger)
	}
	return c
}

func (c *FlowEngine) SetEndLogger(logger INodeEndLogger) *FlowEngine {
	if len(c.nodes) != 0 {
		c.nodes[len(c.nodes)-1].SetEndLogger(logger)
	}
	return c
}

func (c *FlowEngine) SetGlobalBeginLogger(logger INodeBeginLogger) *FlowEngine {
	for _, note := range c.nodes {
		if note.GetBeginLogger() == nil {
			note.SetBeginLogger(logger)
		}
	}
	return c
}

func (c *FlowEngine) SetGlobalEndLogger(logger INodeEndLogger) *FlowEngine {
	for _, note := range c.nodes {
		if note.GetEndLogger() == nil {
			note.SetEndLogger(logger)
		}
	}
	return c
}

func (c *FlowEngine) OnFail(functor IOnFailFunc) *FlowEngine {
	c.onFailFunc = functor
	return c
}

func (c *FlowEngine) OnSuccess(functor IOnFailFunc) *FlowEngine {
	c.onSuccessFunc = functor
	return c
}

//END FlowEngine

//ElseFlowEngine implementation

type ElseFlowEngine struct {
	data          *_Data
	nodes         *[]IBasicFlowNode
	result        **_Result
	invoker       *FlowEngine
	onFailFunc    IOnFailFunc
	onSuccessFunc IOnSuccessFunc
}

func NewElseFlowEngine(data *_Data, invoker *FlowEngine, result **_Result, nodes *[]IBasicFlowNode) *ElseFlowEngine {
	res := &ElseFlowEngine{
		data:    data,
		nodes:   nodes,
		result:  result,
		invoker: invoker,
	}
	return res
}

func (e *ElseFlowEngine) Do(functors ...ICallable) *FlowEngine {
	node := NewNormalNode(e.data, e.result, functors...)
	if len(*e.nodes) != 0 {
		(*e.nodes)[len(*e.nodes)-1].SetNext(node)
	}
	*e.nodes = append(*e.nodes, node)
	return e.invoker
}

func (e *ElseFlowEngine) For(times int, functors ...ICallable) *FlowEngine {
	node := NewForNode(times, e.data, e.result, functors...)
	if len(*e.nodes) != 0 {
		(*e.nodes)[len(*e.nodes)-1].SetNext(node)
	}
	*e.nodes = append(*e.nodes, node)
	return e.invoker
}

func (e *ElseFlowEngine) Parallel(functors ...ICallable) *FlowEngine {
	node := NewParallelNode(e.data, e.result, functors...)
	if len(*e.nodes) != 0 {
		(*e.nodes)[len(*e.nodes)-1].SetNext(node)
	}
	*e.nodes = append(*e.nodes, node)
	return e.invoker
}

func (e *ElseFlowEngine) If(condition IBoolFunc, functors ...ICallable) *ElseFlowEngine {
	node := NewIfNode(e.data, e.result, condition, functors...)
	if len(*e.nodes) != 0 {
		(*e.nodes)[len(*e.nodes)-1].SetNext(node)
	}
	*e.nodes = append(*e.nodes, node)
	return e
}

func (e *ElseFlowEngine) ElseIf(condition IBoolFunc, functors ...ICallable) *ElseFlowEngine {
	node := NewElseIfNode(e.data, e.result, condition, functors...)
	(*e.nodes)[len(*e.nodes)-1].SetNext(node)
	*e.nodes = append(*e.nodes, node)
	return e
}

func (e *ElseFlowEngine) Else(functors ...ICallable) *FlowEngine {
	node := NewElseNode(e.data, e.result, functors...)
	(*e.nodes)[len(*e.nodes)-1].SetNext(node)
	*e.nodes = append(*e.nodes, node)
	return e.invoker
}

func (e *ElseFlowEngine) Wait() *_Result {
	for _, node := range *e.nodes {
		node.Run()
	}
	if e.onSuccessFunc != nil {
		if (*e.result).Err == nil && (*e.result).StatusCode == 0 {
			e.onSuccessFunc(e.data, *e.result)
		}
	}
	if e.onFailFunc != nil {
		if (*e.result).Err != nil || (*e.result).StatusCode != 0 {
			e.onFailFunc(e.data, *e.result)
		}
	}
	return *e.result
}

func (e *ElseFlowEngine) SetNote(note string) *ElseFlowEngine {
	if len(*e.nodes) != 0 {
		(*e.nodes)[len(*e.nodes)-1].SetNote(note)
	}
	return e
}

func (e *ElseFlowEngine) SetBeginLogger(logger INodeBeginLogger) *ElseFlowEngine {
	if len(*e.nodes) != 0 {
		(*e.nodes)[len(*e.nodes)-1].SetBeginLogger(logger)
	}
	return e
}

func (e *ElseFlowEngine) SetEndLogger(logger INodeEndLogger) *ElseFlowEngine {
	if len(*e.nodes) != 0 {
		(*e.nodes)[len(*e.nodes)-1].SetEndLogger(logger)
	}
	return e
}

func (e *ElseFlowEngine) SetGlobalBeginLogger(logger INodeBeginLogger) *ElseFlowEngine {
	for _, note := range *e.nodes {
		if note.GetBeginLogger() == nil {
			note.SetBeginLogger(logger)
		}
	}
	return e
}

func (e *ElseFlowEngine) SetGlobalEndLogger(logger INodeEndLogger) *ElseFlowEngine {
	for _, note := range *e.nodes {
		if note.GetEndLogger() == nil {
			note.SetEndLogger(logger)
		}
	}
	return e
}

func (e *ElseFlowEngine) OnFail(functor IOnFailFunc) *ElseFlowEngine {
	e.onFailFunc = functor
	return e
}

func (e *ElseFlowEngine) OnSuccess(functor IOnFailFunc) *ElseFlowEngine {
	e.onSuccessFunc = functor
	return e
}

//END ElseFlowEngine
