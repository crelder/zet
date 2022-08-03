package view

import (
	"fmt"
	"github.com/crelder/zet"
)

// IndexPersister creates symlinks to zettel. These symlinks that serve as access points into your zettelkasten.
// The path where the links lie is specified within the IndexPersister.
//
// Persist creates a tree-like symlink structure, so-called "Folgezettel".
// This represents the physical representation of how Niklas Luhmann arranged his Zettel in his
// wooden zettelkasten boxes. This is used for creating chains of thoughts.

// CreateInfo persists some information like a list of keywords used in your zettelkasten and the number of occurrences.
type IndexPersister interface {
	Persist(links map[string]string) error // links[linkName]targetID
}

type InfoPersister interface {
	CreateInfo(prefix string, m map[string][]string) error
}

// Viewer contains the application entry point for all operations regarding views upon your zettelkasten.
// Viewer satisfies the zet.Viewer interface.
type Viewer struct {
	IndexPersister IndexPersister
	InfoPersister  InfoPersister
	Repo           zet.Repo
}

func New(ip IndexPersister, infoP InfoPersister, r zet.Repo) Viewer {
	return Viewer{
		IndexPersister: ip,
		InfoPersister:  infoP,
		Repo:           r,
	}
}

// CreateViews creates a folder with different access points (links).
func (v Viewer) CreateViews() error {
	zettel, err := v.Repo.GetZettel()
	if err != nil {
		return fmt.Errorf("error creating views: %w", err)
	}
	index, err2 := v.Repo.GetIndex()
	if err2 != nil {
		return fmt.Errorf("error creating index: %w", err2)
	}

	// INDEX
	// Create a method, which returns all paths like "Komplexität/180215a - Komplexität, ..../180215a - Komplexität, ..."
	folgezettelMap, err := getFolgezettelMap(zettel, index)
	if err != nil {
		return err
	}

	// Persist all these paths via a call v.IndexPersister.Persist(map[paths][]ids). It creates already everything in "zettelkasten/INDEX/"
	// Concrete Implementierung heißt FsIndexPersister.
	err = v.IndexPersister.Persist(folgezettelMap)
	if err != nil {
		return err
	}

	// INFO
	// Method where you get the info: ids
	// Method where you get the info: keywords
	// Method where you get the info: context
	// Method where you get the info: references
	// Method where you get the info: unlinked
	// Method where you get the info: unindexed
	// always in the format: [id, Häufigkeit]
	//
	// Call method that persists all these info v.InfoPersister.Persist(name, []string).
	// Concrete Implementierung heißt CSVPersister.

	//err3 := v.createInfos(zettel)
	//if err3 != nil {
	//	return err3
	//}

	return nil
}

func getZettel(id string, zettel []zet.Zettel) (zet.Zettel, error) {
	for _, z := range zettel {
		if z.Id == id {
			return z, nil
		}
	}
	return zet.Zettel{}, fmt.Errorf("zettel with id %v not found", id)
}
