package parse

import "github.com/crelder/zet"

// Parser is a wrapper for the exported functions in this package.
// Parser satisfies the zet.Parser interface.
type Parser struct{}

func New() Parser {
	return Parser{}
}

func (p Parser) Content(filename string, zettel []zet.Zettel) (string, error) {
	return Content(filename, zettel)
}

func (p Parser) Filename(s string) (zet.Zettel, error) {
	return Filename(s)
}
func (p Parser) Index(content string) (zet.Index, []zet.InconErr) {
	return Index(content)
}

func (p Parser) Reference(d string) []string {
	return Reference(d)
}
