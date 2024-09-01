package gttype

import (
	"github.com/BeginerAndProgresses/generalized-tools/utils"
	"testing"
)

func TestNewHeap(t *testing.T) {
	heap := NewHeap[any]()
	heap.Insert(1)
	heap.Insert(2)
	heap.PrintHeap()
}

func TestNewZipList(t *testing.T) {
	a := byte(63)
	a = a >> 2
	t.Logf("%v", a)
}

func TestParseZipList_parseHeader(t *testing.T) {
	list := NewZipList[[]byte]()
	a := make([]byte, 255)
	parseElement, err := list.parseElement(nil, a)
	if err != nil {
		t.Error(err)
	}
	t.Logf("%v", parseElement)
	t.Logf("%v", utils.FormatBytesAsBinary([]byte{parseElement.head.thisEntryLen[len(parseElement.head.thisEntryLen)-1]}))
	//t.Logf("%v", utils.FormatBytesAsBinary(parseElement.context))
}
