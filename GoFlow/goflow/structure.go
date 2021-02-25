package goflow

import "context"

//************************DEFINE YOUR STRUCTURE BELOW****************************//
// The name starts with underscore means replaceable.
// [IMPORTANT] Notice that even though _Result can be replace with other type, the Err, StatusCode and StatusMsg must be provided

type _Data struct {
	Ctx context.Context  `json:"-"`
	FunctionName string  `json:"-"`
}

type _Result struct {
	Err        error
	StatusCode int64
	StatusMsg  string
}

type _PrepareInput struct {
	Ctx context.Context `json:"-"`
	FunctionName string `json:"-"`
}

func InitStructure(data *_Data){
	//Initialize your data here if needed
}

//************************DEFINE YOUR STRUCTURE ABOVE****************************//
