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
		Do(Func4, Func5, Func6).
		If(CondFalse, Func2, Func3).
		ElseIf(CondFalse, Func6, Func5).
		Else(Func1, Func2).
		For(3, Func9).
		Parallel(Func9, Func2, Func1, Func6).
		If(CondTrue, Func2).
		Wait()
	fmt.Println(result)
}

func main() {
	fmt.Println("Welcome to GoFlow!")

	SimpleWorkFlow()

}
