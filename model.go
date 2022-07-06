package zet

// Zettel holds one thought.
type Zettel struct {
	Id          string
	Keywords    []string
	Folgezettel []string
	Predecessor []string
	References  []Reference
	Context     []string
	Name        string // the filename, e.g. '170212g - Go.txt'
}

// Index represents thematic entry points into a line of thoughts within your zettelkasten.
//
// The power of the zettelkasten, which Niklas Luhmann used,
// came from working with such an index in combination with folgezettel,
// which results in lines of thought.
//
// At the beginning, your zettelkasten does not contain many zettel. Therefore, working with the
// index at the beginning might not be that useful. In the beginning,
// it is more helpful just to use keywords as entry points
// to own thoughts (= zettel) within the zettelkasten.
//
// Map assigns a topic to one or more ids (map[topic][]ids).
type Index map[string][]string

// Reference has a bibkey which refers to a literature reference (e.g. book, paper, etc.). E.g. "welter2011".
// With a literature reference and the location (e.g. page number, chapter, etc.)
// you can precisely define the source of your thought.
type Reference struct {
	Bibkey   string
	Location string
}
