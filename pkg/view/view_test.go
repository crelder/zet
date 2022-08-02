package view

import (
	"github.com/crelder/zet/pkg/parse"
	"github.com/crelder/zet/pkg/transport/repo"
	"os"
	"testing"
)

func TestCreateIndexViews(t *testing.T) {
	// Arrange
	wd, err := os.Getwd()
	if err != nil {
		t.Errorf("could not get the current working dir")
	}
	var pathTestRepo = wd + "/testdata/zettelkasten"
	parser := parse.New()
	r := repo.New(pathTestRepo, parser)
	viewer := New(r, r)

	// Remove this directory, which might got created in a previous test
	viewPath := pathTestRepo + "/VIEWS"
	clearPath(viewPath)

	// Act
	err = viewer.CreateViews()
	if err != nil {
		t.Errorf("Could not generate views: %v", err)
	}

	// Assert
	// testcases for checking if views are correctly generated
	testcases := []string{
		// index folder contains chains of zettel - these chains can branch (= tree structure)
		"VIEWS/index/Komplexität/190119e/00_190119e - Komplexität.txt",
		"VIEWS/index/Komplexität/190119e/01_190119d/00_190119d - Testing - clausen2021 87 - 190119e.txt",
		"VIEWS/index/Komplexität/190119e/02_180522a - Komplexität, Thermodynamische Tiefe - 190119e.png",
		"VIEWS/index/Komplexität/190119e/03_210520var - Varietät, Komplexität - 180522a.txt",

		"VIEWS/index/Programmieren/220115p/00_220115p - Refactoring, Programmieren - Marco Fitz, clausen2021 5.pdf",
		"VIEWS/index/Programmieren/220115p/01_220116s - Spezifikation - Marco Fitz - 220115p.pdf",

		// These zettel are not linked to other zettel (=missing zettel id at the end of the filename).
		"VIEWS/unlinked/190119e/190119e - Komplexität.txt",
		// There should be suggestions for an unlinked zettel...
		// ...by keywords
		"VIEWS/unlinked/190119e/keywords/180522a - Komplexität, Thermodynamische Tiefe - 190119e.png",
		"VIEWS/unlinked/190119e/keywords/210520var - Varietät, Komplexität - 180522a.txt",

		"VIEWS/unlinked/220115p/220115p - Refactoring, Programmieren - Marco Fitz, clausen2021 5.pdf",
		// ... by reference
		"VIEWS/unlinked/220115p/references/190119d - Testing - clausen2021 87 - 190119e.txt",
		// ... or by context
		"VIEWS/unlinked/220115p/context/220116s - Spezifikation - Marco Fitz - 220115p.pdf",
	}

	for _, tc := range testcases {
		if _, err := os.Stat(pathTestRepo + "/" + tc); err != nil {
			t.Errorf("Symlink was not created: %+v, ", tc)
		}
	}

	// We also need to check that not any more links got created.
	testcases2 := []struct {
		path  string // Path to folder
		count int    // Number of symlinks in this folder
	}{
		{"VIEWS/unlinked/220115p/references", 1},
		{"VIEWS/unlinked/220115p/context", 1},
		{"VIEWS/unlinked/190119e/keywords", 2},
	}

	for _, tc2 := range testcases2 {
		dir, err2 := os.ReadDir(pathTestRepo + "/" + tc2.path)
		if err2 != nil {
			t.Errorf("An error ocurred: %v", err2)
		}
		if len(dir) > tc2.count {
			t.Errorf("Too many symlinks created in path: %v", tc2.path)
		}
	}
}

func clearPath(path string) {
	err := os.RemoveAll(path)
	if err != nil {
		println("Error occurred: %v", err)
	}
}
