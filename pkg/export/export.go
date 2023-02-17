package export

import (
	"fmt"
	"github.com/crelder/zet"
	"sort"
	"strconv"
	"strings"
)

// ExportPersister creates links to zettel. These links that serve as access points into your zettelkasten.
// The path where the links lie is specified within the ExportPersister.
//
// PersistIndex creates a tree-like link structure, so-called "Folgezettel".
// This represents the physical representation of how Niklas Luhmann arranged his Zettel in his
// wooden zettelkasten boxes. This is used for creating chains of thoughts.

// PersistInfo persists some information like a list of keywords used in your zettelkasten and the number of occurrences.
type ExportPersister interface {
	PersistInfo(m map[string][]byte) error
}

// Exporter contains the application entry point for all operations regarding views upon your zettelkasten.
// Exporter satisfies the zet.Exporter interface.
type Exporter struct {
	Persister ExportPersister
	Repo      zet.Repo
}

func New(ip ExportPersister, r zet.Repo) Exporter {
	return Exporter{
		Persister: ip,
		Repo:      r,
	}
}

// CreateViews creates a folder with different access points (links).
func (e Exporter) Export() error {
	zettel, _, err := e.Repo.GetZettel()
	if err != nil {
		return fmt.Errorf("error creating views: %w", err)
	}
	index, _, err := e.Repo.GetIndex()
	if err != nil {
		return fmt.Errorf("error creating index: %w", err)
	}

	references, err := e.Repo.GetBibkeys()
	if err != nil {
		return err
	}

	// Call method that persists all these info e.InfoPersister.PersistIndex(name, []string).
	// Concrete Implementierung heißt CSVPersister.
	infos, errs := getInfos(zettel, index, references)
	if errs != nil {
		fmt.Println(errs) // TODO: Better error handling
	}

	err = e.Persister.PersistInfo(infos)
	if err != nil {
		return err
	}

	return nil
}

func getZettel(id string, zettel []zet.Zettel) (zet.Zettel, error) {
	// TODO: Make map out of it?
	for _, z := range zettel {
		if z.Id == id {
			return z, nil
		}
	}
	return zet.Zettel{}, fmt.Errorf("export: zettel with id %v not found", id)
}

func getInfos(zettel []zet.Zettel, index zet.Index, bibkeys []string) (map[string][]byte, []error) {
	infos := make(map[string][]byte)

	ids := addFrequency(getIds(zettel))
	if len(ids) > 0 {
		infos["ids.csv"] = ids
	}

	keywords := addFrequency(getKeywords(zettel))
	if len(keywords) > 0 {
		infos["keywords.csv"] = keywords
	}

	context := addFrequency(getContext(zettel))
	if len(context) > 0 {
		infos["context.csv"] = context
	}

	references := addFrequency(getReferences(zettel))
	if len(references) > 0 {
		infos["references.csv"] = references
	}

	pathDepths, errs := getPathDepths(zettel)
	fmt.Println(pathDepths)
	if pathDepths != nil {
		infos["pathDepths.csv"] = convertToByteSlice(pathDepths)
	}

	unindexed := convertToByteSlice(getUnindexedIds(pathDepths, index))
	if unindexed != nil {
		infos["unindexed.csv"] = unindexed
	}

	infos["bibkeys.csv"] = addFrequency(bibkeys)

	// ein großes Json mit zettel, index, references

	// ein gephi

	return infos, errs
}

func getPathDepths(zettels []zet.Zettel) (map[string]int, []error) {
	pathDepths := make(map[string]int)
	var maxDepth int
	var errs []error
	for _, zettel := range zettels {
		rootId, depth, err := getRootAndPathDepth(zettel, zettels)
		if err != nil {
			errs = append(errs, err)
		}
		if depth > maxDepth {
			pathDepths[rootId] = depth
			continue
		}
	}
	return pathDepths, errs
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

func convertToByteSlice(unindexedIds map[string]int) []byte {
	var results []byte
	for s, i := range unindexedIds {
		fmt.Printf("s:%v \ti:%v\n", s, i)
	}
	for id, n := range unindexedIds {
		results = append(results, []byte(fmt.Sprintf("%v;%v", id, n))...)
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i] < results[j]
	})
	return results
}

func getRootAndPathDepth(zettel zet.Zettel, zettels []zet.Zettel) (string, int, error) {
	var (
		z         = zettel
		travelled = make(map[string]bool)
		count     int
		err       error
		ok        bool
	)

	for {
		if _, ok = travelled[z.Id]; ok {
			return "", 0, fmt.Errorf("tracking the root of zettel %v contains a loop", zettel.Id)
		}
		if z.Predecessor == "" { // No Predecessor
			return z.Id, count, nil
		}
		travelled[z.Id] = true                     // Workaround if there are circles in the graph
		z, err = getZettel(z.Predecessor, zettels) // TODO: Handle error!
		if err != nil {
			return z.Id, count, nil
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

func addFrequency(s []string) []byte {
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

	return []byte(strings.Join(result, "\n"))
}
