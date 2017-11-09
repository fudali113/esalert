package javascript

import (
	"github.com/robertkrimen/otto"
	"strings"
	"strconv"
)

func GetVM() *otto.Otto {
	vm := otto.New()
	AddFunction(vm)
	return vm
}

func stringToFloatSlice(s string) []float64 {
	numStrs := strings.Split(s, ",")
	nums := make([]float64, 0, len(numStrs))
	for _, v := range numStrs {
		num, err := strconv.ParseFloat(v, 64)
		if err != nil {
			continue
		}
		nums = append(nums, num)
	}
	return nums
}
