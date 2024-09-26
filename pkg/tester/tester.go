package tester

import (
	"testing"
)

type Tester struct {
	T    *testing.T
	name string
}

func New(t *testing.T, name string) *Tester {
	return &Tester{T: t, name: name}
}

// set the prefix for the error message
func (tt *Tester) SetName(name string) {
	tt.T.Helper()
	tt.name = name
}

func (tt *Tester) AssertEqual(in, expect any) {
	tt.T.Helper()

	if in != expect {
		tt.T.Fatalf("%s: Assert failed: %v != %v", tt.name, in, expect)
	}
}

func (tt *Tester) AssertNotNil(v any) {
	tt.T.Helper()

	if v == nil {
		tt.T.Fatalf("%s: NonNil Assert failed", tt.name)
	}
}

func (tt *Tester) AssertNil(v any) {
	tt.T.Helper()

	if v == nil {
		tt.T.Fatalf("%s: Nil Assert failed", tt.name)
	}
}
