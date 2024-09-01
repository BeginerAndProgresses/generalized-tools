package gttype

import (
	"testing"
)

/*
 * 说明：
 * 作者：吕元龙
 * 时间 2024/8/31 19:08
 */

func TestAdkZipList_Insert(t *testing.T) {
	z := NewZipList[string]()
	err := z.Push("你好")
	if err != nil {
		t.Error(err)
	}
	t.Logf("%v", z)
	z.Push("我好")
}

func TestAdkZipList_Pop(t *testing.T) {
	z := NewZipList[string]()
	err := z.Push("你好")
	if err != nil {
		t.Error(err)
	}
	z.Push("我好")
	t.Logf("%v", z.Pop())
	t.Logf("%v", z.Pop())
	t.Logf("%v", z.Pop())
}

func TestAdkZipList_String(t *testing.T) {
	z := NewZipList[string]()
	z.Push("你好")
	z.Push("我好")
	t.Logf("%v", z.String())
}
