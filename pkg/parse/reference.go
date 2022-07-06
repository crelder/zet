package parse

import (
	"regexp"
)

// Reference parses the string and returns a slice of bibkeys.
// It uses a two-step process with regex to filter out the bibkeys.
func Reference(s string) []string {
	re1 := regexp.MustCompile(`[@][a-z]+[{][a-z]+\d{4}[a-z]?[,]`)
	litKeys := re1.FindAll([]byte(s), -1)

	re2 := regexp.MustCompile(`[a-z]+\d{4}[a-z]?`)
	var result []string
	for _, litKey := range litKeys {
		result = append(result, string(re2.Find(litKey)))
	}

	return result
}
