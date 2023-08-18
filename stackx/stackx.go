package stackx

import (
	"runtime"
	"strconv"
	"strings"
	"unsafe"

	"github.com/pkg/errors"
)

const stackMaxDepth = 8

type StackTracer interface {
	StackTrace() errors.StackTrace
}

// Stack represents a Stack of program counters.
type Stack []uintptr

func (s *Stack) StackTrace() errors.StackTrace {
	f := make([]errors.Frame, len(*s))
	for i := 0; i < len(f); i++ {
		f[i] = errors.Frame((*s)[i])
	}
	return f
}

// Callers returns the current call stack.
func Callers(skip int) *Stack {
	var pcs [stackMaxDepth]uintptr
	n := runtime.Callers(skip+2, pcs[:])
	var st Stack = pcs[:n]
	return &st
}

// GetStackFromError attempts to retrieve the stack trace of an error.
func GetStackFromError(err error) StackTracer {
	stackTracer := tryFindErrStackTacker(err)
	if stackTracer == nil {
		return nil
	}
	return stackTracer
}

// StackToString converts runtime.Frames to string.
func StackToString(st StackTracer) string {
	if st == nil || st.StackTrace() == nil {
		return ""
	}
	tmp := st.StackTrace()
	callersFrames := runtime.CallersFrames(*(*[]uintptr)(unsafe.Pointer(&tmp)))
	i := 0
	var sb strings.Builder
	for f, again := callersFrames.Next(); again; f, again = callersFrames.Next() {
		sb.WriteString(f.Function)
		sb.WriteByte('\n')
		sb.WriteByte('\t')
		sb.WriteString(f.File)
		sb.WriteByte(':')
		sb.WriteString(strconv.Itoa(f.Line))
		sb.WriteByte('\n')
		i++
		if i == stackMaxDepth {
			break
		}
	}
	if i == 0 {
		return ""
	}
	s := sb.String()
	return s[:len(s)-1]
}

// GetStackStringFromError returns the error stack string.
func GetStackStringFromError(err error) string {
	return StackToString(GetStackFromError(err))
}

// tryFindErrStackTacker attempts to find the earliest error that implements errStackTracer.
//
//nolint:errorlint
func tryFindErrStackTacker(err error) StackTracer {
	var st StackTracer
	for err != nil {
		v, ok := err.(StackTracer)
		if ok {
			st = v
		}
		switch x := err.(type) {
		case interface{ Unwrap() error }:
			err = x.Unwrap()
		case interface{ Cause() error }:
			err = x.Cause()
		default:
			return st
		}
	}
	return st
}
