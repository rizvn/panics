package panics

import (
	"fmt"
	"log/slog"
	"net/http"
	"runtime"
	"runtime/debug"
	"strings"
)

// OnError panics if err is not nil, including an optional message and stack trace.
func OnError(err error, message string) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		panic(fmt.Sprintf("(%s:%d): %s %v", file, line, message, err))
	}
}

// OnNil panics if value is nil, including an optional message and stack trace.
func OnNil(value any, message string) {
	if value == nil {
		_, file, line, _ := runtime.Caller(1)
		panic(fmt.Sprintf("(%s:%d): %s %v", file, line, message, "nil value"))
	}
}

// OnFalse panics if condition is false, including an optional message and stack trace.
func OnFalse(condition bool, message string) {
	if !condition {
		_, file, line, _ := runtime.Caller(1)
		panic(fmt.Sprintf("(%s:%d): %s %v", file, line, message, ""))
	}
}

// OnBlank panics if the string value is blank (empty or whitespace), including an optional message and stack trace.
func OnBlank(value string, message string) {
	if strings.TrimSpace(value) == "" {
		_, file, line, _ := runtime.Caller(1)
		panic(fmt.Sprintf("(%s:%d): %s %v.", file, line, message, "blank string"))
	}
}

// WithTrace panics with the provided message and a stack trace.
func WithTrace(message string) {
	panic(fmt.Errorf("panic: %s\nStacktrace: %s\n---", message, debug.Stack()))
}

// Recover is a helper to recover from panics and log the error and stack trace.
func Recover() {
	if r := recover(); r != nil {
		if err, ok := r.(error); ok {
			slog.Error("Recovered from panic: %w", err)
		} else {
			slog.Error("recovered from panic: %v", r)
		}
	}
}

// RecoverAndHandle recovers from a panic and passes the error to the provided handler function.
//
// Example usage:
//
//	func mayPanic() {
//	    panic("something went wrong")
//	}
//	func main() {
//	    defer RecoverAndHandle(func(err error) {
//	        fmt.Println("Recovered error:", err)
//	    })
//	    mayPanic()
//	}
func RecoverAndHandle(fn func(err error)) {
	if r := recover(); r != nil {
		if err, ok := r.(error); ok {
			fn(err)
		} else {
			fn(fmt.Errorf("recovered from panic: %v", r))
		}
	}
}

// Retry executes the provided function, retrying up to maxRetries times if it panics.
//
// Example usage:
//
//	Retry(3, func() {
//	    // code that may panic
//	    fmt.Println("Trying...")
//	    panic("fail")
//	})
func Retry(maxRetries int, fn func()) {
	retries := maxRetries
	for retries > 0 {
		err := Try(fn)

		if err == nil {
			return
		}

		slog.Error("Retrying function due to error: %v", err)
		retries--
	}
}

// Try executes the provided function and returns an error if it panics.
// It uses RecoverAndHandle to capture any panic as an error.
//
// Example usage:
//
//	err := Try(func() {
//	    // code that may panic
//	    panic("fail")
//	})
//	if err != nil {
//	    fmt.Println("Recovered error:", err)
//	}
func Try(fn func()) error {
	var err error
	func() {
		// Use a deferred function to recover from panic and handle retries
		defer RecoverAndHandle(func(err2 error) {
			err = err2
		})

		fn()
	}()
	return err
}

// RecoveryMiddleware is an HTTP middleware that recovers from panics in handlers,
// logs the error and stack trace, and returns a 500 Internal Server Error response.
func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				if err, ok := rec.(error); ok {
					slog.Error("recovered from panic: %w", err)
				} else {
					slog.Error("recovered from panic: %v", r)
				}
				slog.Error("Stacktrace: %s\n", debug.Stack())
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
