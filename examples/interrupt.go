package examples

import (
	"os"
	"fmt"
	"syscall"
)

func main() {
	ninazu.WaitInterrupt(interruptCallback, []os.Signal{
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
		syscall.SIGKILL,
		syscall.SIGSTOP,
		syscall.SIGHUP,
		syscall.SIGUSR1,
		syscall.SIGUSR2,
	})
}

func interruptCallback(s os.Signal) bool {
	fmt.Println(s)

	switch s {
	case syscall.SIGUSR1, syscall.SIGUSR2:
		return false
	default:
		return true
	}
}
