package as

import (
	"github.com/crelder/zet/internal/core/bl"
	"github.com/crelder/zet/internal/core/port"
)

// Importer provides functionality around importing new text zettel.
type Importer struct {
	persister port.ImportRepo
	repo      port.Repo
}

func NewImporter(persister port.ImportRepo, repo port.Repo) Importer {
	return Importer{
		persister: persister,
		repo:      repo}
}

var zk bl.Zettelkasten

// CreateImports takes as parameter a path and reads all the file contents from this path.
// For every content a proper filename is generated is persisted in a separate import folder
// within the zettelkasten.
// The filename is automatically generated and
// is a valid filename with a valid id containing all the zettel's metadata.
//
// In case of success CreateImports returns the number of files created and a nil error.
// In case of an error, 0 is returned and an error (no files are created / persisted).
func (i Importer) CreateImports(path string) (int, []error) {
	contents := i.persister.GetContents(path)

	zettel, err := i.repo.GetZettel()
	if err != nil {
		return 0, nil // TODO: Return errors
	}
	zk = bl.NewZk(zettel, nil)

	var e []error
	zettelFiles := make(map[string]string)

	for _, content := range contents {
		filename, err := zk.ParseFileContent(content)
		if err != nil {
			e = append(e, err)
			continue
		}
		zettelFiles[filename] = content
	}

	// Persisting new zettel should only take place, when there are no errors at all.
	if e != nil {
		return 0, e
	}
	n, err := i.persister.SaveImports(zettelFiles)
	if err != nil {
		s := make([]error, 1)
		s[0] = err
		return 0, s
	}

	return n, nil
}
