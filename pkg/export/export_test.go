package export

import (
	"encoding/json"
	"encoding/xml"
	"github.com/crelder/zet"
	"github.com/crelder/zet/pkg/parse"
	"github.com/crelder/zet/pkg/transport/fs"
	"github.com/google/go-cmp/cmp"
	"os"
	"path"
	"testing"
)

func TestCreateExport(t *testing.T) {
	// Arrange
	wd, err := os.Getwd()
	if err != nil {
		t.Errorf("could not get the current working dir")
	}
	var pathTestRepo = wd + "/testdata/zettelkasten"
	parser := parse.New()
	r := fs.New(pathTestRepo, parser)
	viewer := New(r, r)

	// Remove this directory, which might got created in a previous test
	infoPath := pathTestRepo + "/EXPORT"
	clearPath(infoPath)

	// Act
	err = viewer.Export()
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

	// Create a map[filename]content from the files created
	got := make(map[string]string)
	for _, dirEntry := range dir {
		if dirEntry.Name() == "zettelkasten.json" || dirEntry.Name() == "zettelkasten.gexf" { // This will be checked in another part of this test.
			continue
		}
		file, err := os.ReadFile(path.Join(infoPath, dirEntry.Name()))
		if err != nil {
			t.Errorf("error reading file: %v", err)
		}
		got[dirEntry.Name()] = string(file)
	}

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf(diff)
	}

	for _, dirEntry := range dir {
		if dirEntry.Name() != "zettelkasten.json" {
			continue
		}

		expected := getExpectedJson()
		var actual = zettelkasten{}

		file, err := os.ReadFile(path.Join(infoPath, dirEntry.Name()))
		if err != nil {
			t.Errorf("error reading file: %v", err)
		}
		_ = json.Unmarshal([]byte(file), &actual)
		if diff := cmp.Diff(actual, expected); diff != "" {
			t.Errorf(diff)
		}

	}

	for _, dirEntry := range dir {
		if dirEntry.Name() != "zettelkasten.gexf" {
			continue
		}

		expected := getExpectedGephi()
		var actual = Gexf{}

		file, err := os.ReadFile(path.Join(infoPath, dirEntry.Name()))
		if err != nil {
			t.Errorf("error reading file: %v", err)
		}

		_ = xml.Unmarshal([]byte(file), &actual)
		if diff := cmp.Diff(actual, expected); diff != "" {
			t.Errorf(diff)
		}
	}
}

func TestUnindexed(t *testing.T) {
	// Arrange
	wd, err := os.Getwd()
	if err != nil {
		t.Errorf("could not get the current working dir")
	}
	var pathTestRepo = wd + "/testdata/zettelkasten2"
	parser := parse.New()
	r := fs.New(pathTestRepo, parser)
	exporter := New(r, r)

	// Remove this directory, which might got created in a previous test
	exportPath := pathTestRepo + "/EXPORT"
	clearPath(exportPath)

	// Act
	err = exporter.Export()
	if err != nil {
		t.Errorf("Could not generate views: %v", err)
	}

	// Assert
	// There is no index entry for this zettel.
	// There are two chains of thoughts branching of this zettel with id 170224a.
	// It returns the max length, and the total amount of zettel under this branch.
	want := "190119e;3"

	unindexedFileName := "unindexed.csv"
	got, err := os.ReadFile(path.Join(exportPath, unindexedFileName))
	if err != nil {
		t.Errorf("error reading file %v: %v", path.Join(exportPath, unindexedFileName), err)
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

type zettelkasten struct {
	Zettel  []zet.Zettel
	Index   []zet.Index
	Bibkeys []string
}

func getExpectedJson() zettelkasten {
	var expected = `{
	"zettel": [
		{
			"Id": "170224a",
			"Keywords": [
				"Polymorphism",
				"Interface"
			],
			"Folgezettel": [
				"180522a"
			],
			"Predecessor": "190119e",
			"References": null,
			"Context": [
				"GopherCon"
			],
			"Name": "170224a - Polymorphism, Interface - GopherCon - 190119e.png"
		},
		{
			"Id": "180522a",
			"Keywords": [
				"Testing",
				"Complexity"
			],
			"Folgezettel": null,
			"Predecessor": "170224a",
			"References": [
				{
					"Bibkey": "clausen2021",
					"Location": "87"
				}
			],
			"Context": null,
			"Name": "180522a - Testing, Complexity - clausen2021 87 - 170224a.txt"
		},
		{
			"Id": "190119e",
			"Keywords": [
				"Complexity"
			],
			"Folgezettel": [
				"170224a"
			],
			"Predecessor": "",
			"References": null,
			"Context": [
				"GopherCon"
			],
			"Name": "190119e - Complexity - GopherCon.txt"
		}
	],
	"index": {
		"Complexity": [
			"220122a"
		]
	},
	"bibkeys": [
		"kernighan1999",
		"sedgewick2011"
	]
}`
	var out = zettelkasten{}
	_ = json.Unmarshal([]byte(expected), &out)
	return out
}

func getExpectedGephi() Gexf {
	a := `<?xml version="1.0" encoding="UTF-8"?>
<gexf xmlns="http://gexf.net/1.2" version="1.2">
  <meta lastmodifieddate="2010-03-03+23:44">
    <creator>zet</creator>
    <description>zettelkasten</description>
  </meta>
  <graph defaultedgetype="directed" idtype="string" type="static">
    <nodes count="2">
      <node id="170915a" label="Neuronales Netz, Features und Beispiele, Nicht-numersiche Daten " />
      <node id="171109a" label="Pflanzen, Signal, Sensibilität " />
    </nodes>
    <edges count="1">
      <edge id="1" source="171109a" target="170915a" />
    </edges>
  </graph>
</gexf>`
	var expected Gexf
	xml.Unmarshal([]byte(a), &expected)
	return expected
}
