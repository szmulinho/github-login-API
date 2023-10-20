package utils

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func wait(ctx context.Context, wg *sync.WaitGroup) {
	<-ctx.Done()
	timeout := time.Minute
	start := time.Now()
	waited := make(chan struct{})
	go func() {
		defer close(waited)
		wg.Wait()
	}()

	select {
	case <-waited:
		log.Printf("shutdown completed gracefully in %s", time.Since(start))
	case <-time.After(timeout):
		log.Fatalf("shutdown exceeded wait time of %s - exiting immediately", timeout)
	}
}

type WaitFunc func()

func Gracefully() (context.Context, *sync.WaitGroup, WaitFunc) {
	ctx, cancel := context.WithCancel(context.Background())
	wg := new(sync.WaitGroup)
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigCh
		log.Printf("received signal '%s', beginning shutdown...", sig.String())
		cancel()
		sig = <-sigCh
		log.Fatalf("received signal '%s' during shutdown - exiting immediately", sig.String())
	}()
	return ctx, wg, func() { wait(ctx, wg) }
}
