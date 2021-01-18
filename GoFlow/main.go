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

func EndLogger(note string, data *DataSet, result *Result) {
	fmt.Println("[END]", note, "data = ", data, "result=", result)
}

func main() {
	result := NewFlow().Do(Fun1WithoutData, Fun2WithoutData).SetNote("开始检查数据").
		Prepare(InputParam{}, PrepareData).SetNote("开始准备数据").
		Do(Fun1WithData, Fun2WithData).SetNote("带着数据检查").
		Parallel(Fun2WithData,Fun1WithData).
		SetGlobalBeginLogger(BeginLogger).
		SetGlobalEndLogger(EndLogger).
		Wait()


	fmt.Println("Result=", result)
}
