package validate

import (
	"errors"
	"fmt"
	"github.com/crelder/zet"
	"sort"
)

// Validator analyzes any inconsistencies the zettelkasten has.
// Validator satisfies the zet.Validator interface.
type Validator struct {
	Repo zet.Repo
}

func New(r zet.Repo) Validator {
	return Validator{
		Repo: r,
	}
}

// Val returns all inconsistencies that your zettelkasten has in form of a slice of inconsistencies.
// If the same inconsistency occurs several times, only one is returned, not several.
// If there are none, it returns nil.
func (v Validator) Val() ([]error, error) {
	var incons []error
	zettel, parsErrs, err := v.Repo.GetZettel()
	if err != nil {
		return nil, err
	}
	incons = append(incons, parsErrs...)

	// indexParsingErrors
	index, parsErrs2, err2 := v.Repo.GetIndex()
	if err2 != nil {
		return nil, err2
	}
	incons = append(incons, parsErrs2...)

	bibkeys, err3 := v.Repo.GetBibkeys()
	if err3 != nil {
		return nil, err3
	}

	incons = append(incons, validate(zettel, index, bibkeys)...)
	incons = makeUnique(incons)
	sort.Slice(incons, func(i, j int) bool {
		return incons[i].Error() < incons[j].Error()
	})

	return incons, nil
}

// makeUnique returns a unique list of errors, where a specific error string only occurs once.
func makeUnique(incons []error) []error {
	m := make(map[string]bool)
	for _, i := range incons {
		m[i.Error()] = true
	}

	var result []error
	for str := range m {
		result = append(result, errors.New(str))
	}
	return result
}

// validate returns a slice of inconsistencies.
// If there are no inconsistencies, it returns nil.
func validate(zettel []zet.Zettel, index zet.Index, bibkeys []string) []error {

	var incons []error

	deadLinks := getDeadLinks(zettel)
	for _, deadLink := range deadLinks {
		incons = append(incons, fmt.Errorf("zettel: link to id %v not existing", deadLink))
	}

	deadIndexLinks := getDeadIndexLinks(zettel, index)
	for _, deadIndexLink := range deadIndexLinks {
		incons = append(incons, fmt.Errorf("index: link to id %v not existing", deadIndexLink))
	}

	// Missing Bibkey
	missingBibKeys := getMissingBibKeys(zettel, bibkeys)
	for _, missingBibKey := range missingBibKeys {
		incons = append(incons, fmt.Errorf("reference: missing bibkey %q", missingBibKey))
	}

	// More than one predecessor
	tooManyPredecessorsIds := getTooManyPredecessorsIds(zettel)
	for _, id := range tooManyPredecessorsIds {
		incons = append(incons, fmt.Errorf("zettel: more than one predecessor: %v", id))
	}

	return incons
}

func getTooManyPredecessorsIds(zettel []zet.Zettel) []string {
	var ids []string
	for _, z := range zettel {
		if len(z.Predecessor) > 1 {
			ids = append(ids, z.Id)
		}
	}
	return ids
}

func getMissingBibKeys(zettel []zet.Zettel, bibkeys []string) []string {
	var missing []string
	m := make(map[string]bool)
	for _, bibkey := range bibkeys {
		m[bibkey] = true
	}

	for _, z := range zettel {
		for _, reference := range z.References {
			if _, ok := m[reference.Bibkey]; !ok {
				missing = append(missing, reference.Bibkey)
			}
		}
	}
	return missing
}

func getDeadLinks(zettel []zet.Zettel) []string {
	var deadLinks []string
	for _, z := range zettel {
		for _, pred := range z.Predecessor {
			if !idExist(pred, zettel) {
				deadLinks = append(deadLinks, pred)
			}
		}
	}
	deadLinks = removeDuplicates(deadLinks)

	return deadLinks
}

// removeDuplicates takes a list of ids or keywords and returns a list where every id or keyword has only once occurrence.
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

func idExist(id string, zettel []zet.Zettel) bool {
	for _, z := range zettel {
		if z.Id == id {
			return true
		}
	}
	return false
}

func getDeadIndexLinks(zettel []zet.Zettel, index zet.Index) []string {
	var deadLinks []string
	for _, ids := range index {
		for _, id := range ids {
			if !idExist(id, zettel) {
				deadLinks = append(deadLinks, id)
			}
		}
	}
	deadLinks = removeDuplicates(deadLinks)

	return deadLinks
}
