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

func NewFlow()*Flow{
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

func (i *IfNode) SetParentResult(result *_Result) {
	*i.parentResult = result
}

func (i *IfNode) GetParentResult() *_Result {
	return *i.parentResult
}

func (i *IfNode) Run() {
	if i.ShouldSkip || i.GetParentResult().Err != nil || i.GetParentResult().StatusCode != 0 {
		return
	}
	if i.BeginLogger != nil {
		i.BeginLogger(i.Note, i.Data)
	}

	result := i.ImplTask()
	i.SetParentResult(result)

	if i.EndLogger != nil {
		i.EndLogger(i.Note, i.Data, i.GetParentResult())
	}
}

func (i *IfNode) GetNext() IBasicFlowNode {
	return i.Next
}

func (i *IfNode) SetNext(node IBasicFlowNode) {
	i.Next = node
}

func (i *IfNode) GetNodeType() NodeType {
	return i.NodeType
}

func (i *IfNode) SetShouldSkip(shouldSkip bool) {
	i.ShouldSkip = shouldSkip
}

func (i *IfNode) SetNote(note string) {
	i.Note = note
}

func (i *IfNode) GetNote() string {
	return i.Note
}

func (i *IfNode) SetBeginLogger(logger INodeBeginLogger) {
	i.BeginLogger = logger
}

func (i *IfNode) GetBeginLogger() INodeBeginLogger {
	return i.BeginLogger
}

func (i *IfNode) SetEndLogger(logger INodeEndLogger) {
	i.EndLogger = logger
}

func (i *IfNode) GetEndLogger() INodeEndLogger {
	return i.EndLogger
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

func (e *ElseNode) SetParentResult(result *_Result) {
	*e.parentResult = result
}

func (e *ElseNode) GetParentResult() *_Result {
	return *e.parentResult
}

func (e *ElseNode) Run() {
	if e.ShouldSkip || e.GetParentResult().Err != nil || e.GetParentResult().StatusCode != 0 {
		return
	}
	if e.BeginLogger != nil {
		e.BeginLogger(e.Note, e.Data)
	}

	result := e.ImplTask()
	e.SetParentResult(result)

	if e.EndLogger != nil {
		e.EndLogger(e.Note, e.Data, e.GetParentResult())
	}
}

func (e *ElseNode) GetNext() IBasicFlowNode {
	return e.Next
}

func (e *ElseNode) SetNext(node IBasicFlowNode) {
	e.Next = node
}

func (e *ElseNode) GetNodeType() NodeType {
	return e.NodeType
}

func (e *ElseNode) SetShouldSkip(shouldSkip bool) {
	e.ShouldSkip = shouldSkip
}

func (e *ElseNode) SetNote(note string) {
	e.Note = note
}

func (e *ElseNode) GetNote() string {
	return e.Note
}

func (e *ElseNode) SetBeginLogger(logger INodeBeginLogger) {
	e.BeginLogger = logger
}

func (e *ElseNode) GetBeginLogger() INodeBeginLogger {
	return e.BeginLogger
}

func (e *ElseNode) SetEndLogger(logger INodeEndLogger) {
	e.EndLogger = logger
}

func (e *ElseNode) GetEndLogger() INodeEndLogger {
	return e.EndLogger
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
			if result.Err != nil || result.StatusCode != 0 {
				return result
			}
		}
	}

	current := e.Next
	for current != nil && (current.GetNodeType() == ElseIfNodeType || current.GetNodeType() == ElseNodeType) {
		current.SetShouldSkip(true)
		current = current.GetNext()
	}

	return e.GetParentResult()
}

func (e *ElseIfNode) SetParentResult(result *_Result) {
	*e.parentResult = result
}

func (e *ElseIfNode) GetParentResult() *_Result {
	return *e.parentResult
}

func (e *ElseIfNode) Run() {
	if e.ShouldSkip || e.GetParentResult().Err != nil || e.GetParentResult().StatusCode != 0 {
		return
	}
	if e.BeginLogger != nil {
		e.BeginLogger(e.Note, e.Data)
	}

	result := e.ImplTask()
	e.SetParentResult(result)

	if e.EndLogger != nil {
		e.EndLogger(e.Note, e.Data, e.GetParentResult())
	}
}

func (e *ElseIfNode) GetNext() IBasicFlowNode {
	return e.Next
}

func (e *ElseIfNode) SetNext(node IBasicFlowNode) {
	e.Next = node
}

func (e *ElseIfNode) GetNodeType() NodeType {
	return e.NodeType
}

func (e *ElseIfNode) SetShouldSkip(shouldSkip bool) {
	e.ShouldSkip = shouldSkip
}

func (e *ElseIfNode) SetNote(note string) {
	e.Note = note
}

func (e *ElseIfNode) GetNote() string {
	return e.Note
}

func (e *ElseIfNode) SetBeginLogger(logger INodeBeginLogger) {
	e.BeginLogger = logger
}

func (e *ElseIfNode) GetBeginLogger() INodeBeginLogger {
	return e.BeginLogger
}

func (e *ElseIfNode) SetEndLogger(logger INodeEndLogger) {
	e.EndLogger = logger
}

func (e *ElseIfNode) GetEndLogger() INodeEndLogger {
	return e.EndLogger
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
		if result.Err != nil || result.StatusCode != 0 {
			return result
		}
	}
	return n.GetParentResult()
}

