package test

import (
	"fmt"
	"github.com/BeginerAndProgresses/generalized-tools/slice"
	"testing"
)

func TestBase(t *testing.T) {
	insert, _ := slice.Insert[int](0, 3, []int{0, 0, 0})
	fmt.Println(insert)
}
