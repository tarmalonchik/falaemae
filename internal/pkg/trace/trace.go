package trace

import (
	"errors"
	"fmt"
	"runtime"
	"strings"
)

func FuncName() error {
	pc := make([]uintptr, 15)
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()
	return errors.New(getFunctionLast(frame.Function))
}

func FuncNameWithError(err error) error {
	pc := make([]uintptr, 15)
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()
	if err == nil {
		return errors.New(getFunctionLast(frame.Function))
	}
	return fmt.Errorf("%s: %w", getFunctionLast(frame.Function), err)
}

func FuncNameWithErrorMsg(err error, message string) error {
	pc := make([]uintptr, 15)
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()
	outMsg := fmt.Sprintf("%s %s", getFunctionLast(frame.Function), message)
	if err == nil {
		return errors.New(outMsg)
	}
	return fmt.Errorf("%s: %w", outMsg, err)
}

func getFunctionLast(in string) string {
	split := strings.Split(in, "/")
	return split[len(split)-1]
}
