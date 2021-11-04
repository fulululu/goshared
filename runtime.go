package goshared

import (
	"runtime"
	"strings"
)

func FunctionName() string {
	pc, _, _, _ := runtime.Caller(1)
	nameChain := runtime.FuncForPC(pc).Name()
	names := strings.Split(nameChain, ".")
	return names[len(names)-1]
}
