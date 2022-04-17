package main

import (
	"errors"
	"fmt"
	"unicode/utf8"
)

// Reverse I hate this function, it's hard to understand.
// Note: In GoLand, press alt+insert and click on "Test on function" to automatically create a test file for this function
func Reverse(s string) (string, error) {
	// Adding some error-protection for non utf8 valid strings
	if !utf8.ValidString(s) {
		return s, errors.New("input is not valid UTF-8")
	}
	r := []rune(s)
	for i, j := 0, len(r)-1; i < len(r)/2; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r), nil
}

func main() {
	input := "The quick brown fox jumped over the lazy dog"
	rev, revErr := Reverse(input)
	doubleRev, doubleRevErr := Reverse(rev)
	fmt.Printf("original: %q\n", input)
	fmt.Printf("reversed: %q, err: %v\n", rev, revErr)
	fmt.Printf("reversed again: %q, err: %v\n", doubleRev, doubleRevErr)
}
