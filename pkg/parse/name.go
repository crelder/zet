package parse

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/crelder/zet"
	"regexp"
	"strings"
	"time"
)

// Filename parses a filename into a Zettel.
// Only the id is mandatory in a filename. Everything else is optional.
// If you want to provide context, you must at least provide one keyword though.
func Filename(filename string) (zet.Zettel, error) {
	if filename == "" {
		return zet.Zettel{}, errors.New("parse Filename: could not parse empty filename string")
	}
	context := parseContextFromFilename(filename)
	id := parseId(filename)
	if id == "" {
		return zet.Zettel{}, fmt.Errorf("parse Filename: could not parse id from filename %q", filename)
	}
	return zet.Zettel{
		Id:          id,
		Keywords:    parseKeywords(filename),
		Folgezettel: nil,
		Predecessor: context.Predecessor,
		References:  context.References,
		Context:     context.Context,
		Name:        filename,
	}, nil
}

func toFilename(z zet.Zettel) (string, error) {
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

	if len(z.Context) > 0 || len(z.References) > 0 {
		fn += " - "
	}

	if len(z.Context) > 0 {
		fn += strings.Join(z.Context, ", ")
	}

	if len(z.References) > 0 {
		for i, l := range z.References {
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
	// sepLen refers to the length of the separator " - ", which is three characters long.
	const sepLen = 3
	end := strings.Index(filename[start+sepLen:], " - ")
	if end == -1 {
		end = strings.LastIndex(filename, ".")
		keywords = strings.Split(filename[start+sepLen:end], ",")
	} else {
		keywords = strings.Split(filename[start+sepLen:end+start+sepLen], ",")
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
	lines := bytes.Split([]byte(s), []byte("\n"))
	if len(lines) >= 3 {
		header.keywords = string(lines[0])
		header.date = string(lines[1])
		header.contexts = string(lines[2])
		return header
	}
	if len(lines) == 2 {
		header.keywords = string(lines[0])
		header.date = string(lines[1])
		return header
	}
	if len(lines) == 1 && s != "" {
		header.keywords = string(lines[0])
		return header
	}
	return header
}

func parseKeywordsFromHeader(header string) []string {
	var keywords []string
	kw := bytes.Split([]byte(header), []byte(","))
	for _, b := range kw {
		keywords = append(keywords, strings.TrimSpace(string(b)))
	}
	return keywords
}

func parseDate(header string) (time.Time, error) {
	layouts := []string{
		"2.1.06",
		"2.1.2006",
		"060102",
		"January 2, 2006",
		"01/02/06",
	}
	for _, l := range layouts {
		t, err := time.Parse(l, header)
		if err == nil {
			return t, nil
		}
	}
	return time.Time{}, errors.New("parseDate: could not parse date")
}

func parseContext(header string) context {
	if header == "" {
		return context{}
	}
	cleanedHeader := clean(strings.Split(header, ","))

	var con context
	for _, spl := range cleanedHeader {
		if isId(spl) {
			con.Predecessor = spl
			continue
		}
		if isLink(spl) {
			con.Links = append(con.Links, spl)
			continue
		}
		if l := getRef(spl); l.Bibkey != "" {
			con.References = append(con.References, l)
			continue
		}

		if header != "" {
			con.Context = append(con.Context, spl)
		}
	}

	return con
}

// clean removes unnecessary whitespaces.
func clean(s []string) []string {
	var clean []string
	for _, c := range s {
		clean = append(clean, strings.TrimSpace(c))
	}
	return clean
}

func isId(s string) bool {
	r, _ := regexp.Compile("^\\d{6}[a-z]{1,3}$")
	return r.Match([]byte(s))
}

func isLink(s string) bool {
	r, _ := regexp.Compile("^@\\d{6}[a-z]{1,3}$")
	return r.Match([]byte(s))
}

func getRef(spl string) zet.Reference {
	var l zet.Reference

	r, _ := regexp.Compile("^[a-z]{2,}\\d{4}[a-z]?")
	var s = strings.Split(spl, " ")
	l.Bibkey = r.FindString(s[0])
	if len(s) > 1 {
		l.Location = s[1]
	}

	return l
}

type context struct {
	Predecessor string
	Links       []string
	References  []zet.Reference
	Context     []string
}
