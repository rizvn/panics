# Panics
Go package for robust panic handling, recovery, and error reporting utilities. 

For many use cases panics are a better alternative to returning errors, whilst also 
adding context messages to panics and stack traces can help with debugging and understanding 
the flow of the program, especially in cases where the error is unrecoverable or indicates a 
bug in the code. 
This package provides utilities to handle panics gracefully, log them, and recover from them.


## Installation
```bash
go get github.com/rizvn/panics
```

## Features
- Panic handling utilities: `OnError`, `OnNil`, `OnFalse`, `OnBlank`, `WithTrace`
- Panic recovery utilities: `Recover`, `RecoverAndHandle`
- Retry and Try utilities: `Retry`, `Try`
- HTTP middleware for panic recovery: `RecoveryMiddleware`
- Stack trace generation for panics
- Customizable panic handling with optional messages and stack traces


## Usage
```go
import  "github.com/rizvn/panics"
```

## Functions and Usage

### OnError

Panics if `err` is not nil, including an optional message and stack trace.

```go
err := errors.New("something went wrong")
panics.OnError(err, "operation failed")
```

### OnNil

Panics if value is nil, including an optional message and stack trace.

```go
var ptr *int = nil
panics.OnNil(ptr, "pointer is nil")
```

### OnFalse

Panics if condition is false, including an optional message and stack trace.

```go
panics.OnFalse(1 > 2, "math is broken")
```

### OnBlank

Panics if the string value is blank (empty or whitespace), including an optional message and stack trace.

```go
panics.OnBlank("   ", "string is blank")
```

### WithTrace

Panics with the provided message and a stack trace.

```go
panics.WithTrace("unexpected situation")
```

### Recover

Helper to recover from panics and log the error and stack trace. This defines the panic boundary and can be placed in
the call stack.

**Example: Using `defer panics.Recover()` at the top of a goroutine, to handle panic before killing the goroutine**

```go
go func() {
    defer panics.Recover()
    // code that may panic inside goroutine
    panic("panic in goroutine")
}()
```

**Example: Using `defer panics.Recover()` at the top of a function, to handle panic before returning from the function**

```go
func doSomething() {
    defer panics.Recover()
    // code that may panic
    panic("panic in function")
}

doSomething()
```

### RecoverAndHandle

Recovers from a panic and passes the error to the provided handler function.

```go
func mayPanic() {
    panic("something went wrong")
}
defer panics.RecoverAndHandle(func(err error) {
    fmt.Println("Recovered error:", err)
})
mayPanic()
```

### Retry

Executes the provided function, retrying up to maxRetries times if it panics.

```go
panics.Retry(3, func() {
    fmt.Println("Trying...")
    panic("fail")
})
```

### Try

Executes the provided function and returns an error if it panics.

```go
err := panics.Try(func() {
    panic("fail")
})
if err != nil {
    fmt.Println("Recovered error:", err)
}
```

### RecoveryMiddleware

HTTP middleware that recovers from panics in handlers, logs the error and stack trace, and returns a 500 Internal Server
Error response.

**Example: Standard net/http**

```go
http.Handle("/", panics.RecoveryMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    panic("handler panic")
})))
```

**Example: Using with gorilla/mux router**

```go
import (
    "github.com/gorilla/mux"
    "net/http"
    "yourmodule/panics"
)

r := mux.NewRouter()
r.Use(panics.RecoveryMiddleware)
r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    panic("handler panic")
})
http.ListenAndServe(":8080", r)
```
