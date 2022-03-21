package port

import "github.com/crelder/zet/internal/core/bl"

// Importer is the instance for accessing all functionality regarding
// importing new text files with thoughts (zettel).
//
// CreateImports will persist the text files that lay in path.
type Importer interface {
	CreateImports(path string) (int, []error)
}

// Viewer is the instance for accessing all functionality regarding
// entry points into your zettelkasten (so-called "views").
//
// CreateViews will create all access points into your zettelkasten.
type Viewer interface {
	CreateViews() error
}

// Validator is the instance for accessing all functionality regarding
// the consistency and health checks for your zettelkasten.
//
// Val returns all consistencies of your zettelkasten.
type Validator interface {
	Val() []bl.ValidatorErr
}
