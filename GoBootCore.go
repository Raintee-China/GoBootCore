package GoBootCore

import "fmt"

// Greet returns a greeting message.
func Greet(name string) string {
	return fmt.Sprintf("Hello, %s from GoBootCore!", name)
}

// Add returns the sum of two integers.
func Add(a, b int) int {
	return a + b
}
