package view

import (
	"fmt"
	"github.com/crelder/zet"
	"math"
	"sort"
)

// Persister creates symlinks to zettel. These symlinks that serve as access points into your zettelkasten.
// The path where the links lie is specified within the Persister.
//
// CreateSyml creates list of zettel links.
// Prefix is the type of the view (e.g. keywords, literature sources, context).
// m contains as keys the type, e.g. keywords, and as []string one or more ids.
//
// CreateFolgezettelStruct creates a tree-like symlink structure, so-called "Folgezettel".
// This represents the physical representation of how Niklas Luhmann arranged his Zettel in his
// wooden zettelkasten boxes. This is used for creating chains of thoughts.
//
// FileExists checks within your Persister if a file or link was set.
//
// CreateInfo persists some information like a list of keywords used in your zettelkasten and the number of occurrences.
type Persister interface {
	CreateSyml(prefix string, m map[string][]string) error
	CreateFolgezettelStruct(prefix string, links map[string]string) error // links[linkName]targetID
	FileExists(link string) bool
	CreateInfo(prefix string, m map[string][]string) error
}

// Viewer contains the application entry point for all operations regarding views upon your zettelkasten.
// Viewer satisfies the zet.Viewer interface.
type Viewer struct {
	Persister Persister
	Repo      zet.Repo
}

func New(vp Persister, r zet.Repo) Viewer {
	return Viewer{
		Persister: vp,
		Repo:      r,
	}
}

// CreateViews creates a folder with different access points (links).
func (v Viewer) CreateViews() error {
	exists := v.Persister.FileExists("VIEWS")
	if exists {
		return fmt.Errorf(
			"folder 'VIEWS' exists in the zettelkasten.\nPlease manually delete folder 'VIEWS' in order to execute the command 'views'.")
	}

	zettel, _, err := v.Repo.GetZettel()
	if err != nil {
		return fmt.Errorf("error creating views: %w", err)
	}
	index, _, err2 := v.Repo.GetIndex()
	if err2 != nil {
		return fmt.Errorf("error creating index: %w", err2)
	}

	err3 := v.createInfos(zettel)
	if err3 != nil {
		return err3
	}

	//Creating index links
	for topic := range index {
		for _, id := range index[topic] {
			links := getFolgezettel(id, zettel)
			err6 := v.Persister.CreateFolgezettelStruct(topic, links)
			if err6 != nil {
				return err6
			}
		}
	}

	// Creating unlinked links
	m := make(map[string][]string)
	unlinked := getUnlinked(20, zettel)
	for _, u := range unlinked {
		m["unlinked/"+u] = []string{u}
		s := getSimKeywords(getZettel(u, zettel), zettel)
		s = substract(s, u)
		if s != nil {
			m["unlinked/"+u+"/keywords"] = s
		}
		s2 := getSimReferences(getZettel(u, zettel), zettel)
		s2 = substract(s2, u)
		if s2 != nil {
			m["unlinked/"+u+"/references"] = s2
		}
		s3 := getSimContext(getZettel(u, zettel), zettel)
		s3 = substract(s3, u)
		if s3 != nil {
			m["unlinked/"+u+"/context"] = s3
		}
	}
	err7 := v.Persister.CreateSyml("", m)
	if err7 != nil {
		return err7
	}

	return nil
}

// substract removes u from s.
// In case s is nil, it returns nil.
// In case u is an empty string, it return s.
func substract(s []string, u string) []string {
	var result []string
	for _, s2 := range s {
		if s2 == u {
			continue
		}
		result = append(result, s2)
	}
	return result
}

func (v Viewer) createInfos(zettel []zet.Zettel) error {
	kv := getKeyLinks(zettel)
	err3 := v.Persister.CreateInfo("keywords", kv)
	if err3 != nil {
		return fmt.Errorf("error creating info for keywords: %w", err3)
	}

	con := getConLinks(zettel)
	err4 := v.Persister.CreateInfo("context", con)
	if err4 != nil {
		return err4
	}

	cit := getRefLinks(zettel)
	err5 := v.Persister.CreateInfo("references", cit)
	if err5 != nil {
		return err5
	}
	return nil
}

// getKeyLinks returns a map with map[keyword][]ids.
// []ids are the ids of zettel that have this keyword.
func getKeyLinks(zettel []zet.Zettel) map[string][]string {
	links := make(map[string][]string)
	for _, z := range zettel {
		for _, keyword := range z.Keywords {
			links[keyword] = append(links[keyword], z.Id)
		}
	}
	return links
}

// getConLinks returns a map with map[context][]ids.
// []ids are the ids of zettel that have this context.
func getConLinks(zettel []zet.Zettel) map[string][]string {
	// key is the context, []string are the ids.
	con := map[string][]string{}
	for _, z := range zettel {
		for _, c := range z.Context {
			con[c] = append(con[c], z.Id)
		}
	}
	return con
}

