package main

import (
	"reflect"
	"testing"
)

func TestSplitArgs(t *testing.T) {
	tests := []struct {
		input    string
		expected []string
	}{
		{`curl "https://api.example.com" -H "Auth: Bearer foo"`, []string{"curl", "https://api.example.com", "-H", "Auth: Bearer foo"}},
		{`-X POST http://test.com`, []string{"-X", "POST", "http://test.com"}},
		{`-d '{"foo": "bar"}'`, []string{"-d", `{"foo": "bar"}`}},
	}
	for _, tc := range tests {
		got := splitArgs(tc.input)
		if !reflect.DeepEqual(got, tc.expected) {
			t.Errorf("splitArgs(%q) = %v; want %v", tc.input, got, tc.expected)
		}
	}
}
