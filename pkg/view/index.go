package view

import (
	"fmt"
	"github.com/crelder/zet"
	"path"
)

// getFolgezettelMap contains the business logic for converting the
// tree structure of a zettelkasten into a flat structure in a file directory.
func getFolgezettelMap(zettel []zet.Zettel, index zet.Index) (map[string]string, error) {
	var result = make(map[string]string)
	for topic := range index {
		for _, id := range index[topic] {
			var err error
			result, err = mergeMaps(result, getFolgezettel(id, topic, zettel))
			if err != nil {
				return nil, err
			}
		}
	}
	result2 := make(map[string]string)
	for p, id := range result {
		const indexFolder = "INDEX/"
		result2[indexFolder+p] = id
	}
	return result2, nil
}

// getFolgezettel returns links that represent the
// order of zettel in the same way Luhmann had it physically
// in his Zettelkasten. See the test for what the output looks like.
func getFolgezettel(id, topic string, zettels []zet.Zettel) map[string]string {
	// We need to make sure that this variable is reseted everytime we generate a new Folgezettel structure.
	traveledIds = make(map[string]bool)
	link := map[string]string{}
	result := addLink(id, link, 0, path.Join(topic, id), zettels)
	return result
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

func addLink(id string, links map[string]string, counter int, path string, zettels []zet.Zettel) map[string]string {
	// Make sure that circular links don't end in an endless loop.
	if _, ok := traveledIds[id]; ok {
		return links
	}
	traveledIds[id] = true

	z, err := getZettel(id, zettels)
	if err != nil {
		return links
	}

	newName := path + "/" + fmt.Sprintf("%02d", counter) + "_" + z.Name
	links[newName] = z.Id
	counter++
	if len(z.Folgezettel) == 1 {
		links = addLink(z.Folgezettel[0], links, counter, path, zettels)
	}

	if len(z.Folgezettel) >= 2 {
		for _, fz := range z.Folgezettel[1:] {
			var newPath string
			if path == "" { // This can be done in another function
				newPath = fmt.Sprintf("%03d", counter) + " " + fz
			} else {
				newPath = path + "/" + fmt.Sprintf("%03d", counter) + " " + fz
			}
			links = addLink(fz, links, 0, newPath, zettels)
			counter++
		}
		links = addLink(z.Folgezettel[0], links, counter, path, zettels)
	}

	return links
}

// Make sure that circular links don't end up in an endless loop.
var traveledIds = make(map[string]bool)
