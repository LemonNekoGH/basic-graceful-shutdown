package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"
)

func main() {
	// graceful shutdown
	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)

	fmt.Println("start!")
	g := errgroup.Group{}

	// a gorutine do process
	g.Go(func() error {
		fmt.Println("start")
		times := 0
		for {
			select {
			// if receive done signal, cancel
			case <-ctx.Done():
				fmt.Println("end")
				return nil
			// if not receive, continue do process
			default:
				fmt.Printf("process start %d\n", times+1)
				for i := 0; i < 10; i++ {
					time.Sleep(time.Second)
					fmt.Printf("child process done %d\n", i+1)
				}
				fmt.Printf("process done %d\n", times+1)
				times++
			}
		}
	})
	// a gorutine receive signal
	g.Go(func() error {
		<-ctx.Done()
		fmt.Printf("stop signal received, waiting for process done.\n")
		return fmt.Errorf("times gorutine stopped")
	})

	_ = g.Wait()
	fmt.Println("end!")
}
