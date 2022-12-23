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
		"zettel: link to id 160122e not existing": true,
		"zettel: id 180112a not unique":           true,
		"parse filename: more than one predecessor for file \"170327f - More than one predecessor - 180112a, 170311f\"": true,
		"parse filename: could not parse id from filename \"noId.txt\"":                                                 true,
		"index: could not parse line \"Water::170312w\"":                                                                true,
		"index: link to id 180317q not existing":                                                                        true,
		"reference: missing bibkey \"knut2012\"":                                                                        true,
	}

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf(diff)
	}
}
