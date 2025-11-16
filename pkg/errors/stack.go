package errors

import (
	"runtime"
	"strconv"
	"strings"
)

type StackTrace []Frame

type Frame uintptr

const maxDepth = 32

func trace(skip int) StackTrace {
	var programCounters [maxDepth]uintptr

	n := runtime.Callers(skip+2, programCounters[:])

	stackFrames := make([]Frame, 0, n)

	for _, programCounter := range programCounters {
		stackFrames = append(stackFrames, Frame(programCounter))
	}

	return stackFrames
}

func (s StackTrace) String() string {
	sb := strings.Builder{}

	for index, frame := range s {
		if index != 0 {
			_, _ = sb.WriteString("\n")
		}

		_, _ = sb.WriteString(frame.Func())
		_, _ = sb.WriteString("\n\t")
		_, _ = sb.WriteString(frame.File())
		_, _ = sb.WriteString(":")
		_, _ = sb.WriteString(strconv.Itoa(frame.Line()))
	}

	return sb.String()
}

func (f Frame) Func() string {
	fn := runtime.FuncForPC(f.pc())
	if fn == nil {
		return ""
	}

	return fn.Name()
}

func (f Frame) File() string {
	fn := runtime.FuncForPC(f.pc())
	if fn == nil {
		return ""
	}

	file, _ := fn.FileLine(f.pc())

	return file
}

func (f Frame) Line() int {
	fn := runtime.FuncForPC(f.pc())
	if fn == nil {
		return 0
	}

	_, line := fn.FileLine(f.pc())

	return line
}

func (f Frame) pc() uintptr {
	return uintptr(f) - 1
}
