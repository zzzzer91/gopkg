package stackx

import (
	"runtime"
	"strconv"
	"strings"
	"unsafe"

	"github.com/pkg/errors"
)

// RecordStack return current stack string.
func RecordStack(skip int) string {
	return CallersFrames2Str(GetCallersFrames(skip + 2))
}

// GetStackFromError return current error stack string.
func GetStackFromError(err error) string {
	return CallersFrames2Str(GetCallersFramesFromError(err))
}

func CallersFrames2Str(callersFrames *runtime.Frames) string {
	if callersFrames == nil {
		return ""
	}
	var sb strings.Builder
	for f, again := callersFrames.Next(); again; f, again = callersFrames.Next() {
		sb.WriteString(f.Function)
		sb.WriteByte(':')
		sb.WriteString(strconv.Itoa(f.Line))
		sb.WriteByte('\n')
	}
	s := sb.String()
	return s[:len(s)-1]
}

// GetCallersFrames returns current stack.
func GetCallersFrames(skip int) *runtime.Frames {
	const maxDepth = 16
	var pcs [maxDepth]uintptr
	n := runtime.Callers(skip+1, pcs[:])
	return runtime.CallersFrames(pcs[:n])
}

// GetCallersFramesFromError try to get error's stack.
func GetCallersFramesFromError(err error) *runtime.Frames {
	stackTracer := tryFindErrStackTacker(err)
	if stackTracer == nil {
		return nil
	}
	st := stackTracer.StackTrace()
	return runtime.CallersFrames(*(*[]uintptr)(unsafe.Pointer(&st)))
}

type errStackTracer interface {
	StackTrace() errors.StackTrace
}

type errCauser interface {
	Cause() error
}

// tryFindErrStackTacker try to find last err that implements errStackTracer.
func tryFindErrStackTacker(err error) errStackTracer {
	var st errStackTracer
	for err != nil {
		v, ok := err.(errStackTracer)
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
