package sodacouplib

// Copypasta: code that the users of stackoverflow (probably) strongly believe
// does the thing that's required.
// Utility functions that are probably part or a library but not worth being
// a dependency on for now.

import (
	"strings"
	"unicode"
)

// remove all spaces/newlines eetc
func removeWhitespace(str string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, str)
}
