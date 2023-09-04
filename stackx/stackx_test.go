package stackx

import (
	"fmt"
	"testing"
)

func Test_Caller(t *testing.T) {
	s := f9()
	if len(*s) != stackMaxDepth {
		t.Errorf("length of stack is not %d", stackMaxDepth)
	}
	fmt.Println(StackToString(s))
}

type T struct {
}

func (t *T) f() *Stack {
	return Callers(0)
}

func f1() *Stack {
	t := T{}
	return t.f()
}

func f2() *Stack {
	return f1()
}

func f3() *Stack {
	return f2()
}

func f4() *Stack {
	return f3()
}

func f5() *Stack {
	return f4()
}

func f6() *Stack {
	return f5()
}

func f7() *Stack {
	return f6()
}

func f8() *Stack {
	return f7()
}

func f9() *Stack {
	return f8()
}
