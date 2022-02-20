package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	server "github.com/rombintu/xxxants/internal"
)

func runServer() {
	server := server.NewServer("localhost", "5000")
	if err := server.Start(); err != nil {
		server.Logger.Fatalf("%v", err)
	}
	server.Logger.Info("Server exit")
}

func programExit(ctx context.Context) {
	exitCh := make(chan os.Signal)
	signal.Notify(exitCh, os.Interrupt, syscall.SIGTERM)
	<-exitCh
	fmt.Println("Exit with 0")
	os.Exit(0)
}

func main() {
	ctx, cancelFunc := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	wg.Add(2)
	defer wg.Done()
	go programExit(ctx)
	go runServer()
	wg.Wait()
	cancelFunc()
}
