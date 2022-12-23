package validate

import (
	"github.com/crelder/zet/pkg/parse"
	"github.com/crelder/zet/pkg/transport/fs"
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
	repo := fs.New(pathTestRepo, parser)
	validator := New(repo)

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
		"zettel: more than one predecessor: 170327f":   true,
		"zettel: link to id 160122e not existing":      true,
		"zettel: id 180112a not unique":                true,
		"zettel: cannot parse filename 'noId.txt'":     true,
		"index: could not parse line 'Water::170312w'": true,
		"reference: missing bibkey \"pike1989\"":       true,
	}

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf(diff)
	}
}
