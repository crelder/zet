package parse

import (
	"errors"
	"fmt"
	"github.com/crelder/zet"
	"strings"
)

// Index parses the content of an index.
// It returns all parsing errors that occurred while parsing each line.
func Index(content string) (zet.Index, []error) {
	var parsErrs []error
	if content == "" {
		parsErrs = append(parsErrs, errors.New("parse Index: cannot parse empty string"))
		return nil, parsErrs
	}
	result := make(map[string][]string)

	lines := strings.Split(content, "\n")
out:
	for lineNumber, line := range lines {
		if line == "" {
			continue
		}
		index := strings.Split(line, ":")
		if len(index) != 2 {
			parsErrs = append(parsErrs, fmt.Errorf("index: could not parse line %q", line))
			continue
		}
		topic := strings.TrimSpace(index[0])
		ids := strings.Split(index[1], ",")

		if len(ids) == 1 && strings.TrimSpace(ids[0]) == "" { // TODO: make so that all potential positions get cleaned
			parsErrs = append(parsErrs, fmt.Errorf("index: could not parse line %q, no ids provided", line))
			continue
		}

		for _, id := range ids {
			if parseId(strings.TrimSpace(id)) == "" {
				parsErrs = append(
					parsErrs,
					fmt.Errorf("index: could not parse line %v, not an id %q", lineNumber, id))
				continue out
			}
		}

		var trimmedIds []string
		for _, id := range ids {
			trimmedIds = append(trimmedIds, strings.TrimSpace(id))
		}

		result[topic] = trimmedIds
	}

	if len(result) == 0 {
		return nil, parsErrs
	}

	return result, parsErrs
}
