package view

import (
	"fmt"
	"github.com/crelder/zet"
)

// ViewPersister creates links to zettel. These links that serve as access points into your zettelkasten.
// The path where the links lie is specified within the ViewPersister.
//
// PersistIndex creates a tree-like link structure, so-called "Folgezettel".
// This represents the physical representation of how Niklas Luhmann arranged his Zettel in his
// wooden zettelkasten boxes. This is used for creating chains of thoughts.

// PersistInfo persists some information like a list of keywords used in your zettelkasten and the number of occurrences.
type ViewPersister interface {
	PersistIndex(links map[string]string) error // links[linkName]targetID
	PersistInfo(m map[string][]string) error
}

// Viewer contains the application entry point for all operations regarding views upon your zettelkasten.
// Viewer satisfies the zet.Viewer interface.
type Viewer struct {
	ViewPersister ViewPersister
	Repo          zet.Repo
}

func New(ip ViewPersister, r zet.Repo) Viewer {
	return Viewer{
		ViewPersister: ip,
		Repo:          r,
	}
}

// CreateViews creates a folder with different access points (links).
func (v Viewer) CreateViews() error {
	zettel, err := v.Repo.GetZettel()
	if err != nil {
		return fmt.Errorf("error creating views: %w", err)
	}
	index, err := v.Repo.GetIndex()
	if err != nil {
		return fmt.Errorf("error creating index: %w", err)
	}

	references, err := v.Repo.GetBibkeys()
	if err != nil {
		return err
	}

	// Create a method, which returns all paths like "Komplexität/180215a - Komplexität, ..../180215a - Komplexität, ..."
	folgezettelMap, err := getFolgezettelMap(zettel, index)
	if err != nil {
		return err
	}
	// PersistIndex all these paths via a call v.ViewPersister.PersistIndex(map[paths][]ids). It creates already everything in "zettelkasten/INDEX/"
	// Concrete Implementierung heißt FsIndexPersister.
	err = v.ViewPersister.PersistIndex(folgezettelMap)
	if err != nil {
		return err
	}

	// Call method that persists all these info v.InfoPersister.PersistIndex(name, []string).
	// Concrete Implementierung heißt CSVPersister.
	infos := getInfos(zettel, index, references)

	err = v.ViewPersister.PersistInfo(infos)
	if err != nil {
		return err
	}

	return nil
}

func getZettel(id string, zettel []zet.Zettel) (zet.Zettel, error) {
	// TODO: Make map out of it?
	for _, z := range zettel {
		if z.Id == id {
			return z, nil
		}
	}
	return zet.Zettel{}, fmt.Errorf("view: zettel with id %v not found", id)
}
