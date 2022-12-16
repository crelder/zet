package view

import (
	"github.com/crelder/zet/pkg/parse"
	"github.com/crelder/zet/pkg/transport/fs"
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
	repo := fs.New(pathTestRepo, parser)
	viewer := New(repo, repo)

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
		// Chain of thoughts branch of. The oldest branch is linear in the same folder.
		// All other branches create a new folder with their subbranch.
		"INDEX/Komplexität/190119e/000 190119e - Komplexität.txt",
		"INDEX/Komplexität/190119e/001 190119d - Testing - clausen2021 87 - 190119e.txt",
		"INDEX/Komplexität/190119e/002 180522a/000 180522a - Komplexität, Thermodynamische Tiefe - 190119d.png",
		"INDEX/Komplexität/190119e/002 180522a/001 220116s/000 220116s - Spezifikation - Marco Fitz - 180522a.pdf",
		"INDEX/Komplexität/190119e/002 180522a/002 210520var - Varietät, Komplexität - 180522a.txt",
		"INDEX/Komplexität/190119e/003 170224a - Polymorphismus, Objektorientierte Programmierung, Schnittstelle - 190119d.png",
		"INDEX/Komplexität/190119e/004 190412d - Presentation, Domain Driven Design, Programmierung - 170224a.txt",

		// An index topic can define more than one entry point into the zettelkasten.
		"INDEX/Komplexität/220122a/000 220122a - Some keyword.txt",

		// Every index entry creates a new folder.
		"INDEX/Programmieren, Objektorientiert/210328obj/000 210328obj - Objektorientiert, Programmierung - kernighan2016 155.pdf",
	}

	for _, tc := range testcases {
		if _, err := os.Stat(pathTestRepo + "/" + tc); err != nil {
			t.Errorf("link was not created: %+v, ", tc)
		}
	}
}
