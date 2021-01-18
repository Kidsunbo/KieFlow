package goflow

import (
	"runtime/debug"
	"sync"
)

type ICallable = func(_data *_Data) *_Result

type IBoolFunc = func(_data *_Data) bool

type IPrepareFunc = func(_data *_Data, input _PrepareInput) *_Data

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

type Flow = FlowEngine

func NewFlow() *Flow {
	return NewFlowEngine()
}

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
	if result != nil {
		b.SetParentResult(result)
	}

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
			if result != nil && (result.Err != nil || result.StatusCode != 0) {
				return result
			}
		}
		current := i.Next
		for current != nil && (current.GetNodeType() == ElseIfNodeType || current.GetNodeType() == ElseNodeType) {
			current.SetShouldSkip(true)
			current = current.GetNext()
		}
	}

	return i.GetParentResult()
}

func (i *IfNode) Run() {
	if i.ShouldSkip || i.GetParentResult().Err != nil || i.GetParentResult().StatusCode != 0 {
		return
	}
	if i.BeginLogger != nil {
		i.BeginLogger(i.Note, i.Data)
	}

	result := i.ImplTask()
	if result != nil {
		i.SetParentResult(result)
	}

	if i.EndLogger != nil {
		i.EndLogger(i.Note, i.Data, i.GetParentResult())
	}
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
		if result != nil && (result.Err != nil || result.StatusCode != 0) {
			return result
		}
	}
	return e.GetParentResult()
}

