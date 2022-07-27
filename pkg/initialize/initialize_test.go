package initialize

import (
	"os"
	"testing"
)

func TestInit(t *testing.T) {
	// Arrange
	wd, err := os.Getwd()
	if err != nil {
		t.Errorf("could not get the current working dir")
	}
	var path = wd + "/tmp_test"
	initiator := New(path)

	// Remove this directory, which might got created in a previous test
	clearPath(path)

	// Act
	err = initiator.Init()
	if err != nil {
		t.Errorf("Could not run command initiate: %v", err)
	}

	// Assert
	testcases := []string{
		// Should create these two files...
		"zettelkasten/references.bib",
		"zettelkasten/index.txt",

		// ... and one folder
		"zettelkasten/zettel/",
	}

	for _, tc := range testcases {
		if _, err := os.Stat(path + "/" + tc); err != nil {
			t.Errorf("Not created: %q, ", tc)
		}
	}
}

//func TestInitExample(t *testing.T) {
//	// Arrange
//	wd, err := os.Getwd()
//	if err != nil {
//		t.Errorf("could not get the current working dir")
//	}
//	var path = wd + "/tmp_test_init_example"
//	clearPath(path)
//	initiator := New(path)
//
//	// Act
//	err2 := initiator.InitExample()
//	if err2 != nil {
//		t.Errorf("error occurred: %v", err2)
//	}
//
//	// Assert
//	testcases := []string{
//		// Should create these two files...
//		"zettelkasten/references.bib",
//		"zettelkasten/index.txt",
//
//		// ... and one folder
//		"zettelkasten/zettel/",
//	}
//
//	for _, tc := range testcases {
//		if _, err := os.Stat(path + "/" + tc); err != nil {
//			t.Errorf("Not created: %q, ", tc)
//		}
//	}
//}

func clearPath(path string) {
	err := os.RemoveAll(path)
	if err != nil {
		println("Error occurred: %v", err)
	}
}
