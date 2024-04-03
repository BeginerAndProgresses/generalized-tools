package test

import (
	gttype "generalized_tools/type"
	"testing"
)

func TestStack(t *testing.T) {
	stack := gttype.NewStack[int]()
	stack.Push(1)
	stack.Push(4)
	stack.Push(2)
	stack.Push(3)
	println(stack.Top())
	println(stack.Size())
	println(stack.Pop())
	println(stack.Size())
	println(stack.Pop())
	println(stack.Pop())
	println(stack.Pop())
	println(stack.Size())
	println(stack.Pop())
}
