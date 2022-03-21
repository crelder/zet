package port

import (
	"github.com/crelder/zet/internal/core/bl"
)

// Repo is the interface that wraps the repository with all the zettelkasten content.
//
// GetZettel returns Zettel entities and all errors that occurred while fetching the zettel,
// e.g. parsing errors.
//
// GetIndex returns a map with [index term][]ids.
type Repo interface {
	GetZettel() ([]bl.Zettel, error)
	GetIndex() ([]bl.Index, error)
}

// ViewRepo creates links that serve as access points into your zettelkasten.
// The path where the links lie is specified within the ViewRepo.
//
// CreateSyml creates list of zettel links
// prefix is the type of the view (e.g. keywords, literature sources, context).
// m contains as keys the type, e.g. keywords, and as []string one or more ids.
//
// CreateFolgezettelStruct creates a tree like link structure, so called "Folgezettel".
// This represents the physical representation how Niklas Luhmann arranged his Zettel in his
// wooden zettelkasten boxes. This is used for creating chains of thoughts.
//
// FileExists checks within your ViewRepo if a file or link was set.
type ViewRepo interface {
	CreateSyml(prefix string, m map[string][]string) error
	CreateFolgezettelStruct(prefix string, links []bl.Symlink) error
	FileExists(link string) bool
}

// ImportRepo persists imported textfiles with a valid filename.
// The target path for persisting is set within the repo as a property.
//
// GetContents takes a path to a folder with textfiles and returns the contents of these files.
//
// SaveImports takes a map[filename]content and persists these.
// In case of success it returns a nil error and the number of files persisted.
// In case of a failure, it returns the error and the number of files it has written until the error ocurred.
type ImportRepo interface {
	GetContents(path string) []string
	SaveImports(f map[string]string) (int, error)
}
