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

// CreateImports reads all the zettel contents from the parameter path.
// For every zettel content, a valid filename is generated containing all the zettel's metadata and
// a unique id.
//
// In case of success CreateImports returns the number of zettel created and a nil error.
// In case of an error CreateImports returns 0 (no zettel are created) and the error.
func (i Importer) CreateImports(path string) (int, error) {
	contents, err := i.reader.GetContents(path)
	if err != nil {
		return 0, err
	}

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

		z, err4 := i.parser.Filename(filename)
		if err4 != nil {
			return 0, err4
		}

		// Make sure that a following import is not using the same id as this zettel.
		zettel = append(zettel, z)
		zettelFiles[filename] = content
		continue
	}

	n, err5 := i.repo.Save(zettelFiles)
	if err5 != nil {
		return n, err5
	}

	return n, nil
}
