package myutils

import (
	"fmt"
	"testing"
)

func TestSampleFunction(t *testing.T) {
	expected := "Hello from myutils!"
	got := SampleFunction()
	if got != expected {
		t.Errorf("SampleFunction() = %q, want %q", got, expected)
	}
	fmt.Println(got)
}
