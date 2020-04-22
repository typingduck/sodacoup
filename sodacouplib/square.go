package sodacouplib

import (
	"errors"
	"strings"
	"unicode"
)

// SudokuSquare does the magic
type SudokuSquare struct {
}

// NewSudokuSquare Create a SudokuSquare struct given a string that
// roughly looks like a sudoku problem. Can have many spaces/newlines, but
// just needs '_' for empty cells or a number 1-9 for filled cells.
func NewSudokuSquare(stringRepresentation string) (*SudokuSquare, error) {
	stringRepresentation = filterValidChars(stringRepresentation)
	if len(stringRepresentation) != 81 {
		return nil, errors.New("doesn't look like a valid sudoku")
	}
	return &SudokuSquare{}, nil
}

// Solve TODO.....
func (sud *SudokuSquare) Solve() error {
	return nil
}

func (sud SudokuSquare) String() string {
	return "todo......"
}

func filterValidChars(s string) string {
	return strings.Map(func(r rune) rune {
		if r == '_' || unicode.IsNumber(r) {
			return r
		}
		return -1
	}, s)
}
