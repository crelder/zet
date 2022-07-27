package validate

import (
	"github.com/crelder/zet/pkg/parse"
	"github.com/crelder/zet/pkg/transport/repo"
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
		"zettel: not unique id \"180112a\"":                                                 true,
		"zettel: target link not existing: 190311d":                                         true,
		"index: target link not existing: 180317q":                                          true,
		"index: could not parse line \"Water::170312w\"":                                    true,
		"zettel: predecessor id 190311d doesn't exist, zettel 170327f":                      true,
		"reference: missing bibkey \"knut2012\"":                                            true,
		"zettel: parse Filename: could not parse id from filename \"noId - Something.txt\"": true,
		"works: parse Filename: could not parse id from filename \"noId - Something.txt\"":  true,
	}

	for str := range want {
		if _, ok := got[str]; !ok {
			t.Errorf("missing inconsistency %q", str)
		}
	}

	// Check if we didn't get too many errors.
	if len(want) < len(inconsErrs) {
		t.Errorf("Got too many inconsistencies: got %v, want %v", len(inconsErrs), len(want))
	}

}
