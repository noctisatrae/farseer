package main

import (
	"sync"
	"testing"
	// "time"
)

func TestGracefulShutdown(t *testing.T) {
	var wg sync.WaitGroup
	stopCh := make(chan struct{})

	// Start the gRPC server in a separate goroutine
	wg.Add(1)
	go Start(&wg, stopCh)

	// time.Sleep(time.Second)

	// close(stopCh)

	// Wait for the server to stop
	wg.Wait()
}