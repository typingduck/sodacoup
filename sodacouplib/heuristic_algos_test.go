package sodacouplib

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNakedSingle(t *testing.T) {
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
		t.Fatal("got unexpected error from valid input:", err)
	}

	cell := &s.cells[0][1]
	assert.Equal(t, false, cell.isSet)
	assert.Equal(t, true, cell.hasCandidate(2))
	for n := 1; n <= 9; n++ {
		if n != 2 {
			assert.Equal(t, false, cell.hasCandidate(n))
		}
	}

	apply(t, nakedSingle, s)

	assert.Equal(t, true, cell.isSet)
	assert.Equal(t, uint8(2), cell.value)

	// second time should be no-op
	noOpCheck(t, nakedSingle, s)
}

func apply(t *testing.T, fn func(s *SudokuSquare) (bool, error), s *SudokuSquare) {
	changes, err := fn(s)
	if err != nil {
		t.Fatal("Got unexpected algorithm error", err)
	}
	if !changes {
		t.Fatal("algorithm did not make changes as expected")
	}

}

func noOpCheck(t *testing.T, fn func(s *SudokuSquare) (bool, error), s *SudokuSquare) {
	changes, err := fn(s)
	if err != nil {
		t.Fatal("Got unexpected algorithm error", err)
	}
	if changes {
		t.Fatal("algorithm made unexpected changes")
	}
}
