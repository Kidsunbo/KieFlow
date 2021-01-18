package main

import "context"

//************************DEFINE YOUR STRUCTURE BELOW****************************//
// The name starts with underscore means replaceable.
// [IMPORTANT] Notice that even though Result can be replace with other type, the Err, StatusCode and StatusMsg must be provided

type DataSet struct {
	Ctx context.Context
	Name string
}

type Result struct {
	Err        error
	StatusCode int64
	StatusMsg  string
}

type InputParam struct {
	Ctx context.Context
}

//************************DEFINE YOUR STRUCTURE ABOVE****************************//
