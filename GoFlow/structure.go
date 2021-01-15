package main

import "context"

//************************DEFINE YOUR STRUCTURE BELOW****************************//
// The name starts with underscore means replaceable.
// [IMPORTANT] Notice that even though ResultTest can be replace with other type, the Err, StatusCode and StatusMsg must be provided

type DataTest struct {
	Ctx  context.Context
	Name string
	Age  int64
	Pet  []string
}

type ResultTest struct {
	Err        error
	StatusCode int64
	StatusMsg  string
	AgePlusOne int64
}

type PrepareTest struct {
	Ctx  context.Context
	Name int64
}

//************************DEFINE YOUR STRUCTURE ABOVE****************************//
