package mo2errors_test

import (
	"fmt"
	"mo2/mo2utils/mo2errors"
	"testing"
)

func TestNewEqual(t *testing.T) {
	// Different allocations should not be equal.
	if mo2errors.New(mo2errors.Mo2LengthRequired, "abc") == mo2errors.New(mo2errors.Mo2LengthRequired, "abc") {
		t.Errorf(`New("abc") == New("abc")`)
	}
	if mo2errors.New(mo2errors.Mo2LengthRequired, "abc") == mo2errors.New(mo2errors.Mo2LengthRequired, "xyz") {
		t.Errorf(`New("abc") == New("xyz")`)
	}

	// Same allocation should be equal to itself (not crash).
	err := mo2errors.New(mo2errors.Mo2LengthRequired, "jkl")
	if err != err {
		t.Errorf(`err != err`)
	}
}

func TestErrorMethod(t *testing.T) {
	err := mo2errors.New(mo2errors.Mo2LengthRequired, "abc")
	if err.Error() != "411: abc" {
		t.Errorf(`New("abc").Error() = %q, want %q`, err.Error(), "abc")
	}
}

func ExampleNew() {
	err := mo2errors.New(mo2errors.Mo2LengthRequired, "abc")
	if err != nil {
		fmt.Print(err)
	}
	// Output: 411: abc
}

func ExampleNewCode() {
	err := mo2errors.NewCode(mo2errors.Mo2LengthRequired)
	if err != nil {
		fmt.Print(err)
	}
	// Output: 411: length not match required
}
