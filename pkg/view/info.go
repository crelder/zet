package view

import (
	"github.com/crelder/zet"
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

	unindexed := getUnindexed(zettel, index)
	infos["unindexed"] = unindexed

	infos["bibkeys"] = AddFrequency(bibkeys)

	return infos
}

func getUnindexed(zettels []zet.Zettel, index zet.Index) []string {
	var rootId string
	var depth int
	m := make(map[string]int)
	for _, zettel := range zettels {
		rootId, depth = getRootIdAndMaxDepth(zettel, zettels)
		if depth == 1 { // There is just one zettel, so no chain of thoughts
			continue
		}
		if isInIndex(rootId, index) {
			continue
		}
		_, ok := m[rootId]
		if !ok {
			m[rootId] = depth
		}
		if ok {
			if m[rootId] < depth {
				m[rootId] = depth
			}
		}
	}

	// Convert m to a slice and sort
	var result []string
	for id, maxDepth := range m {
		result = append(result, id+";"+strconv.Itoa(maxDepth))
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i] > result[j]
	})

	return result

}

func getRootIdAndMaxDepth(zettel zet.Zettel, zettels []zet.Zettel) (string, int) {
	m := make(map[string]zet.Zettel)
	for _, z := range zettels {
		m[z.Id] = z
	}

	currentZettel := zettel
	count := 1
	for {
		if len(currentZettel.Predecessor) == 0 {
			break
		}
		zt, ok := m[currentZettel.Predecessor[0]]
		if !ok {
			break
		}
		currentZettel = zt
		count += 1
	}
	return currentZettel.Id, count
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
