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

func PATH_OF_BINARY() string {
	path, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	
	return path
}

func PATH_OF_SOURCE() string {
	_, file, _, _ := runtime.Caller(1)
	path, _ := filepath.Abs(filepath.Dir(file))
	
	return path
}
