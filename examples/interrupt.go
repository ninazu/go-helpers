package main

import (
	"os"
	"fmt"
	"syscall"
	"github.com/ninazu/go-helpers"
)

func main() {
	ninazu.WaitInterrupt(func(s os.Signal) bool {
		fmt.Println(s)

		switch s {
		case syscall.SIGUSR1, syscall.SIGUSR2:
			return false
			
		default:
			return true
		}
	}, []os.Signal{
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