// getRefLinks returns a map with map[reference][]ids.
// []ids are the ids of zettel that have this reference.
func getRefLinks(zettel []zet.Zettel) map[string][]string {
	// key is the reference bibkey, []string are the ids.
	bibkeys := map[string][]string{}
	for _, z := range zettel {
		for _, reference := range z.References {
			bibkeys[reference.Bibkey] = append(bibkeys[reference.Bibkey], z.Id)
		}
	}
	return bibkeys
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
	if len(result) > 0 {
		sl["root"] = id // This derives from how the map should be filled: with links; this is a workaround forgenerating the correct link structure.
	}
	return result
}

// getSimKeywords takes a zettel and returns the ids of zettel
// that have similar keywords.
//
// "Similar keywords" is a simple comparison:
// if the first five letters between two keywords match, it is considered "similar" by keyword.
// If one of the to-be-compared keywords is shorter than five letters,
// the two keywords get compared just with the length of the shortest keyword.
//
// The returned ids are sorted alphabetically and made unique (one id doesn't appear more than once).
//
// If there are no similarity matches and empty slice is returned.
func getSimKeywords(zettel zet.Zettel, zettels []zet.Zettel) []string {
	if len(zettel.Keywords) == 0 {
		return []string{}
	}

	var result []string
	for _, k := range zettel.Keywords {
		for _, z := range zettels {
			for _, keyword := range z.Keywords {
				if isSimilarKeyword(k, keyword) {
					result = append(result, z.Id)
				}

			}
		}
	}

	result = removeDuplicates(result)

	sort.Slice(result, func(i, j int) bool {
		return result[i] < result[j]
	})

	return result
}

// removeDuplicates takes a list of strgs (ids or keywords) and returns a list where every string appears only once.
// If an empty list is provided it returns an empty list.
func removeDuplicates(strgs []string) []string {
	if len(strgs) == 0 {
		return nil
	}
	idMap := make(map[string]bool, len(strgs))
	var uniqueIds []string

	for _, id := range strgs {
		if _, ok := idMap[id]; !ok {
			idMap[id] = true
			uniqueIds = append(uniqueIds, id)
		}
	}
	return uniqueIds
}

// getSimReferences takes a zettel and returns the ids of zettel
// that have similar references. "Similar references" is a simple comparison:
// If two zettel have at least one bibkey in common, the two zettel are similar via the references.
func getSimReferences(zettel zet.Zettel, zettels []zet.Zettel) []string {
	if len(zettel.References) == 0 {
		return []string{}
	}

	var result []string
	for _, c := range zettel.References {
		for _, z := range zettels {
			for _, reference := range z.References {
				if c.Bibkey == reference.Bibkey {
					result = append(result, z.Id)
				}

			}
		}
	}

	result = removeDuplicates(result)

	sort.Slice(result, func(i, j int) bool {
		return result[i] < result[j]
	})

	return result
}

// getSimContext takes a zettel and returns the ids of zettel
// that have similar context. "Similar context" is a simple comparison:
// If the context is identical, the two zettel are "similar" via the context.
func getSimContext(zettel zet.Zettel, zettels []zet.Zettel) []string {
	if len(zettel.Context) == 0 {
		return []string{}
	}

	var result []string
	for _, c := range zettel.Context {
		for _, z := range zettels {
			for _, context := range z.Context {
				if c == context {
					result = append(result, z.Id)
				}

			}
		}
	}

	result = removeDuplicates(result)

	sort.Slice(result, func(i, j int) bool {
		return result[i] < result[j]
	})

	return result
}

// isSimilarKeyword returns true if the first five letters between two keywords match.
// If one of the to be compared keywords is shorter than five letters,
// the two keywords get compared just with the length of the shortest keyword.
// In all other cases it returns false.
func isSimilarKeyword(kw1, kw2 string) bool {
	if len(kw1) < 5 || len(kw2) < 5 {
		minLen := int(math.Min(float64(len(kw1)), float64(len(kw2))))
		return kw1[:minLen] == kw2[:minLen]
	}
	return kw1[:5] == kw2[:5]
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

func getZettel(id string, zettel []zet.Zettel) zet.Zettel {
	for _, z := range zettel {
		if z.Id == id {
			return z
		}
	}
	return zet.Zettel{}
}

// getUnlinked gets the n lastest (zettel are sorted ascending by date) ids of zettel, that have no reference
// to another zettel.
func getUnlinked(n int, zettel []zet.Zettel) []string {
	var unlinked []string
	if len(zettel) < n {
		n = len(zettel)
	}
	for _, zettel := range zettel[len(zettel)-n:] {
		if len(zettel.Predecessor) == 0 {
			unlinked = append(unlinked, zettel.Id)
		}
	}

	return unlinked
}
