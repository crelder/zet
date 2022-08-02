package validate

import (
	"github.com/crelder/zet/pkg/parse"
	"github.com/crelder/zet/pkg/transport/repo"
	"github.com/google/go-cmp/cmp"
	"os"
	"testing"
)

func TestValidate(t *testing.T) {
	// Arrange
	wd, err := os.Getwd()
	if err != nil {
		t.Errorf("could not get the current working dir")
	}
	var pathTestRepo = wd + "/testdata/zettelkasten"
	parser := parse.New()
	r := repo.New(pathTestRepo, parser)
	validator := New(r)

	// Act
	inconsErrs, err2 := validator.Val()
	if err2 != nil {
		t.Errorf("Err: %v", err2)
	}

	// Assert
	got := make(map[string]bool)
	for _, inconsErr := range inconsErrs {
		got[inconsErr.Error()] = true
	}

	want := map[string]bool{
		"index: link to id 180317q not existing":     true,
		"reference: missing bibkey \"knut2012\"":     true,
		"zettel: more than one predecessor: 170327f": true,
		"zettel: link to id 170311f not existing":    true,
		"zettel: id 180112a not unique":              true,

		// TODO: works needs to get deleted
		//"works: parse Filename: could not parse id from filename \"noId - Something.txt\"": true,
	}

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf(diff)
	}
}

func TestValidateParsingIdError(t *testing.T) {
	// Arrange
	wd, err := os.Getwd()
	if err != nil {
		t.Errorf("could not get the current working dir")
	}
	var pathTestRepo = wd + "/testdata/zettelkasten2"
	parser := parse.New()
	r := repo.New(pathTestRepo, parser)
	validator := New(r)

	// Act
	_, err = validator.Val()

	// Assert
	if err == nil {
		t.Errorf("Should have received error 'parse Filename: could not parse id from filename \"noId - Something.txt\"', but didn't")
	}

	if err != nil && err.Error() != "parse Filename: could not parse id from filename \"noId - Something.txt\"" {
		t.Errorf("received err string: %v. But expected: not unique id 1801112a", err.Error())
	}
}

func TestValidateParsingIndexError(t *testing.T) {
	// Arrange
	wd, err := os.Getwd()
	if err != nil {
		t.Errorf("could not get the current working dir")
	}
	var pathTestRepo = wd + "/testdata/zettelkasten3"
	parser := parse.New()
	r := repo.New(pathTestRepo, parser)
	validator := New(r)

	// Act
	_, err = validator.Val()

	// Assert
	if err == nil {
		t.Errorf("Should have received error 'parse Filename: could not parse id from filename \"noId - Something.txt\"', but didn't")
	}

	if err != nil && err.Error() != "index: could not parse line \"Water::170312w\"" {
		t.Errorf("received err string: %v. But expected: index: could not parse line \"Water::170312w\"", err.Error())
	}
}