func (n *NormalNode) SetParentResult(result *_Result) {
	*n.parentResult = result
}

func (n *NormalNode) GetParentResult() *_Result {
	return *n.parentResult
}

func (n *NormalNode) Run() {
	if n.ShouldSkip || n.GetParentResult().Err != nil || n.GetParentResult().StatusCode != 0 {
		return
	}
	if n.BeginLogger != nil {
		n.BeginLogger(n.Note, n.Data)
	}

	result := n.ImplTask()
	n.SetParentResult(result)

	if n.EndLogger != nil {
		n.EndLogger(n.Note, n.Data, n.GetParentResult())
	}
}

func (n *NormalNode) GetNext() IBasicFlowNode {
	return n.Next
}

func (n *NormalNode) SetNext(node IBasicFlowNode) {
	n.Next = node
}

func (n *NormalNode) GetNodeType() NodeType {
	return n.NodeType
}

func (n *NormalNode) SetShouldSkip(shouldSkip bool) {
	n.ShouldSkip = shouldSkip
}

func (n *NormalNode) SetNote(note string) {
	n.Note = note
}

func (n *NormalNode) GetNote() string {
	return n.Note
}

func (n *NormalNode) SetBeginLogger(logger INodeBeginLogger) {
	n.BeginLogger = logger
}

func (n *NormalNode) GetBeginLogger() INodeBeginLogger {
	return n.BeginLogger
}

func (n *NormalNode) SetEndLogger(logger INodeEndLogger) {
	n.EndLogger = logger
}

func (n *NormalNode) GetEndLogger() INodeEndLogger {
	return n.EndLogger
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
			if result.Err != nil || result.StatusCode != 0 {
				return result
			}
		}
	}
	return f.GetParentResult()
}

func (f *ForNode) SetParentResult(result *_Result) {
	*f.parentResult = result
}

func (f *ForNode) GetParentResult() *_Result {
	return *f.parentResult
}

func (f *ForNode) Run() {
	if f.ShouldSkip || f.GetParentResult().Err != nil || f.GetParentResult().StatusCode != 0 {
		return
	}
	if f.BeginLogger != nil {
		f.BeginLogger(f.Note, f.Data)
	}

	result := f.ImplTask()
	f.SetParentResult(result)

	if f.EndLogger != nil {
		f.EndLogger(f.Note, f.Data, f.GetParentResult())
	}
}

func (f *ForNode) GetNext() IBasicFlowNode {
	return f.Next
}

func (f *ForNode) SetNext(node IBasicFlowNode) {
	f.Next = node
}

func (f *ForNode) GetNodeType() NodeType {
	return f.NodeType
}

func (f *ForNode) SetShouldSkip(shouldSkip bool) {
	f.ShouldSkip = shouldSkip
}

func (f *ForNode) SetNote(note string) {
	f.Note = note
}

func (f *ForNode) GetNote() string {
	return f.Note
}

func (f *ForNode) SetBeginLogger(logger INodeBeginLogger) {
	f.BeginLogger = logger
}

func (f *ForNode) GetBeginLogger() INodeBeginLogger {
	return f.BeginLogger
}

func (f *ForNode) SetEndLogger(logger INodeEndLogger) {
	f.EndLogger = logger
}

func (f *ForNode) GetEndLogger() INodeEndLogger {
	return f.EndLogger
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
		if result.StatusCode != 0 || result.Err != nil {
			continue
		}
		result = item
	}

	return result
}

func (p *ParallelNode) SetParentResult(result *_Result) {
	*p.parentResult = result
}

func (p *ParallelNode) GetParentResult() *_Result {
	return *p.parentResult
}

func (p *ParallelNode) Run() {
	if p.ShouldSkip || p.GetParentResult().Err != nil || p.GetParentResult().StatusCode != 0 {
		return
	}
	if p.BeginLogger != nil {
		p.BeginLogger(p.Note, p.Data)
	}

	result := p.ImplTask()
	p.SetParentResult(result)

	if p.EndLogger != nil {
		p.EndLogger(p.Note, p.Data, p.GetParentResult())
	}
}

func (p *ParallelNode) GetNext() IBasicFlowNode {
	return p.Next
}

func (p *ParallelNode) SetNext(node IBasicFlowNode) {
	p.Next = node
}

func (p *ParallelNode) GetNodeType() NodeType {
	return p.NodeType
}

func (p *ParallelNode) SetShouldSkip(shouldSkip bool) {
	p.ShouldSkip = shouldSkip
}

func (p *ParallelNode) SetNote(note string) {
	p.Note = note
}

func (p *ParallelNode) GetNote() string {
	return p.Note
}

func (p *ParallelNode) SetBeginLogger(logger INodeBeginLogger) {
	p.BeginLogger = logger
}

func (p *ParallelNode) GetBeginLogger() INodeBeginLogger {
	return p.BeginLogger
}

func (p *ParallelNode) SetEndLogger(logger INodeEndLogger) {
	p.EndLogger = logger
}

func (p *ParallelNode) GetEndLogger() INodeEndLogger {
	return p.EndLogger
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
	res.result = &tempResult
	return res
}

func (c *FlowEngine) Prepare(prepareFunc IPrepareFunc, input _PrepareInput) *FlowEngine {
	data := prepareFunc(c.data, input)
	c.data = data
	return c
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
