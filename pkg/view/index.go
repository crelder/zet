package view

import (
	"fmt"
	"github.com/crelder/zet"
)

// getFolgezettelMap contains the business logic for converting the
// tree structure of a zettelkasten into a flat structure in a file directory.
func getFolgezettelMap(zettel []zet.Zettel, index zet.Index) (map[string]string, error) {
	var result = make(map[string]string)
	for topic := range index {
		for _, id := range index[topic] {
			var err error
			result, err = mergeMaps(result, getFolgezettel(id, zettel))
			if err != nil {
				return nil, err
			}
		}
	}
	result2 := make(map[string]string)
	for path, id := range result {
		result2["INDEX/"+path] = id
	}
	return result2, nil
}

func mergeMaps(result map[string]string, folgezettel map[string]string) (map[string]string, error) {
	for pathName, id := range folgezettel {
		if _, ok := result[pathName]; ok {
			return nil, fmt.Errorf("duplicate pathName %v for id %v", pathName, id)
		}
		result[pathName] = id
	}
	return result, nil
}

func addSymlink(id string, symlinks map[string]string, counter int, path string, zettel []zet.Zettel) map[string]string {
	// Make sure that circular links don't end in an endless loop.
	if _, ok := traveledIds[id]; ok {
		return symlinks
	}
	traveledIds[id] = true

	z := getZettel(id, zettel)
	if z.Id == "" {
		return symlinks
	}
	var newName string
	if path == "" { // This can be done in another function
		newName = fmt.Sprintf("%02d", counter) + "_" + z.Name
	} else {
		newName = path + "/" + fmt.Sprintf("%02d", counter) + "_" + z.Name
	}
	symlinks[newName] = z.Id
	counter++
	if len(z.Folgezettel) == 1 {
		symlinks = addSymlink(z.Folgezettel[0], symlinks, counter, path, zettel)
	}

	if len(z.Folgezettel) >= 2 {
		for _, fz := range z.Folgezettel[1:] {
			var newPath string
			if path == "" { // This can be done in another function
				newPath = fmt.Sprintf("%02d", counter) + "_" + fz
			} else {
				newPath = path + "/" + fmt.Sprintf("%02d", counter) + "_" + fz
			}
			symlinks = addSymlink(fz, symlinks, 0, newPath, zettel)
			counter++
		}
		symlinks = addSymlink(z.Folgezettel[0], symlinks, counter, path, zettel)
	}

	return symlinks
}

// Make sure that circular links don't end up in an endless loop.
var traveledIds = make(map[string]bool)

// getFolgezettel returns symlinks that represent the
// order of zettel in the same way Luhmann had it physically
// in his Zettelkasten. See the test for what the output looks like.
func getFolgezettel(id string, zettel []zet.Zettel) map[string]string {
	// We need to make sure that this variable is reseted everytime we generate a new Folgezettel structure.
	traveledIds = make(map[string]bool)
	sl := map[string]string{}
	result := addSymlink(id, sl, 0, "", zettel)
	return result
}
