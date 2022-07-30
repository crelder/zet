package imports

import (
	"fmt"
	"github.com/crelder/zet/pkg/parse"
	"github.com/crelder/zet/pkg/transport/repo"
	"os"
	"testing"
)

func TestCreateImports(t *testing.T) {
	// Arrange
	wd, errors := os.Getwd()
	if errors != nil {
		t.Errorf("could not get the current working dir")
	}
	var pathTestRepo = wd + "/testdata/zettelkasten/"
	parser := parse.New()
	r := repo.New(pathTestRepo, parser)
	importer := New(parser, r, r)

	testFiles := []string{
		pathTestRepo + "zettel/211005g - Prozess, Glück, Enttäuschung, Hoffnung - Paul Watzlawick, Anleitung zum Unglücklich sein.txt",
		pathTestRepo + "zettel/211005a - Probleme - Dave Cheney.txt"}
	err := removeFiles(testFiles) // These will be created again during the test
	if err != nil {
		t.Errorf("error removing test files")
	}

	const impSourcePath = "./testdata/new_zettel_files"
	contents, err := getContentsFromPath(impSourcePath)
	if err != nil {
		t.Errorf("error getting file content: %v", err)
	}

	// Act
	n, err := importer.Import(contents)

	// Assert
	if err != nil {
		t.Errorf("error creating imports. path = %q", err)
	}
	if n != 2 {
		t.Errorf("Imported %v files, should have imported 2", n)
	}

	var testcases = []string{
		// Since ID 211005p already exists in the zettelkasten and therefore building the id from the date + the first
		// letter of the first keyword doesn't work (=211005p), use the letter from the second keyword which will result
		// in 211005g.
		"zettel/211005g - Prozess, Glück, Enttäuschung, Hoffnung - Paul Watzlawick, Anleitung zum Unglücklich sein.txt",
		// 211005p is already in the zettelkasten. Since this zettel has just one keyword, it should try building the
		// id with letters from a to z from the alphabet.
		"zettel/211005a - Probleme - Dave Cheney.txt",
	}
	for _, tc := range testcases {
		if _, err := os.Stat(pathTestRepo + tc); err != nil {
			t.Errorf("File was not created: %v", tc)
		}
	}
}
func removeFiles(testFiles []string) error {
	for _, tf := range testFiles {
		err := os.RemoveAll(tf)
		if err != nil {
			return err
		}
	}
	return nil
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
