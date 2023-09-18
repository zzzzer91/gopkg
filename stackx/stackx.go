package stackx

import (
	"runtime"
	"strconv"
	"strings"
	"unsafe"

	"github.com/pkg/errors"
)

const stackMaxDepth = 7

type StackTracer interface {
	StackTrace() errors.StackTrace
}

// Stack represents a Stack of program counters.
type Stack []uintptr

func (s *Stack) StackTrace() errors.StackTrace {
	if s == nil {
		return nil
	}
	return *(*errors.StackTrace)(unsafe.Pointer(s))
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
	return tryFindErrStackTacker(err)
}

// StackToString converts runtime.Frames to string.
// The format is as follows:
// (*T).f | xxx/stackx/stackx_test.go:20
// f1 | xxx/stackx/stackx_test.go:25
// f2 | xxx/stackx/stackx_test.go:29
func StackToString(stackTracer StackTracer) string {
	if stackTracer == nil {
		return ""
	}
	st := stackTracer.StackTrace()
	if st == nil {
		return ""
	}
	callersFrames := runtime.CallersFrames(*(*[]uintptr)(unsafe.Pointer(&st)))
	i := 0
	var sb strings.Builder
	for f, again := callersFrames.Next(); again; f, again = callersFrames.Next() {
		sb.WriteString(formatFunction(f.Function))
		sb.WriteString(" | ")
		sb.WriteString(f.File)
		sb.WriteString(":")
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

// GetStackStringFromError returns the string of error stack.
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

// formatFunction returns, if possible, the name of the formatFunction.
func formatFunction(name string) string {
	const (
		dunno     = "???"
		centerDot = "·"
		dot       = "."
		slash     = "/"
	)

	// The name includes the path name to the package, which is unnecessary
	// since the file name is already included.  Plus, it has center dots.
	// That is, we see
	//	runtime/debug.*T·ptrmethod
	// and want
	//	*T.ptrmethod
	// Also the package path might contain dot (e.g. code.google.com/...),
	// so first eliminate the path prefix
	if lastSlash := strings.LastIndex(name, slash); lastSlash >= 0 {
		name = name[lastSlash+1:]
	}
	if period := strings.Index(name, dot); period >= 0 {
		name = name[period+1:]
	}
	name = strings.ReplaceAll(name, centerDot, dot)
	return name
}
