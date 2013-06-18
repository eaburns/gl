package ui

import (
	"runtime"
)

func init() {
	runtime.LockOSThread()
}

var (
	// Do is a channel of functions that should be executed on the main thread.
	do = make(chan func())

	// Stop signals that Start should return.
	stop = make(chan struct{})
)

// Start must be called by the main go routine to start the user interface.
// It does not return until Stop is called.
func Start() {
	for {
		select {
		case f := <-do:
			f()
		case <-stop:
			return
		}
	}
}

// Stop stops the user interface, causing the Start function to return.
func Stop() {
	stop <- struct{}{}
}

// Do schedules a function to be executed by the main go routine on the main thread;
// it returns when the function execution is completed.
func Do(f func()) {
	done := make(chan struct{}, 1)
	do <- func() {
		f()
		done <- struct{}{}
	}
	<-done
}
