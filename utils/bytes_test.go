package utils

import "testing"

/*
 * 说明：
 * 作者：吕元龙
 * 时间 2024/8/30 21:23
 */

func TestBytes(t *testing.T) {
	data := []byte{0x08, 0x01}
	binary := FormatBytesAsBinary(data)
	bytes, _ := ShiftBytes(data, 0, 1)
	shiftBytes, _ := ShiftBytes(data, 1, 1)
	t.Logf("%v", binary)
	t.Logf("%v", FormatBytesAsBinary(bytes))
	t.Logf("%v", FormatBytesAsBinary(shiftBytes))
}
