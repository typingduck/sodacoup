package sodacouplib

import (
	"github.com/stretchr/testify/assert"
	"log"
	"math/rand"
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

func TestHiddenSingle(t *testing.T) {
	t.Run("hidden single found in row", func(t *testing.T) {
		s := newEmptySudoku()

		trow := arbitraryRow()
		tcol := arbitraryCol()
		N := arbitraryValue()

		// mark all cells in a row as available for N except one
		for col := 0; col < 9; col++ {
			if col != tcol {
				s.cells[trow][col].removeCandidate(N)
			}
		}

		assert.Equal(t, false, s.cells[trow][neighbour(tcol-1)].hasCandidate(N))
		assert.Equal(t, true, s.cells[trow][tcol].hasCandidate(N))
		assert.Equal(t, false, s.cells[trow][neighbour(tcol+1)].hasCandidate(N))

		apply(t, hiddenSingle, s)

		cell := s.cells[trow][tcol]
		assert.Equal(t, true, cell.isSet)
		assert.Equal(t, uint8(N), cell.value)

		// second time should be no-op
		noOpCheck(t, hiddenSingle, s)
	})
	t.Run("hidden single found in column", func(t *testing.T) {
		s := newEmptySudoku()

		trow := arbitraryRow()
		tcol := arbitraryCol()
		N := arbitraryValue()

		// mark all cells in a column as available for N except one
		for row := 0; row < 9; row++ {
			if row != trow {
				s.cells[row][tcol].removeCandidate(N)
			}
		}

		assert.Equal(t, false, s.cells[neighbour(trow-1)][tcol].hasCandidate(N))
		assert.Equal(t, true, s.cells[trow][tcol].hasCandidate(N))
		assert.Equal(t, false, s.cells[neighbour(trow+1)][tcol].hasCandidate(N))

		apply(t, hiddenSingle, s)

		cell := s.cells[trow][tcol]
		assert.Equal(t, true, cell.isSet)
		assert.Equal(t, uint8(N), cell.value)

		// second time should be no-op
		noOpCheck(t, hiddenSingle, s)
	})
	t.Run("hidden single found in block", func(t *testing.T) {
		s := newEmptySudoku()

		trow := arbitraryRow()
		tcol := arbitraryCol()
		N := arbitraryValue()

		// given all cells in test block as available for N except one
		blockR, blockC := trow-(trow%3), tcol-(tcol%3)
		for row := blockR; row < blockR+3; row++ {
			for col := blockC; col < blockC+3; col++ {
				if row != trow || col != tcol {
					s.cells[row][col].removeCandidate(N)
				}
			}
		}

		// then hiddenSingle
		apply(t, hiddenSingle, s)

		// should fill in that one cell
		cell := s.cells[trow][tcol]
		assert.Equal(t, true, cell.isSet)
		assert.Equal(t, uint8(N), cell.value)

		// second time should be no-op
		noOpCheck(t, hiddenSingle, s)
	})
	t.Run("should fill multiple singles", func(t *testing.T) {
		s := newEmptySudoku()

		trow1, trow2 := arbitraryTwoRows()
		tcol := arbitraryCol()
		N1, N2 := arbitraryTwoValues()

		// given two cells which are both unique to a given number
		for row := 0; row < 9; row++ {
			if row != trow1 {
				s.cells[row][tcol].removeCandidate(N1)
			}
			if row != trow2 {
				s.cells[row][tcol].removeCandidate(N2)
			}
		}

		apply(t, hiddenSingle, s)

		// then they both should be filled
		cell := s.cells[trow1][tcol]
		assert.Equal(t, true, cell.isSet)
		assert.Equal(t, uint8(N1), cell.value)

		cell = s.cells[trow2][tcol]
		assert.Equal(t, true, cell.isSet)
		assert.Equal(t, uint8(N2), cell.value)

		// second time should be no-op
		noOpCheck(t, hiddenSingle, s)
	})
}

func TestPointingPair(t *testing.T) {
	t.Run("horizontal pointing pairs should remove from row", func(t *testing.T) {
		s := newEmptySudoku()

		trow := arbitraryRow()
		tcol1, tcol2 := arbitraryTwoColsSameBlock()
		N := arbitraryValue()

		// when all cells in a block as unavailable for N except two
		// in same block
		blockR, blockC := trow-(trow%3), tcol1-(tcol1%3)
		for row := blockR; row < blockR+3; row++ {
			for col := blockC; col < blockC+3; col++ {
				if row != trow || (col != tcol1 && col != tcol2) {
					s.cells[row][col].removeCandidate(N)
				}
			}
		}

		for col := 0; col < 9; col++ {
			outSideTestBlock := col/3 != blockC/3
			if outSideTestBlock {
				assert.Equal(t, true, s.cells[trow][col].hasCandidate(N))
			}
		}

		// then pointing pair
		apply(t, pointingPair, s)

		// should remove from rest of row
		for col := 0; col < 9; col++ {
			outSideTestBlock := col/3 != blockC/3
			if outSideTestBlock {
				assert.Equal(t, false, s.cells[trow][col].hasCandidate(N))
			}
		}

		// second time should be no-op
		noOpCheck(t, pointingPair, s)
	})
	t.Run("vertical pointing pairs should remove from column", func(t *testing.T) {
		s := newEmptySudoku()

		trow1, trow2 := arbitraryTwoRowsSameBlock()
		tcol := arbitraryCol()
		N := arbitraryValue()

		// mark all cells in a block as unavailable for N except two
		// in same block on same column
		blockR, blockC := trow1-(trow1%3), tcol-(tcol%3)
		for row := blockR; row < blockR+3; row++ {
			for col := blockC; col < blockC+3; col++ {
				if col != tcol || (row != trow1 && row != trow2) {
					s.cells[row][col].removeCandidate(N)
				}
			}
		}

		for row := 0; row < 9; row++ {
			outSideTestBlock := row/3 != blockR/3
			if outSideTestBlock {
				assert.Equal(t, true, s.cells[row][tcol].hasCandidate(N))
			}
		}

		apply(t, pointingPair, s)

		for row := 0; row < 9; row++ {
			outSideTestBlock := row/3 != blockR/3
			if outSideTestBlock {
				assert.Equal(t, false, s.cells[row][tcol].hasCandidate(N))
			}
		}

		// second time should be no-op
		noOpCheck(t, pointingPair, s)
	})
}

func arbitraryRow() int {
	r := rand.Intn(9)
	log.Println("using row:", r)
	return r
}

func arbitraryCol() int {
	c := rand.Intn(9)
	log.Println("using col:", c)
	return c
}

func arbitraryValue() int {
	n := 1 + rand.Intn(9)
	log.Println("using value:", n)
	return n
}

func rand2N() (int, int) {
	a := rand.Intn(9)
	for {
		b := rand.Intn(9)
		if a != b {
			return a, b
		}
	}
}

func rand2NBlock() (int, int) {
	for {
		a := rand.Intn(9)
		b := rand.Intn(9)
		if a != b && a/3 == b/3 {
			return a, b
		}
	}
}

func arbitraryTwoRows() (int, int) {
	r1, r2 := rand2N()
	log.Println("using rows:", r1, r2)
	return r1, r2
}

func arbitraryTwoRowsSameBlock() (int, int) {
	r1, r2 := rand2NBlock()
	log.Println("using rows:", r1, r2)
	return r1, r2
}

func arbitraryTwoColsSameBlock() (int, int) {
	c1, c2 := rand2NBlock()
	log.Println("using cols:", c1, c2)
	return c1, c2
}

func arbitraryTwoValues() (int, int) {
	a, b := rand2N()
	v1, v2 := a+1, b+1
	log.Println("using values:", v1, v2)
	return v1, v2
}

func neighbour(n int) int {
	return (n + 9) % 9
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
