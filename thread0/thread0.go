// Package thread0 package executes functions, guaranteed to be executed on the initial thread.
package thread0

import (
	"runtime"
)

func init() {
	runtime.LockOSThread()
}

// Do is a channel of functions that should be executed on the main thread.
var do = make(chan func(), 10)

// Hijack must be called by the main go routine, usually in the main function
// of the program, to take control of the initial thread.  It never returns.
func Hijack() {
	for f := range do {
		f()
	}
}

// Do schedules a function for execution on the initial thread and waits for its
// execution to complete before returning.
func Do(f func()) {
	done := make(chan struct{}, 1)
	do <- func() {
		f()
		done <- struct{}{}
	}
	<-done
}

// DoAsync schedules a function to be executed on the initial thread.  DoAsync
// may return immediately, before the function has been executed.
func DoAsync(f func()) {
	do <- f
}
