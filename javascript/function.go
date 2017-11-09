package javascript

import (
	"github.com/robertkrimen/otto"
	"strings"
	"strconv"
	"math"
)

func AddFunction(vm *otto.Otto) {
	funcs := map[string]func(otto.FunctionCall) otto.Value{
		"avg":          avg,
		"maxNextDiff":  maxNextDiff,
		"maxNextMulti": maxNextMulti,
	}
	for k, v := range funcs {
		vm.Set(k, v)
	}
}

// avg 获取平均数
// `avg([1,2,3,4,5,6]) = 3.5`
func avg(call otto.FunctionCall) otto.Value {
	numStrs := strings.Split(call.Argument(0).String(), ",")
	len := 0.0
	count := 0.0
	for _, v := range numStrs {
		num, err := strconv.ParseFloat(v, 64)
		if err != nil {
			continue
		}
		len++
		count += num
	}
	result, _ := otto.ToValue(count / len)
	return result
}

// maxNextDiff 间隔最大差值
// `maxNextDiff([1,4,5,14,25,46]) = 21`
func maxNextDiff(call otto.FunctionCall) otto.Value {
	nums := stringToFloatSlice(call.Argument(0).String())
	max := 0.0
	for i, v := range nums {
		if i == len(nums)-1 {
			break
		}
		next := nums[i+1]
		diff := math.Abs(v - next)
		if diff > max {
			max = diff
		}
	}
	res, _ := otto.ToValue(max)
	return res
}

// maxNextMulti 间隔最大相差倍数
// `maxNextMulti([1,4,5,14,25,46]) = 4`
func maxNextMulti(call otto.FunctionCall) otto.Value {
	nums := stringToFloatSlice(call.Argument(0).String())
	max := 0.0
	for i, v := range nums {
		if i == len(nums)-1 {
			break
		}
		next := nums[i+1]
		multi := v / next
		if multi < 1.0 {
			multi = 1.0 / multi
		}
		if multi > max {
			max = multi
		}
	}
	res, _ := otto.ToValue(max)
	return res
}
