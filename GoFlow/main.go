package main

import "fmt"

func Fun1WithoutData(data *DataSet) *Result {
	fmt.Println("Fun1WithoutData, ", data)
	return nil
}

func Fun2WithoutData(data *DataSet) *Result {
	fmt.Println("Fun2WithoutData, ", data)
	return new(Result)
}

func Fun1WithData(data *DataSet) *Result {
	fmt.Println("Fun1WithData, ", data)
	return new(Result)
}

func Fun2WithData(data *DataSet) *Result {
	fmt.Println("Fun2WithData, ", data)
	return nil
}

func PrepareData(data *DataSet, input InputParam) *Result {
	fmt.Println("PrepareData, data=", data, "input=", input)
	data.Name = "Tom is Tom"
	result := new(Result)
	return result
}

func PrepareDataWithErr(data *DataSet, input InputParam) *Result {
	fmt.Println("PrepareData, data=", data, "input=", input)
	data.Name = "Tom is Tom"
	result := new(Result)
	result.StatusCode = 1
	return result
}

func BeginLogger(note string, data *DataSet) {
	fmt.Println("\n[START]", note, "data = ", data)
}

func EndLogger(funcName string) func(int64) INodeEndLogger {
	return func(line int64) INodeEndLogger {
		return func(note string, data *DataSet, result *Result) {
			fmt.Println("[END] func=", funcName, "line=", line, note, "data = ", data, "result=", result)
		}
	}
}

func CondTrue(data *DataSet)bool{
	return true
}

func CondFalse(data *DataSet)bool{
	return false
}

func SimpleFunc1(data* DataSet)*Result{
	fmt.Println("SimpleFunc1")
	return nil
}

func SimpleFunc2(data* DataSet)*Result{
	fmt.Println("SimpleFunc2")
	return nil
}

func SimpleFunc3(data* DataSet)*Result{
	fmt.Println("SimpleFunc3")
	return nil
}

func SimpleFunc4(data* DataSet)*Result{
	fmt.Println("SimpleFunc4")
	return nil
}

func SimpleFunc5(data* DataSet)*Result{
	fmt.Println("SimpleFunc5")
	return nil
}

func SimpleFunc6(data* DataSet)*Result{
	fmt.Println("SimpleFunc6")
	result := new(Result)
	result.StatusCode = 1
	return result
}

func main() {
	flow :=NewFlow()
	result := flow.Do(Fun1WithoutData, Fun2WithoutData).SetNote("开始检查数据").
		Prepare(InputParam{}, PrepareData).SetNote("开始准备数据").
		Do(Fun1WithData, Fun2WithData).SetNote("带着数据检查").
		Parallel(Fun2WithData, Fun1WithData).SetNote("并行跑一跑").
		If(CondFalse,SimpleFunc1).
		//IfSubPath(CondFalse,NewFlow().Do(SimpleFunc1).Do(SimpleFunc2).If(CondTrue,SimpleFunc2).Else(SimpleFunc3)).
		//ElseSubPath(NewFlow().Do(SimpleFunc4).Do(SimpleFunc5).Do(SimpleFunc6)).
		SetGlobalBeginLogger(BeginLogger).
		SetGlobalEndLogger(EndLogger("hello")(59)).
		Wait()

	fmt.Println("Result=", result)
}
