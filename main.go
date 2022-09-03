package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	// graceful shutdown
	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)

	fmt.Println("start!")
	g := sync.WaitGroup{}
	g.Add(2)
	// a gorutine do process
	go func() {
		fmt.Println("start")
		times := 0
		for {
			select {
			// if receive done signal, cancel
			case <-ctx.Done():
				fmt.Println("end")
				g.Done()
				return
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
	}()
	// a gorutine receive signal
	go func() {
		<-ctx.Done()
		fmt.Printf("stop signal received, waiting for process done.\n")
		// should not wait after
		g.Done()
		ctx2, _ := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		for {
			select {
			case <-ctx2.Done():
				fmt.Printf("stop signal received again, force shutting down.\n")
				os.Exit(1)
			default:
				time.Sleep(time.Millisecond * 500)
			}
		}
	}()

	g.Wait()
	fmt.Println("end!")
}
