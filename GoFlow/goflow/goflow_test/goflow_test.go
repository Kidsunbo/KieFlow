package goflow_test

import (
	"context"
	"testing"
)

func NewData()*_Data{
	return &_Data{Ctx: context.Background()}
}


func TestNewBasicFlowNode(t *testing.T) {
	NewBasicFlowNode()
}