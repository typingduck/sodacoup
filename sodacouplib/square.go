package sodacouplib

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
)

// SudokuSquare does the magic
type SudokuSquare struct {
	cells [9][9]byte
}

// NewSudokuSquare Create a SudokuSquare struct given a string that
// roughly looks like a sudoku problem. Can have many spaces/newlines, but
// just needs '_' for empty cells or a number 1-9 for filled cells.
func NewSudokuSquare(stringRepresentation string) (*SudokuSquare, error) {
	stringRepresentation = filterValidChars(stringRepresentation)
	if len(stringRepresentation) != 81 {
		return nil, errors.New("doesn't look like a valid sudoku")
	}
	var sud SudokuSquare

	si := 0
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			r := stringRepresentation[si]
			if r != '_' {
				if e := sud.setCell(i, j, int(r-'0')); e != nil {
					return nil, e
				}
			}
			si++
		}
	}
	return &sud, nil
}

// Solve TODO.....
func (sud *SudokuSquare) Solve() error {
	return nil
}

func (sud SudokuSquare) String() string {
	return sud.asTableString()
}

func filterValidChars(s string) string {
	return strings.Map(func(r rune) rune {
		if r == '_' || unicode.IsNumber(r) {
			return r
		}
		return -1
	}, s)
}

// FormatSudoku takes a sudoku like string and prints it in the square format
// with spaces/empty line between blocks.
func FormatSudoku(s string) (string, error) {
	s = filterValidChars(s)
	if len(s) != 81 {
		return "", errors.New("invalid length")
	}

	var o []byte = make([]byte, 81+9+2+9*2) // numbers + newlines + extra newlines + spaces between squares
	t := 0
	for i := 0; i < 81; i++ {
		o[t] = s[i]
		t++
		if i%27 == 26 && i != 80 {
			o[t] = '\n'
			t++
		}
		if i%9 == 8 {
			o[t] = '\n'
			t++
		} else if i%3 == 2 {
			o[t] = ' '
			t++
		}
	}
	return string(o), nil
}

// Format a sudoku as a table with lines between blocks.
func (sud SudokuSquare) asTableString() string {
	var sb strings.Builder
	const hr = " -------------------------\n"
	sb.WriteString(hr)
	for r := 0; r < 9; r++ {
		sb.WriteString(" |")
		for c := 0; c < 9; c++ {
			sb.WriteByte(' ')
			cell := sud.cells[r][c]
			if cell != 0 {
				fmt.Fprintf(&sb, "%d", cell)
			} else {
				sb.WriteByte('_')
			}
			if c%3 == 2 {
				sb.WriteString(" |")
			}
		}
		sb.WriteByte('\n')
		if r%3 == 2 {
			sb.WriteString(hr)
		}
	}
	return sb.String()
}

func (sud *SudokuSquare) setCell(row int, col int, val int) error {
	if row < 0 || row >= 9 || col < 0 || col >= 9 {
		return errors.New("program logic wrong")
	}
	if val < 1 || val > 9 {
		return fmt.Errorf("invalid input value: %d", val)
	}
	sud.cells[row][col] = byte(val)
	return nil
}
