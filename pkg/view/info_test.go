package view

import (
	"github.com/crelder/zet/pkg/parse"
	"github.com/crelder/zet/pkg/transport/fs"
	"github.com/google/go-cmp/cmp"
	"os"
	"path"
	"testing"
)

func TestCreateInfo(t *testing.T) {
	// Arrange
	wd, err := os.Getwd()
	if err != nil {
		t.Errorf("could not get the current working dir")
	}
	var pathTestRepo = wd + "/testdata/zettelkasten2"
	parser := parse.New()
	r := fs.New(pathTestRepo, parser)
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
		"ids.csv":        "170224a;1\n180522a;1\n190119e;1", // TODO: Remove Häufigkeit, since an error will be listed in errors.csv
		"keywords.csv":   "Complexity;2\nInterface;1\nPolymorphism;1\nTesting;1",
		"context.csv":    "GopherCon;2",
		"references.csv": "clausen2021;1",
		"bibkeys.csv":    "kernighan1999;1\nsedgewick2011;1",
		"pathDepths.csv": "190119e;2",
		"unindexed.csv":  "190119e;2", // TODO: Also test that something which is a chain but in the index, doesn't show up here.
		// TODO: Check if a circular structur is defended against.
	}

	dir, err := os.ReadDir(infoPath)
	if err != nil {
		t.Errorf("error reading dir %v: %v", infoPath, err)
	}

	// map[filename]content
	got := make(map[string]string)
	for _, dirEntry := range dir {
		file, err := os.ReadFile(path.Join(infoPath, dirEntry.Name()))
		if err != nil {
			t.Errorf("error reading file: %v", err)
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
	r := fs.New(pathTestRepo, parser)
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
	want := "190119e;3"

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
