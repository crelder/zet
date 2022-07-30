// Package zet holds the central elements of this application.
//
// It contains three models (see model.go) that make the key concept of this application:
// * Zettel
// * Reference
// * Index
//
// It also holds all interfaces that are
//   a) driving ports of the application therefore making clear what this application does and
//   b) used in more than one package, so that from all other packages the dependencies point only to this package.
package zet

// Importer persists zettel content.
//
// Import takes one or more zettel contents and persists each content.
// In case of an error it returns the number of zettel contents already persisted until the occurrence of the error.
type Importer interface {
	Import(zettelContents []string) (int, error)
}

// Initiator supports starting with this personal knowledge management system.
//
// Init will create an empty zettelkasten.
//
// InitExample will download an example zettelkasten.
type Initiator interface {
	Init() error
	InitExample() error
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
// Val returns all inconsistencies of your zettelkasten as the first parameter.
// Inconsistencies can be:
//     * dead links
//     * double ids
//     * missing reference entry
// The second return parameter contains a potential error.
type Validator interface {
	Val() ([]error, error)
}

// Repo gives access to the content of your zettelkasten.
//
// GetZettel returns Zettel entities and all errors that occurred while fetching the zettel,
// e.g. parsing errors.
// The structure is always the same since the raw filenames get sorted by name before processing so that
// the order of links, etc. is always the same; therefore, the VIEWS structure also stays the same.
//
// GetIndex returns a Index that maps thematic entry topics to one or more ids.
//
// GetBibkeys returns a list of bibkeys representing literature references.
//
// Save takes a map[filename]content of zettel and saves these.
// filename is the name of the file that holds the thought. Content is the text content of your thought.
// In case of success it returns a nil error and the number of zettel persisted.
// In case of a failure, it returns the error and the number of zettel it has written until the error occurred.
type Repo interface {
	GetZettel() ([]Zettel, []error, error)
	GetIndex() (Index, []error, error)
	GetBibkeys() ([]string, error)
	Save(content map[string]string) (int, error)
}

// Parser handles all functionality regarding parsing from and
// sometimes to raw data like filenames, literature entries and index entries to zettel.
type Parser interface {
	Content(string, []Zettel) (string, error)
	Filename(string) (Zettel, error)
	Index(content string) (Index, []error)
	Reference(d string) []string
}
