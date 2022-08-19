package imports

import (
	"github.com/crelder/zet/pkg/parse"
	"github.com/crelder/zet/pkg/transport/fs"
	"os"
	"testing"
)

func TestCreateImports(t *testing.T) {
	// Arrange
	wd, err := os.Getwd()
	if err != nil {
		t.Errorf("could not get the current working dir")
	}
	var pathTestRepo = wd + "/testdata/zettelkasten"
	p := parse.New()
	repo := fs.New(pathTestRepo, p)
	importer := New(p, repo, repo)

	// Rebuild a clean state of the zettel folder
	err2 := os.RemoveAll(pathTestRepo + "/zettel")
	if err2 != nil {
		t.Errorf("could not remove zettel folder for recreating it")
	}

	// Put one zettel in your zettelkasten.
	zettel := `Post-capitalism
11.8.21

Some thought...`
	err3 := os.MkdirAll(pathTestRepo+"/zettel", 0755)
	if err3 != nil {
		t.Errorf("could not create zettel folder: %v", err3)
	}
	fnName := pathTestRepo + "/zettel/" + "211005p - Post-capitalism.txt"
	err4 := os.WriteFile(fnName, []byte(zettel), 0755)
	if err4 != nil {
		t.Errorf("could not write zettel file: %v", err4)
	}

	// Load new content that will get persisted in the zettel folder.
	const impSourcePath = "./testdata/new_zettel_files"

	// Assert
	n, err := importer.Import(impSourcePath)
	if err != nil {
		t.Errorf("error creating import: %v", err)
	}
	if n != 2 {
		t.Errorf("Imported %v files, should have imported 2", n)
	}

	var testcases = []string{
		// Since ID 211005p already exists in the zettelkasten and therefore building the id from the date + the first
		// letter of the first keyword doesn't work (=211005p), use the letter from the second keyword which will result
		// in 211005g.
		"211005g - Prozess, Glück, Enttäuschung, Hoffnung - Paul Watzlawick, Anleitung zum Unglücklich sein.txt",
		// 211005p is already in the zettelkasten. Since this zettel has just one keyword, it should try building the
		// id with letters from a to z from the alphabet.
		"211005a - Probleme - Dave Cheney.txt",
	}
	pathZettelTestRepo := pathTestRepo + "/zettel/"
	for _, tc := range testcases {
		if _, err := os.Stat(pathZettelTestRepo + tc); err != nil {
			t.Errorf("File was not created: %v", tc)
		}
	}
}
