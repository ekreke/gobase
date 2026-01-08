package osutil

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// GetProcessName returns the name of the current process
func GetProcessName() string {
	return filepath.Base(os.Args[0])
}

// ExitWithError exit process with error
func ExitWithError(err error) {
	progName := filepath.Base(os.Args[0])
	fmt.Fprintf(os.Stderr, "%s exit -1: %+v\n", progName, err)
	os.Exit(-1)
}

// SIGTERMExit exit process with warning message & clean funcs
// TODO: add context for clean up funcs
func SIGTERMExit(cleanups ...func()) {
	progName := filepath.Base(os.Args[0])
	fmt.Fprintf(os.Stderr, "Warning %s receive process terminal SIGTERM exit 0\n", progName)
	for _, cleanup := range cleanups {
		if cleanup != nil {
			cleanup()
		}
	}
	fmt.Fprintf(os.Stderr, "Cleanup finished. %s exit 0\n", progName)
	os.Exit(0)
}

// GetCurrentGoroutineIDFromStack returns the current goroutine ID from the stack trace.
// Caveats:
//  1. Performance Impact: Calling runtime.Stack is expensive as it involves scanning the stack
//     and can lead to significant latency. It is NOT suitable for high-concurrency hot paths.
//  2. Fragility: It relies on the specific string format of the Go runtime stack trace.
//     Future Go versions might change this format, causing the function to fail or panic.
//  3. GC Pressure: Frequent memory allocations (byte slices) and string manipulations
//     increase garbage collection overhead.
//  4. Anti-pattern: Go purposely hides the goroutine ID to discourage its use as a key
//     for Thread-Local Storage (TLS). Use context.Context for request-scoped values instead.
//
// It is only for debugging purpose, if really need to get the goroutine ID, can use context.Context and request id to track the goroutine, or can follow the project below:
// https://github.com/petermattis/goid
func GetCurrentGoroutineIDFromStack() string {
	buf := make([]byte, 128)
	buf = buf[:runtime.Stack(buf, false)]
	stackInfo := string(buf)
	return strings.TrimSpace(strings.Split(strings.Split(stackInfo, "[running]")[0], "goroutine")[1])
}
