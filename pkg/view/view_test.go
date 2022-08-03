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
	viewPath := pathTestRepo + "/INDEX"
	clearPath(viewPath)

	// Act
	err = viewer.CreateViews()
	if err != nil {
		t.Errorf("Could not generate views: %v", err)
	}

	// Assert
	// testcases for checking if views are correctly generated
	testcases := []string{
		"INDEX/Komplexität/190119e/00_190119e - Komplexität.txt",
		"INDEX/Komplexität/190119e/01_190119d/00_190119d - Testing - clausen2021 87 - 190119e.txt",
		"INDEX/Komplexität/190119e/02_180522a - Komplexität, Thermodynamische Tiefe - 190119e.png",
		"INDEX/Komplexität/190119e/03_210520var - Varietät, Komplexität - 180522a.txt",
		"INDEX/Programmieren/220115p/00_220115p - Refactoring, Programmieren - Marco Fitz, clausen2021 5.pdf",
		"INDEX/Programmieren/220115p/01_220116s - Spezifikation - Marco Fitz - 220115p.pdf",
	}

	for _, tc := range testcases {
		if _, err := os.Stat(pathTestRepo + "/" + tc); err != nil {
			t.Errorf("Symlink was not created: %+v, ", tc)
		}
	}
}

func clearPath(path string) {
	err := os.RemoveAll(path)
	if err != nil {
		println("Error occurred: %v", err)
	}
}
