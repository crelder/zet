package bl

import (
	"bytes"
	"errors"
	"regexp"
	"strings"
	"time"
)

// Zettel holds one thought.
type Zettel struct {
	Id          string
	Keywords    []string
	Folgezettel []string // folgezettel ids
	Citations   []citation
	Context     []string
	Name        string // the filename, e.g. '170212g - Golang.txt'
}

func parseToFilename(z Zettel) (string, error) {
	var fn string
	const fileExt = ".txt"

	if z.Id == "" {
		return "", errors.New("id is missing, but needed for creation of filename")
	}
	fn += z.Id

	if len(z.Keywords) == 0 {
		return "", errors.New("at least one keyword is needed for creation of filename")
	}
	fn += " - " + strings.Join(z.Keywords, ", ")

	if len(z.Context) > 0 || len(z.Citations) > 0 {
		fn += " - "
	}

	if len(z.Context) > 0 {
		fn += strings.Join(z.Context, ", ")
	}

	if len(z.Citations) > 0 {
		for i, l := range z.Citations {
			if i == 0 && len(z.Context) == 0 {
				fn += l.Bibkey
				if l.Location != "" {
					fn += " " + l.Location
				}
			} else {
				fn += ", " + l.Bibkey
				if l.Location != "" {
					fn += " " + l.Location
				}
			}
		}
	}

	if len(z.Folgezettel) > 0 {
		fn += " - "
		fn += strings.Join(z.Folgezettel, ", ")
	}

	fn += fileExt

	return fn, nil
}

// ParseFilename parses a filename into a Zettel.
// If the mandatory value for id could not get parsed,
// an error and an empty Zettel{} is returned.
func ParseFilename(filename string) (Zettel, error) {
	if filename == "" {
		return Zettel{}, errors.New("could not parse empty string")
	}
	context := parseContextFromFilename(filename)
	id := parseId(filename)
	if id == "" {
		return Zettel{}, errors.New("could not parse id")
	}
	return Zettel{
		Id:          id,
		Keywords:    parseKeywords(filename),
		Folgezettel: context.folgezettel,
		Citations:   context.lit,
		Context:     context.context,
		Name:        filename,
	}, nil
}

func parseId(filename string) string {
	r, _ := regexp.Compile("^\\d{6}[a-z]{1,3}")
	id := r.FindString(filename)
	return id
}

func parseKeywords(filename string) []string {
	if filename == "" {
		return nil
	}
	// If the separator " - " does not exist, there can not be any keywords
	start := strings.Index(filename, " - ")
	if start == -1 {
		return nil
	}

	var keywords []string
	end := strings.Index(filename[start+3:], " - ")
	if end == -1 {
		end = strings.LastIndex(filename, ".")
		keywords = strings.Split(filename[start+3:end], ",")
	} else {
		keywords = strings.Split(filename[start+3:end+start+3], ",")
	}

	for i, keyword := range keywords {
		keywords[i] = strings.TrimLeft(keyword, " ")
	}
	return keywords
}

func parseContextFromFilename(fn string) context {
	start := strings.Index(fn, " - ")
	relevantString := fn[start+3 : len(fn)-4]
	pos := strings.Index(relevantString, " - ")
	if pos == -1 {
		return context{}
	}
	relevantString = relevantString[pos+3:]
	split := strings.Split(relevantString, " - ")
	s := strings.Join(split, ",")
	return parseContext(s)
}

type header struct {
	keywords string
	date     string
	contexts string
}

func getHeader(s string) header {
	var header header
	l := bytes.Split([]byte(s), []byte("\n"))
	if len(l) >= 3 {
		header.keywords = string(l[0])
		header.date = string(l[1])
		header.contexts = string(l[2])
		return header
	}
	if len(l) == 2 {
		header.keywords = string(l[0])
		header.date = string(l[1])
		return header
	}
	if len(l) == 1 && s != "" {
		header.keywords = string(l[0])
		return header
	}
	return header
}

func parseKeywordsFromHeader(s string) []string {
	var keywords []string
	kw := bytes.Split([]byte(s), []byte(","))
	for _, b := range kw {
		keywords = append(keywords, strings.TrimSpace(string(b)))
	}
	return keywords
}

func parseDate(s string) time.Time {
	//Layout      = "01/02 03:04:05PM '06 -0700" // The reference time, in numerical order.
	layouts := []string{
		"2.1.06",
		"2.1.2006",
		"060102",
		"January 2, 2006",
		"01/02/06",
	}
	for _, l := range layouts {
		t, err := time.Parse(l, s)
		if err == nil {
			return t
		}
	}
	// This is the nil time and can be checked via IsZero()
	return time.Date(1, time.Month(1), 1, 0, 0, 0, 0, time.UTC)
}

// citation has a bibkey which refers to a Literature entity (e.g. book, paper, etc.).
// With a Literature entity and the location (e.g. page number, chapter, etc.) you can exactly define the source.
type citation struct {
	Bibkey   string
	Location string
}

type context struct {
	folgezettel []string
	links       []string
	lit         []citation
	context     []string
}

func parseContext(s string) context {
	if s == "" {
		return context{}
	}
	var splitted = split(s)

	var con context
	// Get Bibkeys
	for _, spl := range splitted {
		if isFolgezettel(spl) {
			con.folgezettel = append(con.folgezettel, spl)
			continue
		}
		if isLink(spl) {
			con.links = append(con.links, spl)
			continue
		}
		if l := getLit(spl); l.Bibkey != "" {
			con.lit = append(con.lit, l)
			continue
		}

		if s != "" {
			con.context = append(con.context, spl)
		}
	}

	return con
}

func split(s string) []string {
	var splitted []string
	for _, c := range strings.Split(s, ",") {
		splitted = append(splitted, strings.TrimSpace(c))
	}
	return splitted
}

func isFolgezettel(s string) bool {
	r, _ := regexp.Compile("^\\d{6}[a-z]{1,3}$")
	return r.Match([]byte(s))
}

func isLink(s string) bool {
	r, _ := regexp.Compile("^@\\d{6}[a-z]{1,3}$")
	return r.Match([]byte(s))
}

func getLit(spl string) citation {
	var l citation

	r, _ := regexp.Compile("^[a-zA-Z]{2,}\\d{4}[a-z]?")
	var s = strings.Split(spl, " ")
	l.Bibkey = r.FindString(s[0])
	if len(s) > 1 {
		l.Location = s[1]
	}

	return l
}
