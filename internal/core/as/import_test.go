package as

import (
	"github.com/crelder/zet/internal/adapter/driven"
	"os"
	"testing"
)

func TestCreateImports(t *testing.T) {
	// Arrange
	wd, errors := os.Getwd()
	if errors != nil {
		t.Errorf("Could not get the current working dir")
	}
	var pathTestRepo = wd + "/testdata/import/zettelkasten/"
	r := repo.NewRepo(pathTestRepo)
	importer := NewImporter(r, r)

	impPath := pathTestRepo + "IMPORT"
	clearPath(impPath) // Remove this directory, which might got created in a previous test

	// Act
	const impSourcePath = "./testdata/import/new_zettel_files"

	// Assert
	n, errs := importer.CreateImports(impSourcePath)
	if errs != nil {
		t.Errorf("Error creating imports. path = %q", impSourcePath)
	}
	if n != 2 {
		t.Errorf("Should have imported 2 files, imported %v", n)
	}

	var tcs = []string{
		// Since ID 211005p already exists in the zettelkasten, ID 211005g should be assigned to the new zettel.
		"IMPORT/211005g - Prozess, Glück, Enttäuschung, Hoffnung - Paul Watzlawick, Anleitung zum Unglücklich sein.txt",
		"IMPORT/210218g - Guter Programmierer, Business, Probleme, Bezahlung - Dave Cheney.txt",
	}
	for _, tc := range tcs {
		if _, err := os.Stat(pathTestRepo + tc); err != nil {
			t.Errorf("File was not created: %+v", tc)
		}
	}
}
