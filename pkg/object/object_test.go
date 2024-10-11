package object

import (
	"fmt"
	"testing"

	"github.com/fredrikkvalvik/temp-lang/pkg/tester"
)

func TestHashes(t *testing.T) {
	tests := []struct {
		inA         Hashable
		inB         Hashable
		shouldEqual bool
	}{
		{
			&NumberObj{Value: 10.10},
			&NumberObj{Value: 10.10},
			true,
		},
		{
			&NumberObj{Value: 10.11111},
			&NumberObj{Value: 10.11111},
			true,
		},
		{
			&StringObj{Value: "Hello"},
			&StringObj{Value: "Hello"},
			true,
		},
		{
			&StringObj{Value: "10"},
			&NumberObj{Value: 10},
			false,
		},
		{
			&NumberObj{Value: 10.1},
			&NumberObj{Value: 10},
			false,
		},
		{
			&StringObj{Value: "Hello "},
			&StringObj{Value: "Hello"},
			false,
		},
		{
			&BooleanObj{Value: true},
			&BooleanObj{Value: true},
			true,
		},
		{
			&BooleanObj{Value: true},
			&BooleanObj{Value: false},
			false,
		},
	}

	for idx, tt := range tests {
		t.Run(fmt.Sprint(idx), func(t *testing.T) {
			tr := tester.New(t, "")

			eq := tt.inA.HashKey() == tt.inB.HashKey()

			tr.T.Logf("a=%v b=%v", tt.inA.HashKey(), tt.inB.HashKey())
			tr.AssertEqual(eq, tt.shouldEqual)
		})
	}
}
