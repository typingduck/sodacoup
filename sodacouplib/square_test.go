package sodacouplib

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

const basicFormat = `1_3 _56 789
_56 789 1_3
789 1_3 _56

_3_ 567 891
567 891 _3_
891 _3_ 567

3_5 678 91_
678 91_ 3_5
91_ 3_5 678
`

const tableFormat = ` -------------------------
 | 1 _ 3 | _ 5 6 | 7 8 9 |
 | _ 5 6 | 7 8 9 | 1 _ 3 |
 | 7 8 9 | 1 _ 3 | _ 5 6 |
 -------------------------
 | _ 3 _ | 5 6 7 | 8 9 1 |
 | 5 6 7 | 8 9 1 | _ 3 _ |
 | 8 9 1 | _ 3 _ | 5 6 7 |
 -------------------------
 | 3 _ 5 | 6 7 8 | 9 1 _ |
 | 6 7 8 | 9 1 _ | 3 _ 5 |
 | 9 1 _ | 3 _ 5 | 6 7 8 |
 -------------------------
`

func TestFormatSudoku(t *testing.T) {
	t.Run("should format", func(t *testing.T) {

		input := removeWhitespace(basicFormat)

		result, err := FormatSudoku(input)
		if err != nil {
			t.Fatal("error formatting input:", err)
		}

		assert.Equal(t, basicFormat, result)
	})
	t.Run("bad input should give error", func(t *testing.T) {
		badInput := "a crossword!"
		if _, err := FormatSudoku(badInput); err == nil {
			t.Fatal("got unexpected error formatting string ", err)
		}
		if _, err := FormatSudoku(badInput); err == nil {
			t.Fatal("expected error from bad input but none given")
		}
	})
}

func TestNewSudokuSquare(t *testing.T) {
	t.Run("invalid string", func(t *testing.T) {

		invalidString := "nonsense"
		if _, err := NewSudokuSquare(invalidString); err == nil {
			t.Errorf("expected error from bad input but none given")
		}
	})
	t.Run("bad numbers", func(t *testing.T) {
		invalidNumbers := strings.ReplaceAll(basicFormat, "1", "0")
		if _, err := NewSudokuSquare(invalidNumbers); err == nil {
			t.Errorf("expected error from bad input but none given")
		}

	})
	t.Run("not valid sudoku", func(t *testing.T) {
		t.Skip() // TODO...
		invalidPuzzle := `
			555 555 555
			555 555 555
			555 555 555

			555 555 555
			555 555 555
			555 555 555

			555 555 555
			555 555 555
			555 555 555
		`
		if _, err := NewSudokuSquare(invalidPuzzle); err == nil {
			t.Errorf("expected error from bad input but none given")
		}
	})
	t.Run("valid string", func(t *testing.T) {
		if s, err := NewSudokuSquare(basicFormat); err != nil {
			t.Errorf("got unexpected error from creating table %s", err)
		} else {
			result, _ := FormatSudoku(s.String())
			expected := basicFormat

			assert.Equal(t, expected, result)
		}
	})
}

func TestToString(t *testing.T) {
	if s, err := NewSudokuSquare(basicFormat); err != nil {
		t.Errorf("got unexpected error creating table %s", err)
	} else {
		result := s.String()
		expected := tableFormat

		assert.Equal(t, expected, result)
	}
}
