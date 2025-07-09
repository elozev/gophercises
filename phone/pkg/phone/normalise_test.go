package normalise

import (
	"testing"
)

var errorFormat = "expected: %s; got: %s"

func failIfDifferent(t *testing.T, expected any, got any) {
	if expected != got {
		t.Errorf(errorFormat, expected, got)
	}
}

func TestCleanRemoveCharacters(t *testing.T) {
	testCases := []struct {
		input string
		want  string
	}{
		{input: "", want: ""},
		{input: "(123) 456-7890", want: "1234567890"},
		{input: "+1 (555) 123-4567", want: "0015551234567"},
		{input: "    867-5309    ", want: "8675309"},
		{input: "(800) FLOWERS", want: "800"},
		{input: "+44.123.456.7890", want: "00441234567890"},
		{input: "1-800-GOFEDEX", want: "1800"},
		{input: "123-ABC-4567", want: "1234567"},
		{input: "  +1 (123) 456-7890  ", want: "0011234567890"},
		{input: "12.34.56.78.90", want: "1234567890"},
		{input: "+86 123 456 7890", want: "00861234567890"},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			actual := Clean(tc.input)
			failIfDifferent(t, tc.want, actual)
		})
	}
}

func TestRemovePlusAtStartPlusToDoubleZeros(t *testing.T) {
	testCases := []struct {
		input string
		want  string
	}{
		{input: "+44 77 88 123 123", want: "0044 77 88 123 123"},
		{input: "+1 555 123 4567", want: "001 555 123 4567"},
		{input: "+33 1 23 45 67 89", want: "0033 1 23 45 67 89"},
		{input: "+81 3 1234 5678", want: "0081 3 1234 5678"},
		{input: "+86 10 6988 6543", want: "0086 10 6988 6543"},
		{input: "+61 2 1234 5678", want: "0061 2 1234 5678"},
		{input: "+49 30 1234567", want: "0049 30 1234567"},
		{input: "+7 495 123 4567", want: "007 495 123 4567"},
		{input: "+55 11 1234 5678", want: "0055 11 1234 5678"},
		{input: "+91 22 1234 5678", want: "0091 22 1234 5678"},
		{input: "+phone", want: "00phone"},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			actual := removePlusAtStart(tc.input)
			failIfDifferent(t, tc.want, actual)
		})
	}
}
