package export

import (
	"encoding/json"
	"encoding/xml"
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
	PersistInfo(m map[string][]string) error
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

func getInfos(zettel []zet.Zettel, index zet.Index, bibkeys []string) (map[string][]string, []error) {
	infos := make(map[string][]string)

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
	if pathDepths != nil {
		infos["pathDepths.csv"] = convertToStringSlice(pathDepths)
	}

	unindexed := convertToStringSlice(getUnindexedIds(pathDepths, index))
	if unindexed != nil {
		infos["unindexed.csv"] = unindexed
	}

	infos["bibkeys.csv"] = addFrequency(bibkeys)

	j := getJson(zettel, index, bibkeys)
	if j != nil {
		infos["zettelkasten.json"] = j
	}

	g, _ := getGephi(zettel)
	if g != nil {
		infos["zettelkasten.gephi"] = g
	}

	return infos, errs
}

type Gexf struct {
	XMLName xml.Name `xml:"gexf"`
	Xmlns   string   `xml:"xmlns,attr"`
	Graph   Graph    `xml:"graph"`
}

type Graph struct {
	XMLName         xml.Name `xml:"graph"`
	DefaultEdgeType string   `xml:"defaultedgetype,attr"`
	IdType          string   `xml:"idtype,attr"`
	Type            string   `xml:"type,attr"`
	Nodes           Nodes    `xml:"nodes"`
	Edges           Edges    `xml:"edges"`
}

type Nodes struct {
	XMLName xml.Name `xml:"nodes"`
	Count   int      `xml:"count,attr"`
	Nodes   []Node   `xml:"node"`
}

type Node struct {
	XMLName xml.Name `xml:"node"`
	Id      string   `xml:"id,attr"`
	Label   string   `xml:"label,attr"`
}

type Edges struct {
	XMLName xml.Name `xml:"edges"`
	Count   int      `xml:"count,attr"`
	Edges   []Edge   `xml:"edge"`
}

type Edge struct {
	XMLName xml.Name `xml:"edge"`
	Id      int      `xml:"id,attr"`
	Source  string   `xml:"source,attr"`
	Target  string   `xml:"target,attr"`
}

func getGephi(zettel []zet.Zettel) ([]string, error) {
	var result = []string{}

	var n = []Node{}
	var e = []Edge{}
	var i = 1
	for _, z := range zettel {
		n = append(n, Node{
			Id:    z.Id,
			Label: strings.Join(z.Keywords, ", "),
		})
		i++

		for _, f := range z.Folgezettel {
			e = append(e, Edge{
				Id:     i,
				Source: z.Id,
				Target: f,
			})
			i++
		}
	}

	var edges = Edges{
		Count: i,
		Edges: e,
	}

	var nodes = Nodes{
		Count: len(zettel),
		Nodes: n,
	}

	var graph = Graph{
		DefaultEdgeType: "directed",
		IdType:          "string",
		Type:            "static",
		Nodes:           nodes,
		Edges:           edges,
	}
	var gexf = Gexf{
		Xmlns: "xmlns=\"http://gexf.net/1.2\" version=\"1.2\"",
		Graph: graph,
	}

	xmlString, err := xml.Marshal(gexf)
	if err != nil {
		return nil, err
	}

	result = append(result, "<?xml version=\"1.0\" encoding=\"UTF-8\"?>")
	split := strings.Split(string(xmlString), "\n")
	for _, s := range split {
		result = append(result, s)
	}

	return result, nil
}

func getJson(zettel []zet.Zettel, index zet.Index, bibkeys []string) []string {
	zk := struct {
		Zettel  []zet.Zettel `json:"zettel"`
		Index   zet.Index    `json:"index"`
		Bibkeys []string     `json:"bibkeys"`
	}{
		zettel,
		index,
		bibkeys,
	}
	zkJson, _ := json.MarshalIndent(zk, "", "\t")
	return strings.Split(fmt.Sprintf("%s", zkJson), "\n")
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

func addFrequency(s []string) []string {
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
