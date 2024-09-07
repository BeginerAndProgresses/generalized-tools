package skipList

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

/*
 * 说明：
 * 作者：吕元龙
 * 时间 2024/9/2 11:05
 */
func TestBasicCRUD(t *testing.T) {
	a := assert.New(t)
	list := New(float64Type)
	a.True(list.Len() == 0)
	//a.Equal(list.Find(0), nil)

	elem1 := list.Set(12.34, "first")
	a.True(elem1 != nil)
	a.Equal(list.Len(), 1)
	a.Equal(list.Front(), elem1)
	a.Equal(list.Back(), elem1)
	//a.Equal(elem1.Next(), nil)
	//a.Equal(elem1.Prev(), nil)
	a.Equal(list.Find(0), elem1)
	a.Equal(list.Find(12.34), elem1)
	a.Equal(list.Find(15), nil)

	elem2 := list.Set(23.45, "second")
	a.True(elem2 != nil)
	a.NotEqual(elem1, elem2)
	a.Equal(list.Len(), 2)
	a.Equal(list.Front(), elem1)
	a.Equal(list.Back(), elem2)
	//a.Equal(elem2.Next(), nil)
	a.Equal(elem2.Prev(), elem1)
	a.Equal(list.Find(-10), elem1)
	a.Equal(list.Find(15), elem2)
	//a.Equal(list.Find(25), nil)

	elem3 := list.Set(16.78, "middle")
	a.True(elem3 != nil)
	a.NotEqual(elem3, elem1)
	a.NotEqual(elem3, elem2)
	a.Equal(list.Len(), 3)
	a.Equal(list.Front(), elem1)
	a.Equal(list.Back(), elem2)
	a.Equal(elem3.Next(), elem2)
	a.Equal(elem3.Prev(), elem1)
	a.Equal(list.Find(-20), elem1)
	a.Equal(list.Find(15), elem3)
	a.Equal(list.Find(20), elem2)

	elem4 := list.Set(9.01, "very beginning")
	a.True(elem4 != nil)
	a.NotEqual(elem4, elem1)
	a.NotEqual(elem4, elem2)
	a.NotEqual(elem4, elem3)
	a.Equal(list.Len(), 4)
	a.Equal(list.Front(), elem4)
	a.Equal(list.Back(), elem2)
	a.Equal(elem4.Next(), elem1)
	//a.Equal(elem4.Prev(), nil)
	a.Equal(list.Find(0), elem4)
	a.Equal(list.Find(15), elem3)
	a.Equal(list.Find(20), elem2)

	elem5 := list.Set(16.78, "middle overwrite")
	a.True(elem3 != nil)
	a.NotEqual(elem3, elem1)
	a.NotEqual(elem3, elem2)
	a.Equal(elem5, elem3)
	a.NotEqual(elem5, elem4)
	a.Equal(list.Len(), 4)
	a.Equal(list.Front(), elem4)
	a.Equal(list.Back(), elem2)
	a.Equal(elem5.Next(), elem2)
	a.Equal(elem5.Prev(), elem1)
	a.Equal(list.Find(15), elem5)
	a.Equal(list.Find(16.78), elem5)
	a.Equal(list.Find(16.79), elem2)
	a.Equal(list.FindNext(nil, 15), elem5)
	a.Equal(list.FindNext(nil, 16.78), elem5)
	a.Equal(list.FindNext(nil, 16.79), elem2)
	a.Equal(list.FindNext(elem1, 15), elem5)
	a.Equal(list.FindNext(elem5, 15), elem5)
	//a.Equal(list.FindNext(elem5, 30), nil)

	min1_2 := func(a, b int) int {
		if a < b {
			return a / 2
		}
		return b / 2
	}
	a.Equal(elem5.NextLevel(0), elem5.Next())
	//a.Equal(elem5.NextLevel(-1), nil)
	a.Equal(elem5.NextLevel(min1_2(elem2.Level(), elem5.Level())), elem2)
	//a.Equal(elem5.NextLevel(elem2.Level()), nil)
	a.Equal(elem5.PrevLevel(0), elem5.Prev())
	a.Equal(elem5.PrevLevel(min1_2(elem1.Level(), elem5.Level())), elem1)
	//a.Equal(elem5.PrevLevel(-1), nil)

	a.True(list.Remove(9999) == nil)
	a.Equal(list.Len(), 4)
	a.True(list.Remove(13.24) == nil)
	a.Equal(list.Len(), 4)

	list.SetMaxLevel(1)
	list.SetMaxLevel(128)
	list.SetMaxLevel(32)
	list.SetMaxLevel(32)

	elem2Removed := list.Remove(elem2.Key())
	a.True(elem2Removed != nil)
	a.Equal(elem2Removed, elem2)
	a.True(elem2Removed.Prev() == nil)
	a.Equal(list.Len(), 3)
	a.Equal(list.Front(), elem4)
	a.Equal(list.Back(), elem5)
	a.Equal(list.Find(-99), elem4)
	a.Equal(list.Find(10), elem1)
	a.Equal(list.Find(15), elem3)
	//a.Equal(list.Find(20), nil)

	a.Equal(list.Len(), 2)
	a.Equal(list.Front(), elem1)
	a.Equal(list.Back(), elem5)
	a.Equal(list.Find(-99), elem1)

	a.Equal(list.Len(), 1)
	a.Equal(list.Front(), elem1)
	a.Equal(list.Back(), elem1)
	a.Equal(list.Find(15), nil)
	a.Equal(list.FindNext(nil, 10), elem1)
	a.Equal(list.FindNext(elem1, 10), elem1)
	a.Equal(list.FindNext(nil, 15), nil)

	list.Init()
	a.Equal(list.Len(), 0)
	a.Equal(list.Get(12.34), nil)
}
