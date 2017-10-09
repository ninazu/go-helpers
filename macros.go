package ninazu

import (
	"runtime"
	"strconv"
)

func COMPILER_FILE() string {
	_, file, _, _ := runtime.Caller(1)

	return file
}

func COMPILER_LINE() string {
	_, _, line, _ := runtime.Caller(1)

	return strconv.Itoa(line)
}
