package view

import (
	"fmt"
	"github.com/crelder/zet"
	"log"
	"sort"
	"strconv"
	"strings"
)

func getInfos(zettel []zet.Zettel, index zet.Index, bibkeys []string) map[string][]string {
	infos := make(map[string][]string)

	ids := AddFrequency(getIds(zettel))
	infos["ids"] = ids

	keywords := AddFrequency(getKeywords(zettel))
	infos["keywords"] = keywords

	context := AddFrequency(getContext(zettel))
	infos["context"] = context

	references := AddFrequency(getReferences(zettel))
	infos["references"] = references

	unlinked := AddFrequency(getUnlinked(zettel, index))
	if unlinked != nil {
		infos["unlinked"] = unlinked
	}

	// Debugg Strategy
	// 1. Change Model to only one Predecessor posssible
	// 2. Add logging around and inside unindexed
	//unindexed := getUnindexed(zettel, index)
	//infos["unindexed"] = unindexed

	infos["bibkeys"] = AddFrequency(bibkeys)

	return infos
}

// getUnindexed returns a list of ids, that are not in the index.
// The method returns for every id the number of follow up zettel it has (Folgezettel), e.g. 190315d;6.
// This means that zettel with id 190315d has 6 Folgezettel.
func getUnindexed(zettels []zet.Zettel, index zet.Index) []string {
	unindexedIds := make(map[string]int)
	var maxDepth int
	for _, zettel := range zettels {
		id, depth := getRootAndMaxLength(zettel, zettels, index)
		if isIrrelevant(id, depth) { // TODO: Rework interface, so that no "" ids are returned.
			continue
		}
		if depth > maxDepth {
			unindexedIds[id] = depth
			continue
		}
	}
	return convertToStringSlice(unindexedIds)
}

func isIrrelevant(id string, depth int) bool {
	return depth == 0 || id == ""
}

func convertToStringSlice(unindexedIds map[string]int) []string {
	var results []string
	for id, n := range unindexedIds {
		results = append(results, fmt.Sprintf("%v;%v", id, n))
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i] < results[j]
	})
	return results
}

func getRootAndMaxLength(zettel zet.Zettel, zettels []zet.Zettel, index zet.Index) (string, int) {
	var (
		z         = zettel
		travelled = make(map[string]bool)
		count     int
		err       error
	)

	for {
		if _, ok := traveledIds[z.Id]; ok {
			return "", 0 // Is it good to handle it like this?
		}
		if z.Predecessor == "" {
			if isInIndex(z.Id, index) {
				return "", 0
			}
			return z.Id, count
		}
		travelled[z.Id] = true
		z, err = getZettel(z.Predecessor, zettels) // TODO: Handle error!
		if err != nil {
			log.Fatalf("%v", err)
		}
		count += 1
	}
}

//func getRootIdAndMaxDepth(zettel zet.Zettel, zettels []zet.Zettel) (string, int) {
//	m := make(map[string]zet.Zettel)
//	for _, z := range zettels {
//		m[z.Id] = z
//	}
//
//	currentZettel := zettel
//	count := 1
//	for {
//		if len(currentZettel.Predecessor) == 0 {
//			break
//		}
//		zt, ok := m[currentZettel.Predecessor[0]]
//		if !ok {
//			break
//		}
//		currentZettel = zt
//		count += 1
//	}
//	return currentZettel.Id, count
//}

func getIds(zettels []zet.Zettel) []string {
	var ids []string
	for _, zettel := range zettels {
		ids = append(ids, zettel.Id)
	}
	return ids
}

func getKeywords(zettel []zet.Zettel) []string {
	var keywords []string
	for _, z := range zettel {
		for _, keyword := range z.Keywords {
			keywords = append(keywords, keyword)
		}
	}
	return keywords
}

func getContext(zettel []zet.Zettel) []string {
	var contexts []string
	for _, z := range zettel {
		for _, context := range z.Context {
			contexts = append(contexts, context)
		}
	}
	return contexts
}

func getReferences(zettel []zet.Zettel) []string {
	var references []string
	for _, z := range zettel {
		for _, reference := range z.References {
			references = append(references, reference.Bibkey)
		}
	}
	return references
}

func getUnlinked(zettels []zet.Zettel, index zet.Index) []string {
	var unlinked []string
	for _, zettel := range zettels {
		if len(zettel.Predecessor) != 0 {
			continue
		}
		if isInIndex(zettel.Id, index) {
			continue
		}

		unlinked = append(unlinked, zettel.Id)
	}
	return unlinked
}

func isInIndex(id string, index zet.Index) bool {
	for _, ids := range index {
		for _, i := range ids {
			if id == i {
				return true
			}
		}
	}
	return false
}

func AddFrequency(s []string) []string {
	m := make(map[string]int)

	for _, elem := range s {
		m[elem] += 1
	}

	var result []string
	for entry, frequency := range m {
		result = append(result, strings.Join([]string{entry, strconv.Itoa(frequency)}, ";"))
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i] < result[j]
	})

	return result
}
