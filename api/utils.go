package api

import (
	"strings"
	"math"
)

// SplitAndSelect 分割并选择其中一段
// @param index  分割后数组的下标，支持负数，负数为倒序获取，超过数组的部分下标
func SplitAndSelect(s, seq string, index int) string {
	ss := strings.Split(s, seq)
	if index < 0 {
		index = len(ss) + index
	}
	index = int(math.Min(math.Max(float64(index), 0.0), float64(len(ss)-1)))
	return ss[index%len(ss)]
}