func (e *ElseNode) Run() {
	if e.ShouldSkip || e.GetParentResult().Err != nil || e.GetParentResult().StatusCode != 0 {
		return
	}
	if e.BeginLogger != nil {
		e.BeginLogger(e.Note, e.Data)
	}

	result := e.ImplTask()
	if result != nil {
		e.SetParentResult(result)
	}

	if e.EndLogger != nil {
		e.EndLogger(e.Note, e.Data, e.GetParentResult())
	}
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

func (e *ElseIfNode) ImplTask() *_Result {
	if e.Condition == nil {
		return &_Result{
			Err:        NewConditionNotFoundError(),
			StatusCode: 0,
			StatusMsg:  "",
		}
	}

	if e.Condition(e.Data) {
		for _, functor := range e.Functors {
			result := functor(e.Data)
			if result != nil && (result.Err != nil || result.StatusCode != 0) {
				return result
			}
		}

		current := e.Next
		for current != nil && (current.GetNodeType() == ElseIfNodeType || current.GetNodeType() == ElseNodeType) {
			current.SetShouldSkip(true)
			current = current.GetNext()
		}
	}


	return e.GetParentResult()
}

func (e *ElseIfNode) Run() {
	if e.ShouldSkip || e.GetParentResult().Err != nil || e.GetParentResult().StatusCode != 0 {
		return
	}
	if e.BeginLogger != nil {
		e.BeginLogger(e.Note, e.Data)
	}

	result := e.ImplTask()
	if result != nil {
		e.SetParentResult(result)
	}

	if e.EndLogger != nil {
		e.EndLogger(e.Note, e.Data, e.GetParentResult())
	}
}

//END ElseIfNode

//NormalNode Implementation
type NormalNode struct {
	*BasicFlowNode
}

func NewNormalNode(data *_Data, parentResult **_Result, functors ...ICallable) *ElseNode {
	return &ElseNode{NewBasicFlowNode(data, parentResult, NormalNodeType, functors...)}
}

func (n *NormalNode) ImplTask() *_Result {
	for _, functor := range n.Functors {
		result := functor(n.Data)
		if result != nil && (result.Err != nil || result.StatusCode != 0) {
			return result
		}
	}
	return n.GetParentResult()
}

func (n *NormalNode) Run() {
	if n.ShouldSkip || n.GetParentResult().Err != nil || n.GetParentResult().StatusCode != 0 {
		return
	}
	if n.BeginLogger != nil {
		n.BeginLogger(n.Note, n.Data)
	}

	result := n.ImplTask()
	if result != nil {
		n.SetParentResult(result)
	}

	if n.EndLogger != nil {
		n.EndLogger(n.Note, n.Data, n.GetParentResult())
	}
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

func (f *ForNode) ImplTask() *_Result {
	for i := 0; i < f.Times; i++ {
		for _, functor := range f.Functors {
			result := functor(f.Data)
			if result != nil && (result.Err != nil || result.StatusCode != 0) {
				return result
			}
		}
	}
	return f.GetParentResult()
}

func (f *ForNode) Run() {
	if f.ShouldSkip || f.GetParentResult().Err != nil || f.GetParentResult().StatusCode != 0 {
		return
	}
	if f.BeginLogger != nil {
		f.BeginLogger(f.Note, f.Data)
	}

	result := f.ImplTask()
	if result != nil {
		f.SetParentResult(result)
	}

	if f.EndLogger != nil {
		f.EndLogger(f.Note, f.Data, f.GetParentResult())
	}
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

func (p *ParallelNode) ImplTask() *_Result {
	resultChan := make(chan *_Result, len(p.Functors))

	wg := sync.WaitGroup{}
	wg.Add(len(p.Functors))

	go func(wg *sync.WaitGroup) {
		wg.Wait()
		close(resultChan)
	}(&wg)

	for _, functor := range p.Functors {
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
			result := f(p.Data)
			resultChan <- result
		}(&wg, functor)
	}

	result := p.GetParentResult()
	for item := range resultChan {
		if result != nil && (result.StatusCode != 0 || result.Err != nil) {
			continue
		}
		result = item
	}

	return result
}

func (p *ParallelNode) Run() {
	if p.ShouldSkip || p.GetParentResult().Err != nil || p.GetParentResult().StatusCode != 0 {
		return
	}
	if p.BeginLogger != nil {
		p.BeginLogger(p.Note, p.Data)
	}

	result := p.ImplTask()
	if result != nil {
		p.SetParentResult(result)
	}

	if p.EndLogger != nil {
		p.EndLogger(p.Note, p.Data, p.GetParentResult())
	}
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

func NewFlowEngine() *FlowEngine {
	res := &FlowEngine{
		nodes: make([]IBasicFlowNode, 0, 10),
	}
	tempResult := &_Result{
		Err:        nil,
		StatusCode: 0,
		StatusMsg:  "",
	}
	res.data = new(_Data)
	res.result = &tempResult
	return res
}

func (f *FlowEngine) Prepare(prepareFunc IPrepareFunc, input _PrepareInput) *FlowEngine {
	data := prepareFunc(f.data, input)
	if data != nil {
		f.data = data
	}
	return f
}

func (f *FlowEngine) Do(functors ...ICallable) *FlowEngine {
	node := NewNormalNode(f.data, f.result, functors...)
	if len(f.nodes) != 0 {
		f.nodes[len(f.nodes)-1].SetNext(node)
	}
	f.nodes = append(f.nodes, node)
	return f
}

func (f *FlowEngine) For(times int, functors ...ICallable) *FlowEngine {
	node := NewForNode(times, f.data, f.result, functors...)
	if len(f.nodes) != 0 {
		f.nodes[len(f.nodes)-1].SetNext(node)
	}
	f.nodes = append(f.nodes, node)
	return f
}

func (f *FlowEngine) Parallel(functors ...ICallable) *FlowEngine {
	node := NewParallelNode(f.data, f.result, functors...)
	if len(f.nodes) != 0 {
		f.nodes[len(f.nodes)-1].SetNext(node)
	}
	f.nodes = append(f.nodes, node)
	return f
}

func (f *FlowEngine) If(condition IBoolFunc, functors ...ICallable) *ElseFlowEngine {
	node := NewIfNode(f.data, f.result, condition, functors...)
	if len(f.nodes) != 0 {
		f.nodes[len(f.nodes)-1].SetNext(node)
	}
	f.nodes = append(f.nodes, node)
	return NewElseFlowEngine(&f.data, f, f.result, &f.nodes)
}

func (f *FlowEngine) Wait() *_Result {
	for _, node := range f.nodes {
		node.Run()
	}
	if f.onSuccessFunc != nil {
		if (*f.result).Err == nil && (*f.result).StatusCode == 0 {
			f.onSuccessFunc(f.data, *f.result)
		}
	}
	if f.onFailFunc != nil {
		if (*f.result).Err != nil || (*f.result).StatusCode != 0 {
			f.onFailFunc(f.data, *f.result)
		}
	}
	return *f.result
}

func (f *FlowEngine) SetNote(note string) *FlowEngine {
	if len(f.nodes) != 0 {
		f.nodes[len(f.nodes)-1].SetNote(note)
	}
	return f
}

func (f *FlowEngine) SetBeginLogger(logger INodeBeginLogger) *FlowEngine {
	if len(f.nodes) != 0 {
		f.nodes[len(f.nodes)-1].SetBeginLogger(logger)
	}
	return f
}

func (f *FlowEngine) SetEndLogger(logger INodeEndLogger) *FlowEngine {
	if len(f.nodes) != 0 {
		f.nodes[len(f.nodes)-1].SetEndLogger(logger)
	}
	return f
}

func (f *FlowEngine) SetGlobalBeginLogger(logger INodeBeginLogger) *FlowEngine {
	for _, note := range f.nodes {
		if note.GetBeginLogger() == nil {
			note.SetBeginLogger(logger)
		}
	}
	return f
}

func (f *FlowEngine) SetGlobalEndLogger(logger INodeEndLogger) *FlowEngine {
	for _, note := range f.nodes {
		if note.GetEndLogger() == nil {
			note.SetEndLogger(logger)
		}
	}
	return f
}

func (f *FlowEngine) OnFail(functor IOnFailFunc) *FlowEngine {
	f.onFailFunc = functor
	return f
}

func (f *FlowEngine) OnSuccess(functor IOnFailFunc) *FlowEngine {
	f.onSuccessFunc = functor
	return f
}

//END FlowEngine

//ElseFlowEngine implementation

type ElseFlowEngine struct {
	data          **_Data
	nodes         *[]IBasicFlowNode
	result        **_Result
	invoker       *FlowEngine
	onFailFunc    IOnFailFunc
	onSuccessFunc IOnSuccessFunc
}

func NewElseFlowEngine(data **_Data, invoker *FlowEngine, result **_Result, nodes *[]IBasicFlowNode) *ElseFlowEngine {
	res := &ElseFlowEngine{
		data:    data,
		nodes:   nodes,
		result:  result,
		invoker: invoker,
	}
	return res
}

func (e *ElseFlowEngine) Prepare(prepareFunc IPrepareFunc, input _PrepareInput) *ElseFlowEngine {
	data := prepareFunc(*e.data, input)
	if data != nil {
		*e.data = data
	}
	return e
}

func (e *ElseFlowEngine) Do(functors ...ICallable) *FlowEngine {
	node := NewNormalNode(*e.data, e.result, functors...)
	if len(*e.nodes) != 0 {
		(*e.nodes)[len(*e.nodes)-1].SetNext(node)
	}
	*e.nodes = append(*e.nodes, node)
	return e.invoker
}

func (e *ElseFlowEngine) For(times int, functors ...ICallable) *FlowEngine {
	node := NewForNode(times, *e.data, e.result, functors...)
	if len(*e.nodes) != 0 {
		(*e.nodes)[len(*e.nodes)-1].SetNext(node)
	}
	*e.nodes = append(*e.nodes, node)
	return e.invoker
}

func (e *ElseFlowEngine) Parallel(functors ...ICallable) *FlowEngine {
	node := NewParallelNode(*e.data, e.result, functors...)
	if len(*e.nodes) != 0 {
		(*e.nodes)[len(*e.nodes)-1].SetNext(node)
	}
	*e.nodes = append(*e.nodes, node)
	return e.invoker
}

func (e *ElseFlowEngine) If(condition IBoolFunc, functors ...ICallable) *ElseFlowEngine {
	node := NewIfNode(*e.data, e.result, condition, functors...)
	if len(*e.nodes) != 0 {
		(*e.nodes)[len(*e.nodes)-1].SetNext(node)
	}
	*e.nodes = append(*e.nodes, node)
	return e
}

func (e *ElseFlowEngine) ElseIf(condition IBoolFunc, functors ...ICallable) *ElseFlowEngine {
	node := NewElseIfNode(*e.data, e.result, condition, functors...)
	(*e.nodes)[len(*e.nodes)-1].SetNext(node)
	*e.nodes = append(*e.nodes, node)
	return e
}

func (e *ElseFlowEngine) Else(functors ...ICallable) *FlowEngine {
	node := NewElseNode(*e.data, e.result, functors...)
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
			e.onSuccessFunc(*e.data, *e.result)
		}
	}
	if e.onFailFunc != nil {
		if (*e.result).Err != nil || (*e.result).StatusCode != 0 {
			e.onFailFunc(*e.data, *e.result)
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
