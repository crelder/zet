package as

import (
	"github.com/crelder/zet/internal/adapter/driven"
	"os"
	"testing"
)

func TestCreateIndexViews(t *testing.T) {
	// Arrange
	wd, err := os.Getwd()
	if err != nil {
		t.Errorf("could not get the current working dir")
	}
	var pathTestRepo = wd + "/testdata/view/zettelkasten/"
	r := repo.NewRepo(pathTestRepo)
	viewer := NewViewer(r, r)

	// Remove this directory, which might got created in a previous test
	viewPath := pathTestRepo + "VIEWS"
	clearPath(viewPath)

	// Act
	err = viewer.CreateViews()
	if err != nil {
		t.Errorf("Could not generate views: %v", err)
	}

	// Assert
	// testcases for checking if views are correctly generated
	testcases := []string{
		// index folder contains chains of zettel - these chains can can branch (= tree structure)
		"VIEWS/index/Komplexität/190119e/00_190119e - Komplexität - 180522a, 190119d.txt",
		"VIEWS/index/Komplexität/190119e/01_190119d/00_190119d - Testing - clausen2021 87.txt",
		"VIEWS/index/Komplexität/190119e/02_180522a - Komplexität, Thermodynamische Tiefe - 210520var.png",
		"VIEWS/index/Komplexität/190119e/03_210520var - Varietät, Komplexität.txt",

		"VIEWS/index/Programmieren/220115p/00_220115p - Refactoring, Programmieren - Marco Fitz, clausen2021 5 - 220116s.pdf",
		"VIEWS/index/Programmieren/220115p/01_220116s - Spezifikation - Marco Fitz.pdf",

		// All other access points in views contain lists of zettel

		// Check some links that got created for "explore"
		// Here 220115p has similarity through the keyword "Programmieren" (the check is always against the first five letters of a keyword).
		"VIEWS/explore/220115p/210328obj - Objektorientiert, Programmierung - kernighan2016 155 - 170224a.pdf",
		// Here 220115p has similarity through the literature "clausen2021"
		"VIEWS/explore/220115p/190119d - Testing - clausen2021 87.txt",
		// Here 220115p has similarity through the context "Marco Fitz"
		"VIEWS/explore/220115p/220116s - Spezifikation - Marco Fitz.pdf",

		// Check some links that got created for "citations"
		"VIEWS/citations/clausen2021/190119d - Testing - clausen2021 87.txt",
		"VIEWS/citations/clausen2021/220115p - Refactoring, Programmieren - Marco Fitz, clausen2021 5 - 220116s.pdf",

		// Check some links that got created for "context"
		"VIEWS/context/Marco Fitz/220115p - Refactoring, Programmieren - Marco Fitz, clausen2021 5 - 220116s.pdf",
		"VIEWS/context/Marco Fitz/220116s - Spezifikation - Marco Fitz.pdf",

		// Check some links that got created for "keywords"
		"VIEWS/keywords/Komplexität/180522a - Komplexität, Thermodynamische Tiefe - 210520var.png",
		"VIEWS/keywords/Komplexität/190119e - Komplexität - 180522a, 190119d.txt",
		"VIEWS/keywords/Komplexität/210520var - Varietät, Komplexität.txt",
	}

	for _, tc := range testcases {
		if _, err := os.Stat(pathTestRepo + tc); err != nil {
			t.Errorf("Symlink was not created: %+v, ", tc)
		}
	}
}

func clearPath(path string) {
	err := os.RemoveAll(path)
	if err != nil {
		println("Error occured: %v", err)
	}
}
