package view

import (
	"fmt"
	"github.com/crelder/zet"
)

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
