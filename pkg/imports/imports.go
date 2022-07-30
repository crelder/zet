package imports

import (
	"github.com/crelder/zet"
)

// Importer provides functionality for importing new text zettel.
// Importer satisfies the zet.Importer interface.
type Importer struct {
	parser zet.Parser
	reader Reader
	repo   zet.Repo
}

func New(p zet.Parser, r Reader, repo zet.Repo) Importer {
	return Importer{
		parser: p,
		reader: r,
		repo:   repo}
}

// Reader gets zettel as text.
//
// GetContents takes a path to a folder with textfiles and returns their contents.
type Reader interface {
	GetContents(path string) ([]string, error)
}

// Import creates for every slice entry, a zettel content,
// a valid filename with all the zettel's metadata and a unique id.
//
// In case of success Import returns the number of zettel created and a nil error.
// In case of an error Import returns the error and the number of zettel
// that got persisted until the error occurrence.
func (i Importer) Import(contents []string) (int, error) {
	zettel, _, err2 := i.repo.GetZettel()
	if err2 != nil {
		return 0, err2
	}

	zettelFiles := make(map[string]string)
	for _, content := range contents {
		filename, err3 := i.parser.Content(content, zettel)
		if err3 != nil {
			return 0, err3
		}
		zettelFiles[filename] = content

		// Make sure that a following import is not using the same id as this zettel.
		z, err4 := i.parser.Filename(filename)
		if err4 != nil {
			return 0, err4
		}
		zettel = append(zettel, z)

		continue
	}

	n, err5 := i.repo.Save(zettelFiles)
	if err5 != nil {
		return n, err5
	}

	return n, nil
}
