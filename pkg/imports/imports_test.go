package imports

import (
	"errors"
	"fmt"
	"github.com/crelder/zet/pkg/parse"
	fsRepo "github.com/crelder/zet/pkg/transport/fs"
	"io/fs"
	"os"
	"path"
	"testing"
)

func TestCreatePathImport(t *testing.T) {
	// Arrange
	wd, err := os.Getwd()
	if err != nil {
		t.Errorf("could not get the current working dir")
	}
	var pathTestRepo = path.Join(wd, "testdata", "zettelkasten")
	p := parse.New()
	repo := fsRepo.New(pathTestRepo, p)
	importer := New(p, repo, repo)

	// Rebuild a clean state of the zettel folder
	err2 := os.RemoveAll(path.Join(pathTestRepo, "zettel"))
	if err2 != nil {
		t.Errorf("could not remove zettel folder for recreating it")
	}

	// Put one zettel in your zettelkasten.
	zettel := `Post-capitalism
11.8.21

Some thought...`
	fnName := pathTestRepo + "/zettel/" + "211005p - Post-capitalism.txt"
	err = os.MkdirAll(pathTestRepo+"/zettel", 0755)
	if err != nil {
		t.Errorf("could not create zettel folder: %v", err)
	}
	err = os.WriteFile(fnName, []byte(zettel), 0755)
	if err != nil {
		t.Errorf("could not write zettel file: %v", err)
	}

	// Act
	// Load new content that will get persisted in the zettel folder.
	const impSourcePath = "./testdata/new_zettel_files"
	n, err := importer.Import(impSourcePath)

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
		"211005a - P - A, N.txt",
		// 211005p is already in the zettelkasten. Since this zettel has just one keyword, it should try building the
		// id with letters from a to z from the alphabet.
		"211005b - P - B, M.txt",
		// TODO: Create a testcase with Umlauten, see commit b3b45892
	}
	pathZettelTestRepo := pathTestRepo + "/zettel/"
	de1, _ := os.ReadDir(pathZettelTestRepo)

	fmt.Printf("pathZettelTestRepo: %v\n", pathZettelTestRepo)
	for _, entry := range de1 {
		fmt.Println(entry.Name())
	}
	for _, tc := range testcases {
		if _, err := os.Stat(pathZettelTestRepo + tc); err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				t.Errorf("File was not created: %v", tc)
			} else {
				t.Errorf("error occurred: %v", err)
			}
		}
	}
}

func TestCreateFileImport(t *testing.T) {
	// Arrange
	wd, err := os.Getwd()
	if err != nil {
		t.Errorf("could not get the current working dir")
	}
	var pathTestRepo = path.Join(wd, "testdata", "zettelkasten2")
	p := parse.New()
	repo := fsRepo.New(pathTestRepo, p)
	importer := New(p, repo, repo)

	// Rebuild a clean state of the zettel folder
	err2 := os.RemoveAll(path.Join(pathTestRepo, "zettel"))
	if err2 != nil {
		t.Errorf("could not remove zettel folder for recreating it")
	}

	// Put one zettel in your zettelkasten.
	zettel := `Post-capitalism
5.10.21

Some thought...`
	fnName := pathTestRepo + "/zettel/" + "211005p - Post-capitalism.txt"
	err = os.MkdirAll(pathTestRepo+"/zettel", 0755)
	if err != nil {
		t.Errorf("could not create zettel folder: %v", err)
	}
	err = os.WriteFile(fnName, []byte(zettel), 0755)
	if err != nil {
		t.Errorf("could not write zettel file: %v", err)
	}

	// Act
	// Load new content that will get persisted in the zettel folder.
	const uri = "./testdata/new_zettel_file/new_zettel_file.txt"
	n, err := importer.Import(uri)

	// Assert
	if err != nil {
		t.Errorf("error creating import: %v", err)
	}
	if n != 1 {
		t.Errorf("Imported %v files, should have imported 1", n)
	}

	var testcases = []string{
		// Since ID 211005p already exists in the zettelkasten and therefore building the id from the date + the first
		// letter of the first keyword doesn't work (=211005p), use the letter from the second keyword which will result
		// in 211005g.
		"211005a - P - A, N.txt",
	}
	pathZettelTestRepo := pathTestRepo + "/zettel/"
	for _, tc := range testcases {
		if _, err := os.Stat(pathZettelTestRepo + tc); err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				t.Errorf("File was not created: %v", tc)
			} else {
				t.Errorf("error occurred: %v", err)
			}
		}
	}
}
