package object

import (
	"testing"

	"github.com/fredrikkvalvik/temp-lang/pkg/tester"
)

func TestLenBuiltin(t *testing.T) {
	tests := []struct {
		name          string
		input         []Object
		expectedValue float64
	}{
		{
			"string",
			[]Object{&StringObj{"hello world"}},
			11,
		},
		{
			"list",
			[]Object{&ListObj{
				Values: []Object{
					&BooleanObj{},
					&BooleanObj{},
					&BooleanObj{},
					&BooleanObj{},
				},
			}},
			4,
		},
		{
			"map",
			[]Object{&MapObj{
				Pairs: map[HashKey]KeyValuePair{
					{Hash: 1}: {},
					{Hash: 2}: {},
					{Hash: 3}: {},
					{Hash: 4}: {},
				},
			}},
			4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := tester.New(t, "")

			result := LenBuiltin(tt.input...)

			tr.AssertNotNil(result, "result should not be nil")
			tr.AssertNotEqual(result.Type(), OBJ_ERROR, "result should not be of type error")
			tr.AssertEqual(result.(*NumberObj).Value, tt.expectedValue, "test length equal to expected")
		})
	}
}
