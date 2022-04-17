package main

import (
	"testing"
	"unicode/utf8"
)

// Simple unit test to be used against a set of given test cases.
func TestReverse(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"Testing <Hello, world>", args{"Hello, world"}, "dlrow ,olleH"},
		{"Testing empty string", args{" "}, " "},
		{"Testing <!12345>", args{"!12345"}, "54321!"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := Reverse(tt.args.s); got != tt.want {
				t.Errorf("Reverse() = %v, want %v", got, tt.want)
			}
		})
	}
}

/* Fuzzing test to try against a random, unpredictable inputs. It's important to note that we don't have control over the inputs,
   meaning that, as a result, we cannot predict the expected output as we did previously in TestReverse.
   However, there are a few properties of the Reverse function that you can verify in a fuzz test. The two properties being checked in this fuzz test are:
		1. Reversing a string twice preserves the original value
		2. The reversed string preserves its state as valid UTF-8.
*/
func FuzzReverse(f *testing.F) {
	testcases := []string{"Hello, world", " ", "!12345"}
	for _, tc := range testcases {
		f.Add(tc) // Use f.Add to provide a seed corpus
	}
	f.Fuzz(func(t *testing.T, orig string) {
		rev, err1 := Reverse(orig)
		// Checking returning errors in case of a non utf8 valid string
		if err1 != nil {
			t.Skip() // Using t.Skip to stop the execution of this fuzz input, rather than return.
		}
		// Checking returning errors in case of a non utf8 valid string
		doubleRev, err2 := Reverse(rev)
		if err2 != nil {
			t.Skip()
		}
		if orig != doubleRev {
			t.Errorf("Before: %q, after: %q", orig, doubleRev)
		}
		if utf8.ValidString(orig) && !utf8.ValidString(rev) {
			t.Errorf("Reverse produced invalid UTF-8 string %q", rev)
		}
	})
}
