package parse

import (
	"github.com/crelder/zet"
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestParseIndex(t *testing.T) {
	tcs := []struct {
		input       string
		index       zet.Index
		firstErrMsg string
	}{
		{"",
			nil,
			"parse Index: index is empty"},

		// Simple index with one entry "Leben", which has only one reference id.
		{"Leben: 170713a",
			map[string][]string{"Leben": {"170713a"}},
			""},

		// Simple index with one entry "Leben", which has two reference ids.
		{"Leben: 170713a, 201104c",
			map[string][]string{"Leben": {"170713a", "201104c"}},
			""},

		// Index with two entries "Leben" and "Programmierung, Objektorientiert",
		// both with two reference ids.
		{`Leben: 170713a, 201104c
  				Programmierung, Objektorientiert: 220130f, 120412e`,
			map[string][]string{"Leben": {"170713a", "201104c"},
				"Programmierung, Objektorientiert": {"220130f", "120412e"}},
			"",
		},

		// Not providing any space should work also.
		{"Leben:170713a,210404d",
			map[string][]string{"Leben": {"170713a", "210404d"}},
			""},

		// Invalid id provided should return an error.
		{"Leben:170a",
			nil,
			"index: could not parse line 0, not an id \"170a\""},

		// Wrong format due to two columns ("::"), which should return an error.
		{"Leben::170713a",
			nil,
			"index: could not parse line \"Leben::170713a\""},

		// Wrong format due to missing column (":"), which should return an error.
		{"Leben 170713a",
			nil,
			"index: could not parse line \"Leben 170713a\""},

		// No ids provided, should return an error.
		{"Leben:",
			nil,
			"index: could not parse line \"Leben:\", no ids provided"},
	}

	for _, tc := range tcs {
		got, err := Index(tc.input)
		if diff := cmp.Diff(tc.index, got); diff != "" {
			t.Errorf(diff)
		}

		var errMsg string
		if err != nil {
			errMsg = err[0].Error()
		}

		if errMsg != tc.firstErrMsg {
			t.Errorf("Got %q, wanted %q", errMsg, tc.firstErrMsg)
		}
	}
}
