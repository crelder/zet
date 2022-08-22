package parse

import (
	"github.com/crelder/zet"
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestParseFilename(t *testing.T) {
	var tcs = []struct {
		filename string
		zettel   zet.Zettel
		errMsg   string
	}{
		// There is nothing to parse, when empty string is passed as a parameter.
		{"", zet.Zettel{}, "parse Filename: could not parse empty filename string"},

		// This id is not valid.
		{"1923a - Evolution, Lego bauen", zet.Zettel{}, "parse Filename: could not parse id from filename \"1923a - Evolution, Lego bauen\""},
		{"210404Sof - Software.txt", zet.Zettel{}, "parse Filename: could not parse id from filename \"210404Sof - Software.txt\""},

		// The id is not at the start.
		{"Evolution, Lego bauen - 170402a", zet.Zettel{}, "parse Filename: could not parse id from filename \"Evolution, Lego bauen - 170402a\""},

		// There is no id at all in the filename.
		{"Evolution, Lego bauen", zet.Zettel{}, "parse Filename: could not parse id from filename \"Evolution, Lego bauen\""},

		// The minimal structure of a zettel, hence a filename with only an id.
		{"200110d.txt", zet.Zettel{Id: "200110d", Name: "200110d.txt"}, ""},

		// A correct, simple example o fa filename containing only the id and keywords.
		{
			"170712a - Evolution, Lego bauen, Perfektion.txt",
			zet.Zettel{
				Id:       "170712a",
				Keywords: []string{"Evolution", "Lego bauen", "Perfektion"},
				Name:     "170712a - Evolution, Lego bauen, Perfektion.txt",
			},
			"",
		},

		// A correct example with all possibilities to parse: id, keywords, context, literature, and predecessors.
		// Should also work with no space between the comma and the keyword, e.g. 'Evolution,Lego bauen'.
		{
			"170712a - Evolution,Lego bauen, Perfektion - Gespräch Peter, nick2016, gutmann2000a 14f - 190314a.png",
			zet.Zettel{
				Id:          "170712a",
				Keywords:    []string{"Evolution", "Lego bauen", "Perfektion"},
				Predecessor: "190314a",
				References:  []zet.Reference{{Bibkey: "nick2016", Location: ""}, {Bibkey: "gutmann2000a", Location: "14f"}},
				Context:     []string{"Gespräch Peter"},
				Name:        "170712a - Evolution,Lego bauen, Perfektion - Gespräch Peter, nick2016, gutmann2000a 14f - 190314a.png",
			},
			"",
		},

		// A filename with a reference, but without a keyword
		//{
		//	"170712a - nick2016 - 190314a, 200112ver.png",
		//	zet.Zettel{
		//		Id:          "170712a",
		//		Keywords:    []string{},
		//		Predecessor: []string{"190314a", "200112ver"},
		//		References:  []zet.Reference{{Bibkey: "nick2016", Location: ""}},
		//		Context:     []string{},
		//		Name:        "170712a - nick2016 - 190314a, 200112ver.png",
		//	},
		//	"",
		//},
		//{
		// "150302s - 140304t.txt",
		// zet.Zettel{
		//    Id:          "150302s",
		//    Keywords:    nil,
		//    Folgezettel: nil,
		//    Predecessor: []string{"140304t"},
		//    References:   nil,
		//    context:     nil,
		//    Name:        "",
		// },
		// "",
		//},

	}
	for _, tc := range tcs {
		got, err := Filename(tc.filename)
		want := tc.zettel
		if diff := cmp.Diff(want, got); diff != "" {
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
