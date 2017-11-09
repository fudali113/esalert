package javascript

import (
	"testing"
	"github.com/robertkrimen/otto"
)

func Test_avg(t *testing.T) {
	vm := otto.New()
	vm.Set("avg", avg)
	v, _ := vm.Run(`
		avg([1,2,3,4,5,6,7,8,9])
	`)
	if res, _ := v.ToInteger(); res != 5 {
		t.Error("有错")
	}
}

func Test_maxNextDiff(t *testing.T) {
	vm := otto.New()
	vm.Set("maxNextDiff", maxNextDiff)
	v, _ := vm.Run(`
		maxNextDiff([1,2,31,4,5,6,7,8,9])
	`)
	if res, _ := v.ToInteger(); res != 29 {
		t.Error("有错")
	}
}

func Test_maxNextMulti(t *testing.T) {
	vm := otto.New()
	vm.Set("maxNextMulti", maxNextMulti)
	v, _ := vm.Run(`
		maxNextMulti([1,2,3,4,5,6,7,8,9])
	`)
	if res, _ := v.ToInteger(); res != 2 {
		t.Error("有错")
	}
}
