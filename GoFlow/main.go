package main

import (
	"fmt"
)

func SimpleWorkFlow() {
	result := NewFlow().
		Do(Func1).
		Prepare(Prepare, PrepareTest{}).
		If(CondTrue, Func2, Func3).
		Else(Func4, Func5).
		Do(Func4, Func1, Func6).SetNote("订单类型校验").
		If(CondFalse, Func2, Func3).
		ElseIf(CondFalse, Func6, Func5).
		Else(Func1, Func2).
		For(3, Func7).
		Parallel(Func9, Func2, Func1, Func6).
		If(CondTrue, Func2).
		SetGlobalEndLogger(EndLogger).
		Wait()
	fmt.Println(result)
}

func NoteCheck() {
	result := NewFlow().
		Do(Func1).SetNote("first").SetBeginLogger(BeginLogger).SetEndLogger(EndLogger).
		Do(Func2).
		Do(Func3).
		SetGlobalBeginLogger(BeginLogger).
		Wait()
	fmt.Println(result)
}

func OnSuccessAndFailCheck() {
	result := NewFlow().
		Do(Func1).
		Do(Func2).
		Do(Func3).
		Do(Func4).
		Do(Func5).
		Do(Func6).
		Do(Func7).
		//Do(Func8).
		Do(Func9).
		Do(Func1).
		Do(Func2).
		Do(Func3).
		OnSuccess(OnSuccessHandle).
		OnFail(OnFailHandle).
		Wait()
	fmt.Println(result)
}

func main() {
	fmt.Println("Welcome to GoFlow!")

	SimpleWorkFlow()
	NoteCheck()
	OnSuccessAndFailCheck()
}
