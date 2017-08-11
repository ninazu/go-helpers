package ninazu

import (
	"os"
	"sync"
	"os/signal"
)

type CustomInterruptCallback func(os.Signal) bool

func WaitInterrupt(interruptCallback CustomInterruptCallback, signalList []os.Signal) {
	var signalSemaphore sync.WaitGroup
	signalSemaphore.Add(1)

	var signalChannel chan os.Signal
	signalChannel = make(chan os.Signal, 1)
	signal.Notify(signalChannel, signalList...)

	go func() {
		for {
			s := <-signalChannel

			if interruptCallback(s) {
				signalSemaphore.Done()
				break
			}
		}
	}()

	signalSemaphore.Wait()
}
