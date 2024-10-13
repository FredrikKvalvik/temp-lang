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

func (tt *Tester) AssertEqual(in, expect any, msg ...string) {
	tt.T.Helper()

	if in != expect {
		tt.T.Fatalf("%s: Assert failed: %v != %v", msg, in, expect)
	}
}
func (tt *Tester) AssertNotEqual(in, expect any, msg ...string) {
	tt.T.Helper()

	if in == expect {
		tt.T.Fatalf("%s: Assert failed: %v == %v", msg, in, expect)
	}
}
func (tt *Tester) AssertTrue(in bool, msg ...string) {
	tt.T.Helper()

	if !in {
		tt.T.Fatalf("%s: AssertTrue failed, got=%v", msg, in)
	}
}

func (tt *Tester) AssertNotNil(v any, msg ...string) {
	tt.T.Helper()

	if v == nil {
		tt.T.Fatalf("%s: NonNil Assert failed", msg)
	}
}

func (tt *Tester) AssertNil(v any, msg ...string) {
	tt.T.Helper()

	if v != nil {
		tt.T.Fatalf("%s: Nil Assert failed, got=%+v", msg, v)
	}
}
