package bl

import (
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestParseFilename(t *testing.T) {
	var tcs = []struct {
		filename string
		zettel   Zettel
		errMsg   string
	}{
		// There is nothing to parse (empty string)
		{"", Zettel{}, "could not parse empty string"},

		// This id is not valid
		{"1923a - Evolution, Lego bauen", Zettel{}, "could not parse id"},
		{"210404Sof - Software.txt", Zettel{}, "could not parse id"},

		// The id is not at the start
		{"Evolution, Lego bauen - 170402a", Zettel{}, "could not parse id"},

		//
		{"200110d.txt", Zettel{"200110d", nil, nil, nil, nil, "200110d.txt"}, ""},

		// A correct, simple example
		{
			"170712a - Evolution, Lego bauen, Perfektion.txt",
			Zettel{
				"170712a",
				[]string{"Evolution", "Lego bauen", "Perfektion"},
				nil,
				nil,
				nil,
				"170712a - Evolution, Lego bauen, Perfektion.txt",
			},
			"",
		},

		// A correct example with all possibilities to parse
		// Should also work with no space between comma and the keyword, e.g. 'Evolution,Lego bauen'
		{
			"170712a - Evolution,Lego bauen, Perfektion - Gespräch Peter, nick2016, gutmann2000a 14f - 190314a, 200112ver.png",
			Zettel{
				"170712a",
				[]string{"Evolution", "Lego bauen", "Perfektion"},
				[]string{"190314a", "200112ver"},
				[]citation{{"nick2016", ""}, {"gutmann2000a", "14f"}},
				[]string{"Gespräch Peter"},
				"170712a - Evolution,Lego bauen, Perfektion - Gespräch Peter, nick2016, gutmann2000a 14f - 190314a, 200112ver.png",
			},
			"",
		},
		// Also parse some failing cases, like without ID, without keywords, etc. and then return an error.
	}
	for _, tc := range tcs {
		got, err := ParseFilename(tc.filename)
		want := tc.zettel
		if diff := cmp.Diff(want, got, cmp.AllowUnexported()); diff != "" {
			t.Errorf(diff)
		}
		var errMsg string
		if err != nil {
			errMsg = err.Error()
		}
		if errMsg != tc.errMsg {
			t.Errorf("Expected `%s`, got `%s`", tc.errMsg, errMsg)
		}
	}
}
