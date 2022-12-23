package zet

// Zettel holds the metadata of one thought.
type Zettel struct {
	Id          string
	Keywords    []string
	Folgezettel []string
	Predecessor string
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

// InconErr stands for inconsistency error, indicating that something is not right with your zettelkasten.
// They are different from errors, since the programs just is aware of them but can continue functioning.
// If you want to be sure, that zet operates correctly on your zettelkasten, make sure that
// you don't have any inconsistencies in your zettelkasten. Run `zet validate` to get a list of inconsistencies.
type InconErr struct {
	Message error
}

func (i InconErr) Error() string {
	return i.Error()
}
