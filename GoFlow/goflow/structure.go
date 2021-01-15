package goflow

import "context"

//************************DEFINE YOUR STRUCTURE BELOW****************************//
// The name starts with underscore means replaceable.
// [IMPORTANT] Notice that even though _Result can be replace with other type, the Err, StatusCode and StatusMsg must be provided

type _Data struct {
	Ctx context.Context
}

type _Result struct {
	Err        error
	StatusCode int64
	StatusMsg  string
}

type _PrepareInput struct {
	Ctx context.Context
}

//************************DEFINE YOUR STRUCTURE ABOVE****************************//
