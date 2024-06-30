package app

import (
	"fmt"
	"runtime"
	"strings"
)

func traceError(err error) {
	trace(err.Error())
}

func trace(message string) {
	var pcs [32]uintptr
	var str strings.Builder
	str.WriteString(message + "\nTraceback:")
	n := runtime.Callers(3, pcs[:])
	for _, pc := range pcs[:n] {
		fn := runtime.FuncForPC(pc)
		file, line := fn.FileLine(pc)
		str.WriteString(fmt.Sprintf("\n\t%s:%d", file, line))
	}
	fmt.Println(str.String())
}
