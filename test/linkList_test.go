package test

import (
	"fmt"
	gttype "github.com/BeginerAndProgresses/generalized-tools/type"
	"testing"
)

func TestLinkList(t *testing.T) {
	type hehe struct {
		name string
	}
	list := gttype.NewLinkList[hehe]()
	fmt.Println(list.Size())
	list.Insert(1, hehe{name: "张三"})
	fmt.Println(list.Find(1))
	fmt.Println(list.Size())
	fmt.Println(list.Delete(1))
	fmt.Println(list.Size())
	fmt.Println(list.Delete(1))
	list.Insert(2, hehe{name: "张三"})
}
