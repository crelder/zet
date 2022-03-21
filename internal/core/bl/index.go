package bl

import "strings"

// Index represents entry points into line of thoughts within your zettelkasten.
//
// The power of the zettelkasten which Niklas Luhmann used came from working with such an index in combination with folgezettel
// which result in lines of thought.
//
// When starting to you use a zettelkasten, the zettelkasten does not contain many zettel. Therefore working with the
// index at the beginning might not that useful. In the beginning it is more helpful to just use keywords as entry points
// to own thoughts (= zettel) within the zettelkasten.
type Index struct {
	// Should both variables be with a capital letter?
	Topic string
	Id    []string
}

// ParseIndex parses the content of a textfile (the index) and
// returns an Index struct.
func ParseIndex(c string) []Index {
	if c == "" {
		return nil
	}
	var result []Index
	lines := strings.Split(c, "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		indexEntries := strings.Split(line, ":")
		topic := strings.TrimSpace(indexEntries[0])
		ids := strings.Split(indexEntries[1], ",")
		var trimmedIds []string
		for _, id := range ids {
			trimmedIds = append(trimmedIds, strings.TrimSpace(id))
		}
		result = append(result,
			Index{Topic: topic, Id: trimmedIds})
	}
	return result
}
