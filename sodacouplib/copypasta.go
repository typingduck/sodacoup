package sodacouplib

// Copypasta: code that the users of stackoverflow (probably) strongly believe
// does the thing that's required.
// Utility functions that are probably part of a library but not worth being
// a dependency on for now.

import (
	"log"
	"math/rand"
	"reflect"
	"runtime"
	"strings"
	"time"
	"unicode"
)

// remove all spaces/newlines eetc
func removeWhitespace(s string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, s)
}

func getFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

func RandomSeed() {
	seed := time.Now().UTC().UnixNano()
	log.Println("using seed", seed)
	rand.Seed(seed)
}
