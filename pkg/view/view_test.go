package view

import (
	"github.com/crelder/zet/pkg/parse"
	"github.com/crelder/zet/pkg/transport/repo"
	"github.com/google/go-cmp/cmp"
	"os"
	"path"
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
		// Chain of thoughts branch of. The oldest branch is linear in the same folder.
		// All other branches create a new folder with their subbranch.
		"INDEX/Komplexität/190119e/00_190119e - Komplexität.txt",
		"INDEX/Komplexität/190119e/01_190119d - Testing - clausen2021 87 - 190119e.txt",
		"INDEX/Komplexität/190119e/02_180522a/00_180522a - Komplexität, Thermodynamische Tiefe - 190119d.png",
		"INDEX/Komplexität/190119e/02_180522a/01_220116s/00_220116s - Spezifikation - Marco Fitz - 180522a.pdf",
		"INDEX/Komplexität/190119e/02_180522a/02_210520var - Varietät, Komplexität - 180522a.txt",
		"INDEX/Komplexität/190119e/03_170224a - Polymorphismus, Objektorientierte Programmierung, Schnittstelle - 190119d.png",
		"INDEX/Komplexität/190119e/04_190412d - Presentation, Domain Driven Design, Programmierung - 170224a.txt",

		// An index topic can define more than one entry point into the zettelkasten.
		"INDEX/Komplexität/220122a/00_220122a - Some keyword.txt",

		// Every index entry creates a new folder.
		"INDEX/Programmieren, Objektorientiert/210328obj/00_210328obj - Objektorientiert, Programmierung - kernighan2016 155.pdf",
	}

	for _, tc := range testcases {
		if _, err := os.Stat(pathTestRepo + "/" + tc); err != nil {
			t.Errorf("link was not created: %+v, ", tc)
		}
	}
}

func TestCreateInfo(t *testing.T) {
	// Arrange
	wd, err := os.Getwd()
	if err != nil {
		t.Errorf("could not get the current working dir")
	}
	var pathTestRepo = wd + "/testdata/zettelkasten2"
	parser := parse.New()
	r := repo.New(pathTestRepo, parser)
	viewer := New(r, r)

	// Remove this directory, which might got created in a previous test
	infoPath := pathTestRepo + "/INFO"
	clearPath(infoPath)
	indexPath := pathTestRepo + "/INDEX"
	clearPath(indexPath)

	// Act
	err = viewer.CreateViews()
	if err != nil {
		t.Errorf("Could not generate views: %v", err)
	}

	// Assert
	want := map[string]string{
		"ids.csv":        "170224a;1\n180522a;2\n190119e;1",
		"keywords.csv":   "Complexity;3\nInterface;1\nPolymorphism;1\nTesting;1",
		"context.csv":    "GopherCon;2",
		"references.csv": "clausen2021;1",
		"bibkeys.csv":    "kernighan1999;1\nsedgewick2011;1",

		// 190119e is also unliked, but is references in the index.
		"unlinked.csv": "180522a;1",
	}

	dir, err := os.ReadDir(infoPath)
	if err != nil {
		t.Errorf("error reading dir %v: %v", infoPath, err)
	}

	// map[filename]content
	got := make(map[string]string)
	for _, dirEntry := range dir {
		if dirEntry.Name() == "unindexed.csv" { // For now we will test this in a separate test case
			continue
		}
		file, err := os.ReadFile(path.Join(infoPath, dirEntry.Name()))
		if err != nil {
			t.Errorf("error reading file %v: %v", path.Join(infoPath, dirEntry.Name()), err)
		}
		got[dirEntry.Name()] = string(file)
	}

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf(diff)
	}
}

func TestUnindexed(t *testing.T) {
	// Arrange
	wd, err := os.Getwd()
	if err != nil {
		t.Errorf("could not get the current working dir")
	}
	var pathTestRepo = wd + "/testdata/zettelkasten3"
	parser := parse.New()
	r := repo.New(pathTestRepo, parser)
	viewer := New(r, r)

	// Remove this directory, which might got created in a previous test
	infoPath := pathTestRepo + "/INFO"
	clearPath(infoPath)
	indexPath := pathTestRepo + "/INDEX"
	clearPath(indexPath)

	// Act
	err = viewer.CreateViews()
	if err != nil {
		t.Errorf("Could not generate views: %v", err)
	}

	// Assert
	// There is no index entry for this zettel.
	// There are two chains of thoughts branching of this zettel with id 170224a.
	// It returns the max length, and the total amount of zettel under this branch.
	want := "190119e;4"

	unindexedFileName := "unindexed.csv"
	got, err := os.ReadFile(path.Join(infoPath, unindexedFileName))
	if err != nil {
		t.Errorf("error reading file %v: %v", path.Join(infoPath, unindexedFileName), err)
	}

	if diff := cmp.Diff(string(got), want); diff != "" {
		t.Errorf(diff)
	}
}

func clearPath(path string) {
	err := os.RemoveAll(path)
	if err != nil {
		println("Error occurred: %v", err)
	}
}
