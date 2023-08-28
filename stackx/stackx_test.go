package stackx

import (
	"fmt"
	"testing"
)

func Test_Caller(t *testing.T) {
	f4()
}

func f1() {
	func() {
		fmt.Println(StackToString(Callers(0)))
	}()
}

func f2() {
	f1()
}

func f3() {
	f2()
}

func f4() {
	f3()
}
