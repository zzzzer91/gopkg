package stackx

import (
	"runtime"
	"strconv"
	"strings"
	"unsafe"

	"github.com/pkg/errors"
)

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
	const depth = 16
	var pcs [depth]uintptr
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
	var sb strings.Builder
	for f, again := callersFrames.Next(); again; f, again = callersFrames.Next() {
		sb.WriteString(f.Function)
		sb.WriteByte('\n')
		sb.WriteByte('\t')
		sb.WriteString(f.File)
		sb.WriteByte(':')
		sb.WriteString(strconv.Itoa(f.Line))
		sb.WriteByte('\n')
	}
	s := sb.String()
	return s[:len(s)-1]
}

// GetStackStringFromError returns the error stack string.
func GetStackStringFromError(err error) string {
	return StackToString(GetStackFromError(err))
}

type errCauser interface {
	Cause() error
}

// tryFindErrStackTacker attempts to find the earliest error that implements errStackTracer.
func tryFindErrStackTacker(err error) StackTracer {
	var st StackTracer
	for err != nil {
		v, ok := err.(StackTracer)
		if ok {
			st = v
		}
		cause, ok := err.(errCauser)
		if !ok {
			break
		}
		err = cause.Cause()
	}
	return st
}
