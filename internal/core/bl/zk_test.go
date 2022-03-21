package bl

import (
	"github.com/google/go-cmp/cmp"
	"sort"
	"testing"
)

func TestKeywords(t *testing.T) {
	// Arrange
	zk := NewZk(
		[]Zettel{
			{
				Id:       "170312a",
				Keywords: []string{"Water", "Fire"},
			},

			{
				Id:       "181101f",
				Keywords: []string{"Fire", "Air"},
			}},
		nil,
	)

	tc := map[string][]string{
		"Water": []string{"170312a"},
		"Fire":  []string{"170312a", "181101f"},
		"Air":   []string{"181101f"},
	}

	// Act
	got := zk.GetKeyLinks()

	// Assert
	if diff := cmp.Diff(tc, got); diff != "" {
		t.Errorf(diff)
	}
}

// TestGetFolgezettel tests if the folgezettel structure (= tree structure) structure in the zettelkasten is correctly translated into a directory hierachy.
func TestGetFolgezettel(t *testing.T) {
	// Arange
	var filenames = []string{
		"170224a - Polymorphismus, Objetkorientierte Programmierung, Schnittstelle - Eloquente Javascript.png",
		"180522a - Komplexität, Komplexer Mensch, Menschliche Entwicklung - Thermodynamische Tiefe - csikszentmihalyi2010 328, Pagels - 210520var.png",
		"190119d.txt",
		"190119e - Komplexität - 180522a, 190119d.txt",
		"210328obj - Objektorientiert, Programmierung, Vererbung, Kapselung, Komposition - kernighan2016 155 - 170224a.pdf",
		"210520var - Varietät, Komplexität - ropohl2021 71.txt",
		"220115p - Refactoring, Programmieren, Aufgabe, Testen - clausen2021 5 - 220116s.pdf",
		"220116s - Spezifikation, Signatur, Null.pdf",

		// These zettel filenames are used to show the deep structure when calling "170712d"
		"170719f - Lernen, Vereinfachung.png",
		"170313b - Komplexität, feste Hierarchie - 210111a.txt",
		"210111a - Biologie, Physik - 210919e.txt", // Zettel 210919e doesn't exist on purpose for the test.
		"150211j - Gedanken - 150211j.txt",         // Zettel refers to itself on purpose.
		"170712d - Disruptive Innovation, Makroevolution, Mikroevolution - Vortrag Günter Theißen - 170712g, 170712e, 170904b.png",
		"170712e - EcoDevo, Gene, Entwicklung - Vortrag Günter Theißen - Selbstveränderung als Schlüßel zur Orga Entwicklung - 170815f, 170918a.png",
		"170918a - Komplexitätsforschung, einfache Strukturen, komplexe Strukturen, Problem - ladyman2013 60.png",
		"170815f - Neuere Systemtheorie, Gründe für wachsende Komplexität, Steuerungsmechanismen, hohe Komplexität verarbeiten - willke2000 14.png",
		"170712g - Häufigkeit, Wichtigkeit, Dinosaurier - Vortrag Günter Theißen - 171213a.png",
		"171213a - Kulturelle Evolution, Natürliche Evolution, Sprache, Wissen - Vortrag Christoph Kleinedam Biologie.png",
		"170904b - Evolutionstheorie, Selektion, Reproduktion, Variation, gerichtete Evolution, Stabilität der Selektionskräfte - schurz2011.png",
	}
	var zettels []Zettel
	for _, filename := range filenames {
		z, _ := ParseFilename(filename)
		zettels = append(zettels, z)
	}
	zk := NewZk(zettels, nil)

	// Assert
	var testcases = []struct {
		id       string
		symlinks []Symlink
	}{
		// If no id is provided no symlinks are created.
		{"", nil},
		// A wrong or non-existing id returns no symlinks.
		{"000a", nil},
		{"999999z", nil},

		{"190119e",
			[]Symlink{{"190119e", "00_190119e - Komplexität - 180522a, 190119d.txt"},
				{"190119d", "01_190119d/00_190119d.txt"},
				{"180522a", "02_180522a - Komplexität, Komplexer Mensch, Menschliche Entwicklung - Thermodynamische Tiefe - csikszentmihalyi2010 328, Pagels - 210520var.png"},
				{"210520var", "03_210520var - Varietät, Komplexität - ropohl2021 71.txt"}}},

		{"220115p",
			[]Symlink{{"220115p", "00_220115p - Refactoring, Programmieren, Aufgabe, Testen - clausen2021 5 - 220116s.pdf"},
				{"220116s", "01_220116s - Spezifikation, Signatur, Null.pdf"}},
		},

		{"210328obj",
			[]Symlink{{"210328obj", "00_210328obj - Objektorientiert, Programmierung, Vererbung, Kapselung, Komposition - kernighan2016 155 - 170224a.pdf"},
				{"170224a", "01_170224a - Polymorphismus, Objetkorientierte Programmierung, Schnittstelle - Eloquente Javascript.png"}},
		},
		{"170712d", []Symlink{
			{"170712d", "00_170712d - Disruptive Innovation, Makroevolution, Mikroevolution - Vortrag Günter Theißen - 170712g, 170712e, 170904b.png"},
			{"170712e", "01_170712e/00_170712e - EcoDevo, Gene, Entwicklung - Vortrag Günter Theißen - Selbstveränderung als Schlüßel zur Orga Entwicklung - 170815f, 170918a.png"},
			{"170918a", "01_170712e/01_170918a/00_170918a - Komplexitätsforschung, einfache Strukturen, komplexe Strukturen, Problem - ladyman2013 60.png"},
			{"170815f", "01_170712e/02_170815f - Neuere Systemtheorie, Gründe für wachsende Komplexität, Steuerungsmechanismen, hohe Komplexität verarbeiten - willke2000 14.png"},
			{"170904b", "02_170904b/00_170904b - Evolutionstheorie, Selektion, Reproduktion, Variation, gerichtete Evolution, Stabilität der Selektionskräfte - schurz2011.png"},
			{"170712g", "03_170712g - Häufigkeit, Wichtigkeit, Dinosaurier - Vortrag Günter Theißen - 171213a.png"},
			{"171213a", "04_171213a - Kulturelle Evolution, Natürliche Evolution, Sprache, Wissen - Vortrag Christoph Kleinedam Biologie.png"},
		}},
	}
	for _, tc := range testcases {
		t.Run(tc.id, func(t *testing.T) {
			got := zk.GetFolgezettel(tc.id)
			if diff := cmp.Diff(got, tc.symlinks); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}

func TestValidate(t *testing.T) {
	// Arrange: Set up a zettelkasten,
	// that should trigger all possible errors defined in ValidatorErr.ErrType.
	fn := []string{
		// Non-unique id
		"180112a - Water, Fire, Ice.txt",
		"180112a - Water, Fire, Ice.txt",

		// Link target id '190311d' doesn't exist
		"170327f - Water, Fire - 190311d.txt",
	}
	var zettel []Zettel
	for _, f := range fn {
		z, _ := ParseFilename(f)
		zettel = append(zettel, z)
	}
	zk := NewZk(zettel, nil)

	// Act
	zkerr := zk.Validate()

	// Assert
	errs := []ValidatorErr{
		{NotUniqueId, "180112a"},
		{TargetLinkIdNotThere, "190311d"},
	}
	zkerr = sortErrs(zkerr)
	errs = sortErrs(errs)

	if diff := cmp.Diff(zkerr, errs); diff != "" {
		t.Errorf(diff)
	}
}

func sortErrs(errs []ValidatorErr) []ValidatorErr {
	sort.Slice(errs, func(i, j int) bool {
		return errs[i].ErrType < errs[j].ErrType
	})
	return errs
}

func TestParseFileContent(t *testing.T) {
	var tcs = []struct {
		in     string // The file content
		out    string // The file name
		errMsg string // Returned error message
	}{
		// Only the header (first two to three lines, separated by a blank line from the zettel content)
		// is used for creating the filename (with the meta data).
		{"Modelle, Theorien\n12.1.2020\nropohl2013a 14,121117a\n\nHere the zettel content starts...",
			"200112m - Modelle, Theorien - ropohl2013a 14 - 121117a.txt", ""},

		// This header only consists of keywords (first line) and a date (second line)
		// The third line with context information (bibkeys, context) is optional - and here not present
		// The id '200112m' already exists (see above), therefore the id is build with the first letter of the next keyword.
		{"Minen, Bergbau\n12.1.2020\n\nHere the zettel content starts...",
			"200112b - Minen, Bergbau.txt", ""},

		// When the id can't get generated via the date and the first letter of one of the keywords
		// try using all the letters from the alphabet starting with 'a'.
		{"Mars, Building\n12.1.2020\n\nSome content...",
			"200112a - Mars, Building.txt", ""},

		// Context that only consists of a bibkey, should also be parsed correctly into a filename
		{"Lesen, Index, Zettelkasten\n12.8.21\nadler1972\n\nAdler Wie man ein Buch liest",
			"210812l - Lesen, Index, Zettelkasten - adler1972.txt", ""},

		// This shows the full parsing potential.
		// It should also parse the context and folgezettel, if no ' - ' is provided to separate.
		// It should also handle additional whitespaces between the contexts
		{"Risiko, Unsicherheit\n170312\n Paul Ehrlich,  181201f, kahn1985 12,  Movie Dunkirk, 200812c, greyer1987",
			"170312r - Risiko, Unsicherheit - Paul Ehrlich, Movie Dunkirk, kahn1985 12, greyer1987 - 181201f, 200812c.txt", ""},

		// Check the returned errors
		// If a space is missing between keywords, it should be added in the filename
		{"Steuerung,Regelung, SOLL\n1.3.21\nJürgen Kleppert, Der Dreher\n\nHier fehlt das Leerzeichen",
			"210301s - Steuerung, Regelung, SOLL - Jürgen Kleppert, Der Dreher.txt", ""},

		// Different date formats should be parsed
		{"Date, Format\n14.6.21", "210614d - Date, Format.txt", ""},
		{"Date, Format\n14.07.21", "210714d - Date, Format.txt", ""},
		{"Format\n4.08.21", "210804f - Format.txt", ""},
		{"Picture\n4.02.2020", "200204p - Picture.txt", ""},

		// Returning an error, when the data can't get parsed
		{"Date, Format\nNot a date", "", "cannot parse the date"},

		// When there is no content provided, there is no header that can be parsed into a zettel,
		// which then can be parsed into a filename.
		{"", "", "could not parse neither the date, nor the keywords from the header"},

		// A date is missing, which is used for defining the id.
		{"Risiko, Unsicherheit\nPaul Ehrlich, 200812c, @181201f, @kahn1985 12", "", "cannot parse the date"},
	}

	var zk = NewZk(nil, nil)
	for _, tc := range tcs {
		got, err := zk.ParseFileContent(tc.in)
		if got != tc.out {
			t.Errorf("Wanted: %q, got: %q", tc.out, got)
		}
		var errMsg string
		if err != nil {
			errMsg = err.Error()
		}
		if errMsg != tc.errMsg {
			t.Errorf("Expeted error message %q, got %q", tc.errMsg, errMsg)
		}
	}
}
