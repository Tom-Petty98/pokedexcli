package main

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestCleanInput(t *testing.T) {
	tests := map[string]struct {
		input string
		sep   string
		want  []string
	}{
		"simple":       {input: "a b c", want: []string{"a", "b", "c"}},
		"wrong sep":    {input: "a/b/c", want: []string{"a/b/c"}},
		"no sep":       {input: "abc", want: []string{"abc"}},
		"trailing sep": {input: " a b c ", want: []string{"a", "b", "c"}},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := cleanInput(tc.input)
			diff := cmp.Diff(tc.want, got)
			if diff != "" {
				t.Fatalf("%s", diff)
			}
		})
	}
}
