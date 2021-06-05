package utils

import "fmt"

// Must will panic if the error is not nil.
func Must(err error) {
	if err != nil {
		panic(fmt.Sprintf("encountered unexpected error: %v", err))
	}
}
