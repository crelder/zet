package repo

import (
	"github.com/crelder/zet"
	"github.com/google/go-cmp/cmp"
	"sort"
	"testing"
)

func TestAddFolgezettel(t *testing.T) {
	// Arrange
	tc := []zet.Zettel{
		{
			Id: "190122a",
		},
		{
			Id: "200511b",
		},
		{
			Id:          "211201f",
			Predecessor: []string{"190122a", "200511b"},
		},
		{
			Id:          "221101d",
			Predecessor: []string{"190122a"},
		}}

	// Act
	_, result := addFolgezettel(tc)

	// Assert
	expected := []zet.Zettel{
		{
			Id:          "190122a",
			Folgezettel: []string{"211201f", "221101d"},
		},
		{
			Id:          "200511b",
			Folgezettel: []string{"211201f"},
		},
		{
			Id:          "211201f",
			Predecessor: []string{"190122a", "200511b"},
		},
		{
			Id:          "221101d",
			Predecessor: []string{"190122a"},
		}}

	// Make sure that the order doesn't matter
	sort.Slice(result, func(i, j int) bool {
		return result[i].Id < result[j].Id
	})
	sort.Slice(expected, func(i, j int) bool {
		return expected[i].Id < expected[j].Id
	})

	if diff := cmp.Diff(result, expected); diff != "" {
		t.Errorf(diff)
	}
}
