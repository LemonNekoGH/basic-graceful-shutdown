package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"
)

func main() {
	signalChan := make(chan os.Signal, 8)
	// graceful shutdown
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	defer close(signalChan)

	fmt.Println("start!")
	g := errgroup.Group{}
	signalReceived := false

	// a gorutine do process
	g.Go(func() error {
		fmt.Println("start")
		times := 0
		for !signalReceived {
			fmt.Printf("process start %d\n", times+1)
			for i := 0; i < 10; i++ {
				time.Sleep(time.Second)
				fmt.Printf("small process done %d\n", i+1)
			}
			fmt.Printf("process done %d\n", times+1)
			times++
		}
		fmt.Println("end")
		return nil
	})
	// a gorutine receive signal
	g.Go(func() error {
		sig := <-signalChan
		signalReceived = true
		fmt.Printf("signal received %d, waiting for process done.\n", sig)
		return fmt.Errorf("times gorutine stopped")
	})

	_ = g.Wait()
	fmt.Println("end!")
}
