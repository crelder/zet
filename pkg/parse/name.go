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
	id := parseId(filename)
	if id == "" {
		return zet.Zettel{}, fmt.Errorf("parse Filename: could not parse id from filename %q", filename)
	}
	context, incon := parseContextFromFilename(filename)
	return zet.Zettel{
		Id:          id,
		Keywords:    parseKeywords(filename),
		Folgezettel: nil, // Folgezettel wegmachen und nur bei der Indexgenerierung hinzu, dann in ein zweites embedded struct packen?
		Predecessor: context.Predecessor,
		References:  context.References,
		Context:     context.Context,
		Name:        filename,
	}, incon
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

	if len(z.Predecessor) > 0 {
		fn += " - "
		fn += z.Predecessor
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

func parseContextFromFilename(fn string) (context, error) {
	if fn == "" {
		return context{}, nil
	}
	parts := strings.Split(fn, " - ")
	if len(parts) == 1 {
		// The filename only consists of an id.
		return context{}, nil
	}
	if len(parts) == 2 {
		if isId(parts[1]) {
			// The filename consists of an id and a predecessor id.
			return context{
				Predecessor: parts[1],
				References:  nil,
				Context:     nil,
			}, nil

		}
		// The filename consists of an id and keywords.
		return context{
			Predecessor: "",
			References:  nil,
			Context:     nil,
		}, nil
	}
	if len(parts) == 3 {
		if isId(parts[2]) {
			// We don't have context, only id - keywords - predecessor
			return context{
				Predecessor: parts[2],
				References:  nil,
				Context:     nil,
			}, nil
		}
		// We have id - keywords - context
		c, parseErr := parseContext2(parts[2])
		return context{
			Predecessor: "",
			References:  c.References,
			Context:     c.Context,
		}, parseErr

	}
	if len(parts) == 4 {
		// All parts are filled
		var p string
		if isId(parts[3]) {
			p = parts[3]
		}
		c, parseErr := parseContext2(parts[2])
		return context{
			Predecessor: p,
			References:  c.References,
			Context:     c.Context,
		}, parseErr

	}
	return context{}, fmt.Errorf("parse name: there shouldn't be more than four parts in a filename")
}

type header struct {
	keywords string
	date     string
	contexts string
}

func getHeader(s string) header {
	var header header
	lines := bytes.Split([]byte(s), []byte("\n"))
	for i, _ := range lines {
		if i == 0 {
			header.keywords = string(lines[0])
			continue
		}
		if i == 1 {
			header.date = string(lines[1])
			continue
		}
		if i == 2 {
			header.contexts = string(lines[2])
			continue
		}
		break
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

func parseContext(line string) (context, error) {
	if line == "" {
		return context{}, nil
	}
	cleanedLine := clean(strings.Split(line, ","))

	var con context
	for _, elem := range cleanedLine {
		if isId(elem) {
			if con.Predecessor != "" {
				return context{}, fmt.Errorf("more then one predecessor in line: %v", line)
			}
			con.Predecessor = elem
			continue
		}
		if r := getRef(elem); r.Bibkey != "" {
			con.References = append(con.References, r)
			continue
		}
		if line != "" {
			con.Context = append(con.Context, elem)
		}
	}

	return con, nil
}

func parseContext2(filename string) (context, error) {
	if filename == "" {
		return context{}, nil
	}
	cleanedLine := clean(strings.Split(filename, ","))

	var con context
	for _, ch := range cleanedLine {
		if r := getRef(ch); r.Bibkey != "" {
			con.References = append(con.References, r)
			continue
		}
		if filename != "" {
			con.Context = append(con.Context, ch)
		}
	}

	return con, nil
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
	References  []zet.Reference
	Context     []string
}
