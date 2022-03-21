package bl

import (
	"errors"
	"fmt"
	"math"
	"sort"
	"strings"
	"time"
)

// Zettelkasten contains all entities which are loaded from the repository like
// zettel, index and Literature. It is the central element of the zettelkasten method.
//
// It contains exactly three entities that are explained when introducing these types:
// * Zettel
// * Literature
// * Index
type Zettelkasten struct {
	Zettels []Zettel // Klein, damit private. Nur in Pointer Receiver Methoden
	Index   []Index
}

func NewZk(z []Zettel, i []Index) Zettelkasten {
	return Zettelkasten{
		z,
		i,
	}
}

func (zk Zettelkasten) idExist(id string) bool {
	for _, z := range zk.Zettels {
		if id == z.Id {
			return true
		}
	}
	return false
}

// generateId returns a valid id for a new zettel.
// It makes sure, that the newly generated id is unique,
// therefore does not exist in the zettelkasten yet.
func (zk Zettelkasten) generateId(t time.Time, keywords []string) (string, error) {
	date := t.Format("060102")
	// "010101" is the nil date
	if date == "010101" || len(keywords) == 0 {
		return "", errors.New("date and a keyword is needed for building the Id")
	}

	// Let's try building a unique id via the date and the first letter of one of the keywords
	for i := 0; i < len(keywords); i++ {
		id := date + strings.ToLower(string(keywords[i][0]))
		if zk.idExist(id) {
			continue
		}
		return id, nil
	}
	// We still don't have a keyword. Let's try building a unique id with letters from the alphabet
	for c := 'a'; c < 'z'; c++ {
		id := date + string(c)
		if !zk.idExist(id) {
			return id, nil
		}
		continue
	}

	return "", errors.New("could not build Id from date and first letter of one of the keywords")
}

// GetKeyLinks returns a map with map[keyword][]ids.
// []ids are the ids of zettel that have this keyword.
func (zk Zettelkasten) GetKeyLinks() map[string][]string {
	links := make(map[string][]string)
	for _, z := range zk.Zettels {
		for _, keyword := range z.Keywords {
			links[keyword] = append(links[keyword], z.Id)
		}
	}
	return links
}

// GetConLinks returns a map with map[context][]ids.
// []ids are the ids of zettel that have this context.
func (zk Zettelkasten) GetConLinks() map[string][]string {
	// key is the context, []string are the ids.
	con := map[string][]string{}
	for _, z := range zk.Zettels {
		for _, c := range z.Context {
			con[c] = append(con[c], z.Id)
		}
	}
	return con
}

// GetCitLinks returns a map with map[citation][]ids.
// []ids are the ids of zettel that have this citation.
func (zk Zettelkasten) GetCitLinks() map[string][]string {
	// key is the citation bibkey, []string are the ids.
	bibkeys := map[string][]string{}
	for _, z := range zk.Zettels {
		for _, citation := range z.Citations {
			bibkeys[citation.Bibkey] = append(bibkeys[citation.Bibkey], z.Id)
		}
	}
	return bibkeys
}

// Symlink contains the information needed for both commands in the driven port.
//
// linkName is a relative path. It start within the zettelkasten directory.
// Example:
// Let us assume, that this is the link I want to create,
// `~/zettelkasten/VIEWS/keywords/water/190214d - Water, Fire, Earth.txt`
// which points to the zettel with Id `190214d`.
// Since the zettelkasten lies in `~/zettelkasten`, my linkName would be `VIEWS/keywords/water/`
//
// targetId contains a zettel Id where the link will point to.
type Symlink struct {
	TargetId string
	LinkName string
}

// GetFolgezettel returns symlinks that represent the
// order of zettel in the same way Luhmann had it physically
// in his Zettelkasten. See the test for what the output looks like.
func (zk Zettelkasten) GetFolgezettel(id string) []Symlink {
	symlinks := zk.addSymlink(id, nil, 0, "") // Can I just pass nil here?
	sort.SliceStable(symlinks, func(i, j int) bool {
		return symlinks[i].LinkName < symlinks[j].LinkName
	})
	return symlinks
}

