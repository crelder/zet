package bl

import (
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestParseIndex(t *testing.T) {
	tcs := []struct {
		input  string
		output []Index
	}{
		{"",
			nil},
		{"Leben: 170713a",
			[]Index{{"Leben", []string{"170713a"}}}},
		{"Leben: 170713a, 201104c",
			[]Index{{"Leben", []string{"170713a", "201104c"}}}},
		{`Leben: 170713a, 201104c
				Programmierung, Objektorientiert: 220130f, 120412e`,
			[]Index{{"Leben", []string{"170713a", "201104c"}},
				{"Programmierung, Objektorientiert", []string{"220130f", "120412e"}}}},
	}

	for _, tc := range tcs {
		if out := ParseIndex(tc.input); cmp.Diff(tc.output, out) != "" {
			t.Errorf("Expected %v, got %v", tc.output, out)
		}
	}
}
