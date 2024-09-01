package gttype

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"github.com/BeginerAndProgresses/generalized-tools/utils"
)

/*
 * 说明：
 * 作者：吕元龙
 * 时间 2024/8/29 18:42
 */

const (
	encodingLongLong = 0x80
	encodingLong     = 0x40
	encodingShort    = 0x00
	preSizeShort     = 0x00
	preSizeLong      = 0xfe
	preSizeLongLong  = 0xfe + 0xffffffff
)

// ZipList 压缩列表, 用于存储在Redis中用于存储短链数据
// 详细实现：https://matt.sh/redis-quicklist
type ZipList[T any] interface {
	parseElement(pre []byte, val T) (*element, error)
	//Insert(val T) error
	Push(val T) error
	Pop() T
	String() string
}

// adkZipList 压缩列表实现
type adkZipList[T any] struct {
	data []byte
	// 序列化函数,初始化时使用Json序列化
	marshal   func(T) ([]byte, error)
	unmarshal func(data []byte, t *T) error
	//	数据字节最大长度
	maxLen int
}

type eHead struct {
	preEntryLen  []byte
	thisEntryLen []byte
}

type element struct {
	head    eHead
	context []byte
}

// parseElement 处理元素
// 前一个元素的长度
func (a *adkZipList[T]) parseElement(pre []byte, val T) (*element, error) {
	e := &element{
		head: eHead{},
	}
	pl := len(pre)
	if pl == preSizeShort {
		//第一个元素
		e.head.preEntryLen = make([]byte, 1)
	} else if pl < preSizeLong {
		e.head.preEntryLen = make([]byte, 1)
		//	将int值转成字节数组
		e.head.preEntryLen[0] = byte(pl)
	} else if pl <= preSizeLongLong {
		e.head.preEntryLen = make([]byte, 5)
		//bytes := utils.IntToBytes(pl)
		e.head.preEntryLen[0] = 0xfe
		binary.LittleEndian.PutUint32(e.head.preEntryLen[1:5], uint32(pl)-0xfe)
	} else {
		return nil, errors.New("数据过长")
	}
	marshal, err := a.marshal(val)
	if err != nil {
		return nil, err
	}
	ml := len(marshal)
	if ml <= 63 {
		e.head.thisEntryLen = make([]byte, 1)
		e.head.thisEntryLen[0] = byte(ml)
		e.head.thisEntryLen[0] |= encodingShort
	} else if len(marshal) <= 16383 {
		e.head.thisEntryLen = make([]byte, 2)
		binary.LittleEndian.PutUint16(e.head.thisEntryLen, uint16(ml))
		e.head.thisEntryLen[0] |= encodingLong
	} else if len(marshal) <= 4294967294 {
		e.head.thisEntryLen = make([]byte, 5)
		binary.LittleEndian.PutUint32(e.head.thisEntryLen[:4], uint32(ml))
		e.head.thisEntryLen[0] |= encodingLongLong
	} else {
		return nil, errors.New("数据过长")
	}
	e.context = make([]byte, 0)
	e.context = append(e.context, marshal...)
	return e, nil
}

func NewZipList[T any](maxLen ...int) ZipList[T] {
	a := &adkZipList[T]{
		data: make([]byte, 11),
		marshal: func(t T) ([]byte, error) {
			return json.Marshal(t)
		},
		unmarshal: func(data []byte, t *T) error {
			return json.Unmarshal(data, t)
		},
	}
	if len(maxLen) > 0 {
		a.maxLen = maxLen[0]
	} else {
		a.maxLen = 64
	}
	// z-header
	binary.LittleEndian.PutUint32(a.data[0:4], 11) //z-bytes 比特数
	binary.LittleEndian.PutUint32(a.data[4:8], 10) //z-tail 尾部偏移量
	binary.LittleEndian.PutUint16(a.data[8:10], 0) // z-length 压缩列表长度
	// 将尾部填充一个0xff，表示压缩列表结束
	// z-ender
	a.data[10] = 0xff
	return a
}

// Insert 插入数据
// 采用头插法
func (a *adkZipList[T]) Push(val T) error {
	//如果链表为空直接插入数据
	zl := binary.LittleEndian.Uint16(a.data[8:10])
	zb := binary.LittleEndian.Uint32(a.data[0:4])
	zt := binary.LittleEndian.Uint32(a.data[4:8])
	parseElement, err := a.parseElement(nil, val)
	if err != nil {
		return err
	}
	// 插入数据
	a.data, _ = InsertBytes(a.data, parseElement.context, 10)
	// 插入头部
	a.data, _ = InsertBytes(a.data, parseElement.head.thisEntryLen, 10)
	a.data, _ = InsertBytes(a.data, parseElement.head.preEntryLen, 10)
	vall := uint32(len(parseElement.head.thisEntryLen)) + uint32(len(parseElement.head.preEntryLen)) + uint32(len(parseElement.context))
	if zl == 0 {
		// zl + 1
		binary.LittleEndian.PutUint16(a.data[8:10], zl+1)
		// zb 变动
		binary.LittleEndian.PutUint32(a.data[0:4], zb+vall)
		return nil
	}
	prelen := utils.IntToBytes(int(vall))
	a.data = append(append(a.data[:10+vall], prelen...), a.data[11+vall:]...)
	// 修改头部
	// zb 变动
	binary.LittleEndian.PutUint32(a.data[0:4], zb+vall+uint32(len(prelen)-1))
	binary.LittleEndian.PutUint32(a.data[4:8], zt+vall)
	binary.LittleEndian.PutUint16(a.data[8:10], zl+1)
	return nil
}

