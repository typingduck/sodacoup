package sodacouplib

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"unicode"
)

// SudokuSquare Wraps the square in some useful constructs.
type SudokuSquare struct {
	cells [9][9]SudokuCell
}

// SudokuCell adds a little info to each cell to make heuristic algorithms easier.
// Keeps a track of what are valid candidates for the cell in the candidates bitmask.
// To keep that consistent all updates have to be done through SudokuSquare.setCell.
type SudokuCell struct {
	row, col   int
	value      byte // 0 -> 9 inclusive (0 for unset)
	isSet      bool
	candidates uint16 // bitmask 2^1 -> 2^9 of still valid cell numbers
}

// NewSudokuSquare Create a SudokuSquare struct given a string that
// roughly looks like a sudoku problem. Can have many spaces/newlines, but
// just needs '_' for empty cells or a number 1-9 for filled cells.
func NewSudokuSquare(stringRepresentation string) (*SudokuSquare, error) {
	stringRepresentation = filterValidChars(stringRepresentation)
	if len(stringRepresentation) != 81 {
		return nil, errors.New("doesn't look like a valid sudoku")
	}

	sud := newEmptySudoku()

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
	return sud, nil
}

// An empty sudoku is probably an oxymoron, but it's useful for testing.
func newEmptySudoku() *SudokuSquare {
	var sud SudokuSquare
	for r := 0; r < 9; r++ {
		for c := 0; c < 9; c++ {
			sud.cells[r][c].init(r, c)
		}
	}
	return &sud
}

// Solve does the magic.
func (sud *SudokuSquare) Solve() error {
	heuristicAlgorithms := []sudokuAlgo{
		nakedSingle,
	}

	e := untilTrue(func() (bool, error) {
		log.Println(sud)
		changesMade := false
		for _, fn := range heuristicAlgorithms {
			impacting, err := fn(sud)
			if err != nil {
				return false, err
			}
			if impacting {
				log.Println("...done applying:", getFunctionName(fn))
			}
			changesMade = changesMade || impacting
		}
		return !changesMade, nil
	})
	if e != nil {
		return e
	}

	if !isSolved(sud) {
		log.Println("Unsolved by heuristics. Applying backtracking.")
		_, e := backTrack(sud)
		return e
	}
	return nil
}

func (sud SudokuSquare) String() string {
	return sud.asTableString()
}

func (c SudokuCell) String() string {
	if c.isSet {
		return fmt.Sprintf("[%d,%d => %d]", c.row, c.col, c.value)
	}
	return fmt.Sprintf("[%d,%d ? %d]", c.row, c.col, c.candidates)
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
			if cell.isSet {
				fmt.Fprintf(&sb, "%d", cell.value)
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

func (c *SudokuCell) removeCandidate(val int) {
	if val < 1 || val > 9 {
		panic("invalid")
	}
	var msk uint16 = 1 << val
	c.candidates = c.candidates &^ msk
}

func (c SudokuCell) hasCandidate(val int) bool {
	if val < 1 || val > 9 {
		panic("invalid")
	}
	var msk uint16 = 1 << val
	return !c.isSet && (c.candidates&msk == msk)
}

func (c *SudokuCell) init(row, col int) {
	c.isSet = false
	c.value = 0
	c.candidates = 0b1111111110
	c.row = row
	c.col = col
}

func applyToCells(sud *SudokuSquare, fn func(cell *SudokuCell) (bool, error)) (bool, error) {
	result := false

	for row := 0; row < 9; row++ {
		for col := 0; col < 9; col++ {
			b, e := fn(&sud.cells[row][col])
			if e != nil {
				return false, e
			}
			result = result || b
		}
	}

	return result, nil
}

func (sud *SudokuSquare) setCell(row int, col int, val int) error {
	if row < 0 || row >= 9 || col < 0 || col >= 9 {
		return errors.New("invalid")
	}
	if val < 1 || val > 9 {
		return errors.New("invalid")
	}
	c := &sud.cells[row][col]
	if c.isSet {
		return errors.New("trying to update already set cell")
	}
	if !c.hasCandidate(val) {
		return fmt.Errorf("trying to add %d to %d,%d", val, row, col)
	}
	c.isSet = true
	c.value = byte(val)

	/* update rows */
	for i := 0; i < 9; i++ {
		sud.cells[i][col].removeCandidate(val)
	}

	/* update cols */
	for j := 0; j < 9; j++ {
		sud.cells[row][j].removeCandidate(val)
	}

	/* update squares */
	si := row - row%3
	sj := col - col%3
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			sud.cells[si+i][sj+j].removeCandidate(val)
		}
	}
	return nil
}

func isSolved(sud *SudokuSquare) bool {
	for r := 0; r < 9; r++ {
		for c := 0; c < 9; c++ {
			if !sud.cells[r][c].isSet {
				return false
			}
		}
	}
	return true
}

// Run until true, but wth the safety of failing after some large amount of iterations.
func untilTrue(fn func() (bool, error)) error {
	const maxIter = 10000
	for i := 0; i <= maxIter; i++ {
		done, err := fn()
		if err != nil {
			return err
		}
		if done {
			return nil
		}
	}
	return errors.New("infinite loooooooooooop detected. somebody screwed up. probably the same person who wrote the word loooooooooooop")
}
