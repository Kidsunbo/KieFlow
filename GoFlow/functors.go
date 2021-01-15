package main

import (
	"context"
	"errors"
	"fmt"
)

func Func1(data *DataTest) *ResultTest {
	fmt.Println("Func1: ", data)
	data.Age++
	return &ResultTest{
		Err:        nil,
		StatusCode: 0,
		StatusMsg:  "",
		AgePlusOne: data.Age + 1,
	}
}

func Func2(data *DataTest) *ResultTest {
	fmt.Println("Func2: ", data)
	data.Age++

	return nil
}

func Func3(data *DataTest) *ResultTest {
	fmt.Println("Func3: ", data)
	data.Age++
	data.Pet = append(data.Pet,"Func3 Pet")
	return &ResultTest{
		Err:        nil,
		StatusCode: 0,
		StatusMsg:  "",
		AgePlusOne: data.Age + 1,
	}
}

func Func4(data *DataTest) *ResultTest {
	fmt.Println("Func4: ", data)
	data.Age++

	return nil
}

func Func5(data *DataTest) *ResultTest {
	fmt.Println("Func5: ", data)
	data.Age++
	data.Pet = append(data.Pet,"Func5 Pet")

	return &ResultTest{
		Err:        nil,
		StatusCode: 0,
		StatusMsg:  "",
		AgePlusOne: data.Age + 1,
	}
}

func Func6(data *DataTest) *ResultTest {
	fmt.Println("Func6: ", data)
	data.Age++

	return nil
}

func Func7(data *DataTest) *ResultTest {
	fmt.Println("Func7: ", data)
	data.Age++

	return &ResultTest{
		Err:        nil,
		StatusCode: 0,
		StatusMsg:  "",
		AgePlusOne: data.Age + 1,
	}
}

func Func8(data *DataTest) *ResultTest {
	fmt.Println("Func8: ", data)
	data.Age++
	data.Pet = append(data.Pet,"Func8 Pet")

	return &ResultTest{
		Err:        errors.New("something wrong"),
		StatusCode: 0,
		StatusMsg:  "",
		AgePlusOne: data.Age + 1,
	}
}

func Func9(data *DataTest) *ResultTest {
	fmt.Println("Func9: ", data)
	data.Age++

	return &ResultTest{
		Err:        nil,
		StatusCode: 10000,
		StatusMsg:  "",
		AgePlusOne: data.Age + 1,
	}
}

func CondFalse(data *DataTest) bool {
	fmt.Println("CondFalse: ", data)
	return false
}

func CondTrue(data *DataTest) bool {
	fmt.Println("CondTrue: ", data)
	return true
}

func Prepare(data *DataTest, input PrepareTest) *DataTest {
	fmt.Println("Prepare:", input)
	data.Name = "Prepare"
	data.Age = 1
	return &DataTest{Ctx: context.Background()}
}


func BeginLogger(note string, _data *DataTest){
	fmt.Println("[BEGIN LOGGER]",note,_data)
}

func EndLogger(note string, _data *DataTest, result *ResultTest){
	fmt.Println("[END LOGGER]",note,_data)
}



func OnSuccessHandle(_data *DataTest, _result *ResultTest){
	fmt.Println("[SUCCESS]:)")
}


func OnFailHandle(_data *DataTest, _result *ResultTest){
	fmt.Println("[FAIL]:(")

}