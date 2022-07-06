package parse

import (
	"github.com/crelder/zet"
	"testing"
)

func TestContent(t *testing.T) {
	var tcs = []struct {
		in     string // zettel content
		out    string // zettel filename
		errMsg string // Returned error message
	}{
		// Only the header (first two to three lines, separated by a blank line from the zettel content)
		// is used for creating the filename.
		{"Modelle, Theorien\n12.1.2020\nropohl2013a 14,121117a\n\nHere the zettel content starts...",
			"200112m - Modelle, Theorien - ropohl2013a 14 - 121117a.txt", ""},

		// This shows the full parsing potential.
		// It should also parse the context and folgezettel, if no ' - ' is provided to separate.
		// It should also handle additional whitespaces between the contexts.
		{"Risiko, Unsicherheit\n170312\n Paul Ehrlich,  181201f, kahn1985 12,  Movie Dunkirk, 200812c, greyer1987",
			"170312r - Risiko, Unsicherheit - Paul Ehrlich, Movie Dunkirk, kahn1985 12, greyer1987 - 181201f, 200812c.txt", ""},

		// This header only consists of keywords (first line) and a date (second line).
		// The third line with context information (bibkeys, context) is optional - and here not present.
		{"Bergbau, Minen\n12.1.2020\n\nHere the zettel content starts...",
			"200112b - Bergbau, Minen.txt", ""},

		// context that only consists of a bibkey, should also be parsed correctly into a filename
		{"Lesen, Index, Zettelkasten\n12.8.21\nadler1972\n\nAdler Wie man ein Buch liest",
			"210812l - Lesen, Index, Zettelkasten - adler1972.txt", ""},

		// If a space is missing between keywords, the space should be added in the filename
		{"Steuerung,Regelung, SOLL\n1.3.21\nJan Kleppert, Der Dreher\n\nHier fehlt das Leerzeichen",
			"210301s - Steuerung, Regelung, SOLL - Jan Kleppert, Der Dreher.txt", ""},

		// Different date formats should be parsed
		{"Date, Format\n14.6.21", "210614d - Date, Format.txt", ""},
		{"Date, Format\n14.07.21", "210714d - Date, Format.txt", ""},
		{"Date, Format\n4.08.21", "210804d - Date, Format.txt", ""},
		{"Date, Format\n4.02.2020", "200204d - Date, Format.txt", ""},

		// Returning an error, when the data can't get parsed
		{"Date, Format\nNot a date", "", "parseDate: could not parse date"},

		// When there is no content provided, there is no header that can be parsed into a zettel,
		// which then can be parsed into a filename.
		{"", "", "parse.ToZettel: cannot parse empty content string"},

		// A date is missing, which is used for defining the id.
		{"Risiko, Unsicherheit\nPaul Ehrlich, 200812c, 181201f, kahn1985 12", "", "parseDate: could not parse date"},
	}

	for _, tc := range tcs {
		got, err := Content(tc.in, []zet.Zettel{})
		if got != tc.out {
			t.Errorf("Got: %q, wanted: %q", got, tc.out)
		}
		var errMsg string
		if err != nil {
			errMsg = err.Error()
		}
		if errMsg != tc.errMsg {
			t.Errorf("Got %q, wanted %q", errMsg, tc.errMsg)
		}
	}
}
