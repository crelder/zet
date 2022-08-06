package imports

import (
	"fmt"
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
	importer := New(p, repo)

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
	contents, err5 := getContentsFromPath(impSourcePath)
	if err5 != nil {
		t.Errorf("error getting file content: %v", err5)
	}

	// Act
	n, err := importer.Import(contents)

	// Assert
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

func getContentsFromPath(path string) ([]string, error) {
	var contents []string
	filepaths, err := os.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("cli: error %q reading path: %v", err, path)
	}
	filepaths = filterAllowed(filepaths)

	for _, fp := range filepaths {
		dat, _ := os.ReadFile(path + "/" + fp.Name())
		contents = append(contents, string(dat))
	}
	return contents, nil
}

func filterAllowed(filepaths []os.DirEntry) []os.DirEntry {
	var fps []os.DirEntry
	for _, fp := range filepaths {
		if isAllowed(fp.Name()) {
			fps = append(fps, fp)
		}
	}
	return fps
}

func isAllowed(fn string) bool {
	if fn[len(fn)-3:] == "txt" {
		return true
	}
	return false
}
