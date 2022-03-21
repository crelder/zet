package as

import (
	"fmt"
	"github.com/crelder/zet/internal/core/bl"
	"github.com/crelder/zet/internal/core/port"
)

// Viewer contains the application entry point for all operations regarding views upon your zettelkasten.
type Viewer struct {
	ViewPersister port.ViewRepo
	Repo          port.Repo
}

func NewViewer(vp port.ViewRepo, r port.Repo) Viewer {
	return Viewer{
		ViewPersister: vp,
		Repo:          r,
	}
}

// CreateViews creates a folder with different access points (links).
func (v Viewer) CreateViews() error {
	exists := v.ViewPersister.FileExists("VIEWS")
	if exists {
		return fmt.Errorf(
			"folder 'VIEWS' exists in the zettelkasten.\nPlease manually delete this folder in order to execute the command 'views'.")
	}

	zettel, err := v.Repo.GetZettel()
	if err != nil {
		return fmt.Errorf("error creating views: %w", err)
	}
	index, err := v.Repo.GetIndex()
	if err != nil {
		return fmt.Errorf("error creating views: %w", err)
	}

	zk := bl.NewZk(zettel, index)

	// Creating keyword links
	kv := zk.GetKeyLinks()
	err = v.ViewPersister.CreateSyml("keywords", kv)
	if err != nil {
		return fmt.Errorf("error creating keyword views: %w", err)
	}

	// Creating context links
	con := zk.GetConLinks()
	err = v.ViewPersister.CreateSyml("context", con)
	if err != nil {
		return err
	}

	// Creating citation links
	cit := zk.GetCitLinks()
	err = v.ViewPersister.CreateSyml("citations", cit)
	if err != nil {
		return err
	}

	//Creating index links
	for _, i := range zk.Index {
		for _, id := range i.Id {
			links := zk.GetFolgezettel(id)
			err := v.ViewPersister.CreateFolgezettelStruct(i.Topic, links)
			if err != nil {
				return err
			}
		}
	}

	// Creating explore links
	expl := zk.GetSimilar()
	err = v.ViewPersister.CreateSyml("explore", expl)
	if err != nil {
		return err
	}

	return nil
}
