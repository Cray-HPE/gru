package action

import (
	"fmt"
	"os"
	"sync"
)

// Send sends a command to a host and returns an error
func Send(args []string, f func(host string) error) {
	// Concurrently query hosts
	var wg sync.WaitGroup
	// create a buffered channel to collect errors returned by the function
	// the channel's capacity is equal to the length of the slice to ensure that the goroutines won't block when sending errors to the channel
	errorCh := make(chan error, len(args))

	for _, h := range args {
		wg.Add(1)
		go func(host string) {
			defer wg.Done()
			// if adding a blade failes, it returns an error
			err := f(host)
			// if there is an error, send it to the errorCh channel
			if err != nil {
				errorCh <- fmt.Errorf("[%s] %v", host, err)
			} else {
				fmt.Printf("[%s]: command sent\n", host)
			}
		}(h)

	}

	// After waiting for all goroutines to complete,
	wg.Wait()
	// close the errorCh channel
	close(errorCh)

	// and then iterate over the channel to report the errors
	for err := range errorCh {
		fmt.Printf("%s\n", err)
	}
	if len(errorCh) > 0 {
		os.Exit(1)
	}
}

// Get sends a command to a host and returns a string and an error
func Get(args []string, f func(host string) (string, error)) {
	// Concurrently query hosts
	var wg sync.WaitGroup
	// create a buffered channel to collect errors returned by the function
	// the channel's capacity is equal to the length of the slice to ensure that the goroutines won't block when sending errors to the channel
	errorCh := make(chan error, len(args))

	for _, h := range args {
		wg.Add(1)
		go func(host string) {
			defer wg.Done()
			// if adding a blade failes, it returns an error
			resp, err := f(host)
			// if there is an error, send it to the errorCh channel
			if err != nil {
				errorCh <- fmt.Errorf("[%s] %v", host, err)
			} else {
				fmt.Printf("[%s]: %s\n", host, resp)
			}
		}(h)

	}

	// After waiting for all goroutines to complete,
	wg.Wait()
	// close the errorCh channel
	close(errorCh)

	// and then iterate over the channel to report the errors
	for err := range errorCh {
		fmt.Printf("%s\n", err)
	}
	if len(errorCh) > 0 {
		os.Exit(1)
	}
}