// GetSimilar creates for every zettel (id in map key) similar zettel (values []string containing ids).
func (zk Zettelkasten) GetSimilar() map[string][]string {
	var result = map[string][]string{}
	for _, zettel := range zk.Zettels {
		result[zettel.Id] = append(result[zettel.Id], zk.getSimKeywords(zettel)...)
		result[zettel.Id] = append(result[zettel.Id], zk.getSimCitations(zettel)...)
		result[zettel.Id] = append(result[zettel.Id], zk.getSimContext(zettel)...)
	}

	var resultWithoutDuplicates = map[string][]string{}
	for key, ids := range result {
		resultWithoutDuplicates[key] = removeDuplicates(ids)
	}
	return resultWithoutDuplicates
}

// getSimKeywords takes a zettel and returns the ids of zettel
// that have similar keywords. "Similar keywords" is a simple comparision:
// If there are no similarity matches and empty slice is returned.
func (zk Zettelkasten) getSimKeywords(zettel Zettel) []string {
	if len(zettel.Keywords) == 0 {
		return []string{}
	}

	var result []string
	for _, k := range zettel.Keywords {
		for _, z := range zk.Zettels {
			for _, keyword := range z.Keywords {
				if isSimilarKeyword(k, keyword) {
					result = append(result, z.Id)
				}

			}
		}
	}

	return result
}

// getSimCitations takes a zettel and returns the ids of zettel
// that have similar citations. "Similar citation" is a simple comparision:
// If the bibkey is exactly the same, the two zettel are similar via the citations.
func (zk Zettelkasten) getSimCitations(zettel Zettel) []string {
	if len(zettel.Citations) == 0 {
		return []string{}
	}

	var result []string
	for _, c := range zettel.Citations {
		for _, z := range zk.Zettels {
			for _, citation := range z.Citations {
				if c.Bibkey == citation.Bibkey {
					result = append(result, z.Id)
				}

			}
		}
	}

	return result
}

