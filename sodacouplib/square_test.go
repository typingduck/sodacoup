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

func TestLoadingInitialProblem(t *testing.T) {
	initial := `
		1_3 456 789
		456 789 123
		789 123 456

		234 567 891
		567 891 234
		891 234 567

		345 678 912
		678 912 345
		912 345 678
	`

	s, err := NewSudokuSquare(initial)

	if err != nil {
		t.Fatalf("got unexpected error from valid input: %s", err)
	}

	c1 := s.cells[0][0]
	assert.Equal(t, true, c1.isSet)
	assert.Equal(t, uint8(1), c1.value)
	assert.Equal(t, false, c1.hasCandidate(1))
	assert.Equal(t, false, c1.hasCandidate(2))
	assert.Equal(t, false, c1.hasCandidate(3))

	c2 := s.cells[0][1]
	assert.Equal(t, false, c2.isSet)

	assert.Equal(t, false, c2.hasCandidate(1))
	assert.Equal(t, true, c2.hasCandidate(2))
	assert.Equal(t, false, c2.hasCandidate(3))
}
