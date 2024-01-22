package compiler

import (
	"fmt"
	"runtime"
)

// TracedError is a custom error type that includes a slice of error messages
type TracedError struct {
	Trace []string
}

func (e *TracedError) Error() string {
	return fmt.Sprintf("Trace: %v", e.Trace)
}

// TraceError logs the function name, file, and line number
func TraceError(err error) error {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)

		var trace []string
		if tracedError, ok := err.(*TracedError); ok {
			// If the error is already a TracedError, append to its existing trace
			trace = append(tracedError.Trace, fmt.Sprintf("[error] at %s:%d\n", file, line))
		} else {
			// Otherwise, create a new trace
			trace = append(trace, fmt.Sprintf("[error] at %s:%d: %v\n", file, line, err))
		}

		return &TracedError{Trace: trace}
	}
	return nil
}

// InvalidAccessError is returned when the caller attempted to access on type of lexical element when the tokenizer was at another
type InvalidAccessError struct {
	wasAttempted TokenType // Type the caller attempted to access
	wasValid     TokenType // Type the tokenizer was able to access
	wasValidVal  string    // The value of the valid token
}

func (e *InvalidAccessError) Error() string {
	return fmt.Sprintf("expected a token of type `%v` but found a token \"%v\" of type `%v` instead", e.wasAttempted, e.wasValidVal, e.wasValid)
}
