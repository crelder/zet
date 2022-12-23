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
	if len(ids) > 0 {
		infos["ids"] = ids
	}

	keywords := AddFrequency(getKeywords(zettel))
	if len(keywords) > 0 {
		infos["keywords"] = keywords
	}

	context := AddFrequency(getContext(zettel))
	if len(context) > 0 {
		infos["context"] = context
	}

	references := AddFrequency(getReferences(zettel))
	if len(references) > 0 {
		infos["references"] = references
	}

	//pathDepths := getPathDepths(zettel)
	//if pathDepths != nil {
	//	infos["pathDepths"] = convertToStringSlice(pathDepths)
	//}

	//unindexed := convertToStringSlice(getUnindexedIds(pathDepths, index))
	//if unindexed != nil {
	//	infos["unindexed"] = unindexed
	//}

	infos["bibkeys"] = AddFrequency(bibkeys)

	return infos
}

func getPathDepths(zettels []zet.Zettel) map[string]int {
	pathDepths := make(map[string]int)
	var maxDepth int
	for _, zettel := range zettels {
		rootId, depth := getRootAndPathDepth(zettel, zettels)
		if depth > maxDepth {
			pathDepths[rootId] = depth
			continue
		}
	}
	return pathDepths
}

// getUnindexedIds returns a list of ids, that are not in the index.
// The method returns for every id the number of follow up zettel it has (Folgezettel), e.g. 190315d;6.
// This means that zettel with id 190315d has 6 Folgezettel.
func getUnindexedIds(pathDepths map[string]int, index zet.Index) map[string]int {
	unindexedIds := make(map[string]int)
	for id, depth := range pathDepths {
		if isInIndex(id, index) {
			continue
		}
		unindexedIds[id] = depth
	}
	return unindexedIds
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

func getRootAndPathDepth(zettel zet.Zettel, zettels []zet.Zettel) (string, int) {
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
		if z.Predecessor == "" { // No Predecessor
			return z.Id, count
		}
		travelled[z.Id] = true                     // Workaround if there are circles in the graph
		z, err = getZettel(z.Predecessor, zettels) // TODO: Handle error!
		if err != nil {
			log.Fatalf("%v", err)
		}
		count += 1
	}
}

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
