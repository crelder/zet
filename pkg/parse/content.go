package parse

import (
	"errors"
	"github.com/crelder/zet"
	"strings"
	"time"
)

// Content parses the content of a zettel into a valid filename.
// Each returned filename has a unique id.
func Content(content string, zettel []zet.Zettel) (string, error) {
	z, err := toZettel(content, zettel)
	if err != nil {
		return "", err
	}
	fn, err := toFilename(z)
	if err != nil {
		return "", err
	}
	return fn, nil
}

// toZettel parses the content of a zettel into a zettel instance.
func toZettel(content string, zettel []zet.Zettel) (zet.Zettel, error) {
	var z zet.Zettel
	if content == "" {
		return zet.Zettel{}, errors.New("parse.ToZettel: cannot parse empty content string")
	}

	header := getHeader(content)

	date, err := parseDate(header.date)
	if err != nil {
		return zet.Zettel{}, err
	}

	z.Keywords = parseKeywordsFromHeader(header.keywords)

	con := parseContext(header.contexts)
	z.Folgezettel = con.Folgezettel
	z.References = con.References
	z.Context = con.Context

	id, err2 := generateId(date, z.Keywords, zettel)
	if err2 != nil {
		return zet.Zettel{}, err2
	}
	z.Id = id

	return z, nil
}

// generateId returns a valid, unique id for a zettel (therefore the id does not exist in the zettelkasten yet).
func generateId(t time.Time, keywords []string, zettel []zet.Zettel) (string, error) {
	date := t.Format("060102")

	// Let's try building a unique id via the date and the first letter of one of the keywords
	for i := 0; i < len(keywords); i++ {
		id := date + strings.ToLower(string(keywords[i][0]))
		if idExist(id, zettel) {
			continue
		}
		return id, nil
	}
	// We still don't have a keyword. Let's try building a unique id with letters from the alphabet
	for c := 'a'; c < 'z'; c++ {
		id := date + string(c)
		if !idExist(id, zettel) {
			return id, nil
		}
		continue
	}

	return "", errors.New("generateId: could not build Id from the date and the first letter of one of the keywords")
}

func idExist(id string, zettel []zet.Zettel) bool {
	for _, z := range zettel {
		if z.Id == id {
			return true
		}
	}
	return false
}
