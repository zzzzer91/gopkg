package stackx

import (
	"runtime"
	"strconv"
	"strings"
	"unsafe"

	"github.com/pkg/errors"
)

// RecordStack returns the current call stack string.
func RecordStack(skip int) string {
	return CallersFrames2Str(GetCallersFrames(skip + 2))
}

// GetStackFromError returns the error stack string.
func GetStackFromError(err error) string {
	return CallersFrames2Str(GetCallersFramesFromError(err))
}

// CallersFrames2Str converts runtime.Frames to string.
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

// GetCallersFrames returns the current call stack.
func GetCallersFrames(skip int) *runtime.Frames {
	const maxDepth = 16
	var pcs [maxDepth]uintptr
	n := runtime.Callers(skip+1, pcs[:])
	return runtime.CallersFrames(pcs[:n])
}

// GetCallersFramesFromError attempts to retrieve the stack trace of an error.
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

// tryFindErrStackTacker attempts to find the earliest error that implements errStackTracer.
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