// getSimContext takes a zettel and returns the ids of zettel
// that have similar context. "Similar context" is a simple comparision:
// If the context is identical, the two zettel are "similar" via the context.
func (zk Zettelkasten) getSimContext(zettel Zettel) []string {
	if len(zettel.Context) == 0 {
		return []string{}
	}

	var result []string
	for _, c := range zettel.Context {
		for _, z := range zk.Zettels {
			for _, context := range z.Context {
				if c == context {
					result = append(result, z.Id)
				}

			}
		}
	}

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

func (zk Zettelkasten) addSymlink(id string, symlinks []Symlink, counter int, path string) []Symlink {
	zettel := zk.getZettel(id)
	if zettel.Id == "" {
		return symlinks
	}
	var newName string
	if path == "" { // This can be done in another function
		newName = fmt.Sprintf("%02d", counter) + "_" + zettel.Name
	} else {
		newName = path + "/" + fmt.Sprintf("%02d", counter) + "_" + zettel.Name
	}
	symlinks = append(symlinks, Symlink{zettel.Id, newName})
	counter++
	if len(zettel.Folgezettel) == 1 {
		symlinks = zk.addSymlink(zettel.Folgezettel[0], symlinks, counter, path)
	}

	if len(zettel.Folgezettel) >= 2 {
		for _, fz := range zettel.Folgezettel[1:] {
			var newPath string
			if path == "" { // This can be done in another function
				newPath = fmt.Sprintf("%02d", counter) + "_" + fz
			} else {
				newPath = path + "/" + fmt.Sprintf("%02d", counter) + "_" + fz
			}
			symlinks = zk.addSymlink(fz, symlinks, 0, newPath)
			counter++
		}
		symlinks = zk.addSymlink(zettel.Folgezettel[0], symlinks, counter, path)
	}

	return symlinks
}

func (zk Zettelkasten) getZettel(id string) Zettel {
	for _, zettel := range zk.Zettels {
		if zettel.Id == id {
			return zettel
		}
	}
	return Zettel{}
}

type ZkErrType int

const (
	NotUniqueId ZkErrType = iota + 1
	TargetLinkIdNotThere
)

type ValidatorErr struct {
	ErrType ZkErrType
	Id      string
}

// Validate returns a slice of ValidatorErr,
// hence inconsistencies in your zettelkasten.
// If there are no errors, it returns nil.
func (zk Zettelkasten) Validate() []ValidatorErr {
	var result []ValidatorErr

	doubleIds := zk.checkDoubleIds()
	for _, dIds := range doubleIds {
		result = append(result, ValidatorErr{
			ErrType: NotUniqueId,
			Id:      dIds,
		})
	}
	deadLinks := zk.getDeadLinks()
	for _, deadLink := range deadLinks {
		result = append(result, ValidatorErr{
			ErrType: TargetLinkIdNotThere,
			Id:      deadLink,
		})
	}

	return result
}

// checkDoubleIds returns all ids that exist more than once in the zettelkasten.
// A double id exists only once in the return string.
func (zk Zettelkasten) checkDoubleIds() []string {
	var ids []string
	for _, zettel := range zk.Zettels {
		ids = append(ids, zettel.Id)
	}

	var doubleids []string
	var counter int
	for _, id := range ids {
		counter = 0
		for _, id2 := range ids {
			if id == id2 {
				counter++
			}
		}
		// If counter is 2 or higher, than an id appears twice and therefore is a duplicate
		if counter > 1 {
			doubleids = append(doubleids, id)
		}
	}
	doubleids = removeDuplicates(doubleids)
	return doubleids
}

// removeDuplicates takes a list of ids and returns a list where every id has only once occurrence.
// If an empty list is provided it returns an empty list.
func removeDuplicates(ids []string) []string {
	if len(ids) == 0 {
		return nil
	}
	keys := make(map[string]bool)
	var uniqueIds []string

	for _, id := range ids {
		if _, value := keys[id]; !value {
			keys[id] = true
			uniqueIds = append(uniqueIds, id)
		}
	}
	return uniqueIds
}

// checkDoubleIds returns all ids that exist more than once in the zettelkasten.
// A double id exists only once in the return string.
func (zk Zettelkasten) getDeadLinks() []string {
	var deadLinks []string
	for _, z := range zk.Zettels {
		for _, fz := range z.Folgezettel {
			if !zk.idExist(fz) {
				deadLinks = append(deadLinks, fz)
			}

		}
	}
	deadLinks = removeDuplicates(deadLinks)

	return deadLinks
}

// ParseToZettel parses the content of a text file into a zettel instance.
func (zk *Zettelkasten) ParseToZettel(fileContent string) (Zettel, error) {
	var z Zettel
	if fileContent == "" {
		return Zettel{}, errors.New("could not parse neither the date, nor the keywords from the header")
	}

	header := getHeader(fileContent)

	date := parseDate(header.date)
	nullTime, _ := time.Parse("06", "")
	if date == nullTime {
		return Zettel{}, errors.New("cannot parse the date")
	}
	z.Keywords = parseKeywordsFromHeader(header.keywords)

	con := parseContext(header.contexts)
	z.Folgezettel = con.folgezettel
	z.Citations = con.lit
	z.Context = con.context

	id, err := zk.generateId(date, z.Keywords)
	if err != nil {
	}
	z.Id = id
	zk.Zettels = append(zk.Zettels, z) // Make sure that we don't give the same id for another nextly iported zettel!

	return z, nil
}

// ParseFileContent parses the content of a text file into a filename.
// Each filename will have a unique id.
func (zk *Zettelkasten) ParseFileContent(fileContent string) (string, error) {
	z, err := zk.ParseToZettel(fileContent)
	if err != nil {
		return "", err
	}
	fn, err := parseToFilename(z)
	if err != nil {
		return "", err
	}
	return fn, nil
}