func (a *adkZipList[T]) Pop() T {
	zl := binary.LittleEndian.Uint16(a.data[8:10])
	zb := binary.LittleEndian.Uint32(a.data[0:4])
	zt := binary.LittleEndian.Uint32(a.data[4:8])
	var val T
	if zl == 0 {
		return val
	}
	mask := a.data[11] & 0xC0
	var data []byte
	bl := 0
	if mask == encodingLongLong {
		// 获取长度
		length := binary.LittleEndian.Uint32(a.data[11:15])
		// 获取数据
		data = a.data[15 : 15+length]
		binary.LittleEndian.PutUint32(a.data[0:4], zb-uint32(len(data))-6)
		if zl != 1 {
			binary.LittleEndian.PutUint32(a.data[4:8], zt-uint32(len(data))-6)
		}
		bl = 6 + int(length)
	} else if mask == encodingLong {
		length := binary.LittleEndian.Uint16(a.data[11:13])
		data = a.data[13 : 13+length]
		binary.LittleEndian.PutUint32(a.data[0:4], zb-uint32(len(data))-3)
		if zl != 1 {
			binary.LittleEndian.PutUint32(a.data[4:8], zt-uint32(len(data))-3)
		}
		bl = 3 + int(length)
	} else if mask == encodingShort {
		length := int(a.data[11])
		data = a.data[12 : 12+length]
		binary.LittleEndian.PutUint32(a.data[0:4], zb-uint32(len(data))-2)
		if zl != 1 {
			binary.LittleEndian.PutUint32(a.data[4:8], zt-uint32(len(data))-2)
		}
		bl = 2 + length
	}
	binary.LittleEndian.PutUint16(a.data[8:10], zl-1)
	if err := a.unmarshal(data, &val); err != nil {
		return val
	}
	// 删除数据
	a.data = append(a.data[:10], a.data[10+bl:]...)
	// 修改后一个元素的preEntryLen
	if zl != 1 {
		if bl >= preSizeLong && bl <= preSizeLongLong {
			a.data = append(a.data[:11], a.data[11+4:]...)
		}
		a.data[10] = 0x00
	}
	return val
}

func (a *adkZipList[T]) String() string {
	zl := binary.LittleEndian.Uint16(a.data[8:10])
	if zl == 0 {
		return ""
	}
	res := make([]byte, 0)
	zt := binary.LittleEndian.Uint32(a.data[4:8])
	preDataLen := uint32(0)
	for i := uint16(0); i < zl; i++ {
		if a.data[zt] == 0xfe {
			preDataLen = binary.LittleEndian.Uint32(a.data[zt+1:zt+5]) + 254
		} else {
			preDataLen = uint32(a.data[zt])
		}

		mask := a.data[zt+1] & 0xC0
		switch mask {
		case encodingLongLong:
			length := binary.LittleEndian.Uint32(a.data[zt : zt+5])
			res = append(res, a.data[zt+5:zt+5+length]...)
			zt -= 5 + length
		case encodingLong:
			length := binary.LittleEndian.Uint16(a.data[zt : zt+3])
			res = append(res, a.data[zt+3:zt+3+uint32(length)]...)
			zt -= 3 + uint32(length)
		case encodingShort:
			length := int(a.data[zt+1])
			res = append(res, a.data[zt+2:zt+2+uint32(length)]...)
			zt -= 2 + uint32(length)
		}
		if preDataLen == 0 {
			break
		}
		res = append(res, ',')
	}
	return string(res)
}

func (a *adkZipList[T]) IsEmpty() bool {
	return binary.LittleEndian.Uint16(a.data[8:10]) == 0
}

// InsertBytes 将 insert 字节切片插入到 target 字节切片的指定 index 位置。
// 返回一个新的字节切片，原始的 target 不会被修改。
func InsertBytes(target, insert []byte, index int) ([]byte, error) {
	// 参数验证
	if index < 0 || index > len(target) {
		return nil, errors.New("index out of bounds")
	}

	// 如果 insert 为空，直接返回 target 的副本
	if len(insert) == 0 {
		newSlice := make([]byte, len(target))
		copy(newSlice, target)
		return newSlice, nil
	}

	// 如果 target 为空，返回 insert 的副本
	if len(target) == 0 {
		newSlice := make([]byte, len(insert))
		copy(newSlice, insert)
		return newSlice, nil
	}

	// 创建一个新的切片，长度为 len(target) + len(insert)
	newSlice := make([]byte, 0, len(target)+len(insert))

	// 拷贝 target[:index]
	newSlice = append(newSlice, target[:index]...)

	// 拷贝 insert
	newSlice = append(newSlice, insert...)

	// 拷贝 target[index:]
	newSlice = append(newSlice, target[index:]...)

	return newSlice, nil
}
