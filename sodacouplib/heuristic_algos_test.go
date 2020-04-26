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

func TestClaimingPair(t *testing.T) {
	t.Run("horizontal claiming pair should remove from block", func(t *testing.T) {
		s := newEmptySudoku()

		trow := arbitraryRow()
		tcol := arbitraryCol()
		N := arbitraryValue()

		// mark all cells in a row as unavailable for N except inside
		// the test block
		blockR, blockC := trow-(trow%3), tcol-(tcol%3)
		for col := 0; col < 9; col++ {
			outSideTestBlock := col/3 != blockC/3
			if outSideTestBlock {
				s.cells[trow][col].removeCandidate(N)
			}
		}

		for row := blockR; row < blockR+3; row++ {
			if row != trow {
				for col := blockC; col < blockC+3; col++ {
					assert.Equal(t, true, s.cells[row][col].hasCandidate(N))
				}
			}
		}

		apply(t, claimingPair, s)

		for row := blockR; row < blockR+3; row++ {
			if row != trow {
				for col := blockC; col < blockC+3; col++ {
					assert.Equal(t, false, s.cells[row][col].hasCandidate(N))
				}
			}
		}

		// second time should be no-op
		noOpCheck(t, claimingPair, s)
	})
	t.Run("vertical claiming pairs should remove from block", func(t *testing.T) {
		s := newEmptySudoku()

		trow := arbitraryRow()
		tcol := arbitraryCol()
		N := arbitraryValue()

		// mark all cells in a column as unavailable for N except inside
		// the test block
		blockR, blockC := trow-(trow%3), tcol-(tcol%3)
		for row := 0; row < 9; row++ {
			outSideTestBlock := row/3 != blockR/3
			if outSideTestBlock {
				s.cells[row][tcol].removeCandidate(N)
			}
		}

		for row := blockR; row < blockR+3; row++ {
			for col := blockC; col < blockC+3; col++ {
				if col != tcol {
					assert.Equal(t, true, s.cells[row][col].hasCandidate(N))
				}
			}
		}

		apply(t, claimingPair, s)

		for row := blockR; row < blockR+3; row++ {
			for col := blockC; col < blockC+3; col++ {
				if col != tcol {
					assert.Equal(t, false, s.cells[row][col].hasCandidate(N))
				}
			}
		}

		// second time should be no-op
		noOpCheck(t, claimingPair, s)
	})
}

func TestNakedPair(t *testing.T) {
	t.Run("horizontal naked pair", func(t *testing.T) {
		s := newEmptySudoku()

		trow := arbitraryRow()
		tcol1, tcol2 := arbitraryTwoCols()
		N1, N2 := arbitraryTwoValues()

		for n := 1; n <= 9; n++ {
			if n != N1 && n != N2 {
				s.cells[trow][tcol1].removeCandidate(n)
				s.cells[trow][tcol2].removeCandidate(n)
			}
		}

		for col := 0; col < 9; col++ {
			if col != tcol1 && col != tcol2 {
				assert.Equal(t, true, s.cells[trow][col].hasCandidate(N1))
				assert.Equal(t, true, s.cells[trow][col].hasCandidate(N2))
			} else {
				assert.Equal(t, true, s.cells[trow][col].hasCandidate(N1))
				assert.Equal(t, true, s.cells[trow][col].hasCandidate(N2))
			}
		}

		apply(t, nakedPair, s)

		for col := 0; col < 9; col++ {
			if col != tcol1 && col != tcol2 {
				// should remove the pair as candidates, and leave the rest
				assert.Equal(t, false, s.cells[trow][col].hasCandidate(N1))
				assert.Equal(t, false, s.cells[trow][col].hasCandidate(N2))
				for n := 1; n <= 9; n++ {
					if n != N1 && n != N2 {
						assert.Equal(t, true, s.cells[trow][col].hasCandidate(n))
					}
				}
			} else {
				assert.Equal(t, true, s.cells[trow][col].hasCandidate(N1))
				assert.Equal(t, true, s.cells[trow][col].hasCandidate(N2))
			}
		}

		// should be no-op second time
		noOpCheck(t, nakedPair, s)
	})
	t.Run("vertical naked pair", func(t *testing.T) {
		s := newEmptySudoku()

		trow1, trow2 := arbitraryTwoRows()
		tcol := arbitraryCol()
		N1, N2 := arbitraryTwoValues()

		for n := 1; n <= 9; n++ {
			if n != N1 && n != N2 {
				s.cells[trow1][tcol].removeCandidate(n)
				s.cells[trow2][tcol].removeCandidate(n)
			}
		}

		for row := 0; row < 9; row++ {
			if row != trow1 && row != trow2 {
				assert.Equal(t, true, s.cells[row][tcol].hasCandidate(N1))
				assert.Equal(t, true, s.cells[row][tcol].hasCandidate(N2))
			}
		}

		apply(t, nakedPair, s)

		for row := 0; row < 9; row++ {
			if row != trow1 && row != trow2 {
				assert.Equal(t, false, s.cells[row][tcol].hasCandidate(N1))
				assert.Equal(t, false, s.cells[row][tcol].hasCandidate(N2))
				for n := 1; n <= 9; n++ {
					if n != N1 && n != N2 {
						assert.Equal(t, true, s.cells[row][tcol].hasCandidate(n))
					}
				}
			}
		}

		// should be no-op second time
		noOpCheck(t, nakedPair, s)
	})
	t.Run("block naked pair", func(t *testing.T) {
		s := newEmptySudoku()

		trow1, trow2 := arbitraryTwoRowsSameBlock()
		tcol1, tcol2 := arbitraryTwoColsSameBlock()
		N1, N2 := arbitraryTwoValues()

		for n := 1; n <= 9; n++ {
			if n != N1 && n != N2 {
				s.cells[trow1][tcol1].removeCandidate(n)
				s.cells[trow2][tcol2].removeCandidate(n)
			}
		}

		blockR, blockC := trow1-(trow1%3), tcol1-(tcol1%3)
		for row := blockR; row < blockR+3; row++ {
			for col := blockC; col < blockC+3; col++ {
				if (row != trow1 || col != tcol1) && (row != trow2 || col != tcol2) {
					assert.Equal(t, true, s.cells[row][col].hasCandidate(N1))
					assert.Equal(t, true, s.cells[row][col].hasCandidate(N2))
				} else {
					assert.Equal(t, true, s.cells[row][col].hasCandidate(N1))
					assert.Equal(t, true, s.cells[row][col].hasCandidate(N2))
				}
			}
		}

		apply(t, nakedPair, s)

		for row := blockR; row < blockR+3; row++ {
			for col := blockC; col < blockC+3; col++ {
				if (row != trow1 || col != tcol1) && (row != trow2 || col != tcol2) {
					assert.Equal(t, false, s.cells[row][col].hasCandidate(N1))
					assert.Equal(t, false, s.cells[row][col].hasCandidate(N2))
				} else {
					assert.Equal(t, true, s.cells[row][col].hasCandidate(N1))
					assert.Equal(t, true, s.cells[row][col].hasCandidate(N2))
				}
			}
		}

		// should be no-op second time
		noOpCheck(t, nakedPair, s)
	})
}

func TestHiddenPair(t *testing.T) {
	t.Run("row hidden pair", func(t *testing.T) {
		s := newEmptySudoku()

		trow := arbitraryRow()
		tcol1, tcol2 := arbitraryTwoCols()
		N1, N2 := arbitraryTwoValues()

		for col := 0; col < 9; col++ {
			if col != tcol1 && col != tcol2 {
				s.cells[trow][col].removeCandidate(N1)
				s.cells[trow][col].removeCandidate(N2)
			}
		}

		for n := 1; n <= 9; n++ {
			if n != N1 && n != N2 {
				assert.Equal(t, true, s.cells[trow][tcol1].hasCandidate(n))
				assert.Equal(t, true, s.cells[trow][tcol2].hasCandidate(n))
			}
		}

		apply(t, hiddenPair, s)

		for n := 1; n <= 9; n++ {
			if n != N1 && n != N2 {
				assert.Equal(t, false, s.cells[trow][tcol1].hasCandidate(n))
				assert.Equal(t, false, s.cells[trow][tcol2].hasCandidate(n))
			}
		}

		// should be no-op second time
		noOpCheck(t, hiddenPair, s)
	})
	t.Run("column hidden pair", func(t *testing.T) {
		s := newEmptySudoku()

		trow1, trow2 := arbitraryTwoRows()
		tcol := arbitraryCol()
		N1, N2 := arbitraryTwoValues()

		for row := 0; row < 9; row++ {
			if row != trow1 && row != trow2 {
				s.cells[row][tcol].removeCandidate(N1)
				s.cells[row][tcol].removeCandidate(N2)
			}
		}

		for n := 1; n <= 9; n++ {
			if n != N1 && n != N2 {
				assert.Equal(t, true, s.cells[trow1][tcol].hasCandidate(n))
				assert.Equal(t, true, s.cells[trow2][tcol].hasCandidate(n))
			}
		}

		apply(t, hiddenPair, s)

		for n := 1; n <= 9; n++ {
			if n != N1 && n != N2 {
				assert.Equal(t, false, s.cells[trow1][tcol].hasCandidate(n))
				assert.Equal(t, false, s.cells[trow2][tcol].hasCandidate(n))
			}
		}

		// should be no-op second time
		noOpCheck(t, hiddenPair, s)
	})
	t.Run("block hidden pair", func(t *testing.T) {
		s := newEmptySudoku()

		// given two candidates and two cells inside a block
		trow1, trow2 := arbitraryTwoRowsSameBlock()
		tcol1, tcol2 := arbitraryTwoColsSameBlock()
		N1, N2 := arbitraryTwoValues()

		// and all the other cells do not have have those two candidates
		blockR, blockC := trow1-(trow1%3), tcol1-(tcol1%3)
		for row := blockR; row < blockR+3; row++ {
			for col := blockC; col < blockC+3; col++ {
				if (row != trow1 || col != tcol1) && (row != trow2 || col != tcol2) {
					s.cells[row][col].removeCandidate(N1)
					s.cells[row][col].removeCandidate(N2)
				}
			}
		}

		for n := 1; n <= 9; n++ {
			if n != N1 && n != N2 {
				assert.Equal(t, true, s.cells[trow1][tcol1].hasCandidate(n))
				assert.Equal(t, true, s.cells[trow2][tcol2].hasCandidate(n))
			}
		}

		// then applying hidden pair
		apply(t, hiddenPair, s)

		// should remove all other candidates from the two cells
		for n := 1; n <= 9; n++ {
			if n != N1 && n != N2 {
				assert.Equal(t, false, s.cells[trow1][tcol1].hasCandidate(n))
				assert.Equal(t, false, s.cells[trow2][tcol2].hasCandidate(n))
			}
		}

		// should be no-op second time
		noOpCheck(t, hiddenPair, s)
	})
}

func TestXWing(t *testing.T) {
	t.Run("row based xwing", func(t *testing.T) {
		s := newEmptySudoku()

		// given two rows that have two parallel remaining cells slots
		// for a candidate
		trow1, trow2 := arbitraryTwoRows()
		tcol1, tcol2 := arbitraryTwoCols()
		N := arbitraryValue()

		for col := 0; col < 9; col++ {
			if col != tcol1 && col != tcol2 {
				s.cells[trow1][col].removeCandidate(N)
				s.cells[trow2][col].removeCandidate(N)
			}
		}

		for row := 0; row < 9; row++ {
			if row != trow1 && row != trow2 {
				assert.Equal(t, true, s.cells[row][tcol1].hasCandidate(N))
				assert.Equal(t, true, s.cells[row][tcol2].hasCandidate(N))
			}
		}

		// then xwing algorithm
		apply(t, xWing, s)

		// should remove that candidate from same columns on other rows
		for row := 0; row < 9; row++ {
			if row != trow1 && row != trow2 {
				assert.Equal(t, false, s.cells[row][tcol1].hasCandidate(N))
				assert.Equal(t, false, s.cells[row][tcol2].hasCandidate(N))
			}
		}

		// should be no-op second time
		noOpCheck(t, xWing, s)
	})
	t.Run("col based xwing", func(t *testing.T) {
		s := newEmptySudoku()

		// given two columns that have two parallel remaining cells slots
		// for a candidate
		trow1, trow2 := arbitraryTwoRows()
		tcol1, tcol2 := arbitraryTwoCols()
		N := arbitraryValue()

		for row := 0; row < 9; row++ {
			if row != trow1 && row != trow2 {
				s.cells[row][tcol1].removeCandidate(N)
				s.cells[row][tcol2].removeCandidate(N)
			}
		}

		for col := 0; col < 9; col++ {
			if col != tcol1 && col != tcol2 {
				assert.Equal(t, true, s.cells[trow1][col].hasCandidate(N))
				assert.Equal(t, true, s.cells[trow2][col].hasCandidate(N))
			}
		}

		// then xwing algorithm
		apply(t, xWing, s)

		// should remove that candidate from same rows on other columns
		for col := 0; col < 9; col++ {
			if col != tcol1 && col != tcol2 {
				assert.Equal(t, false, s.cells[trow1][col].hasCandidate(N))
				assert.Equal(t, false, s.cells[trow2][col].hasCandidate(N))
			}
		}

		// and should be no-op second time
		noOpCheck(t, xWing, s)
	})
	t.Run("sample xwing problem", func(t *testing.T) {
		problem := `
			___ __5 __9
			_8_ 4__ ___
			___ _7_ __4
			
			_32 748 9__
			__6 539 _42
			___ 216 37_
			
			__4 _57 293
			__5 _21 4_7
			2__ __4 ___
		`
		expectedStep := `
			___ _85 __9
			_8_ 4__ ___
			___ _7_ __4
			
			_32 748 9__
			__6 539 _42
			___ 216 37_
			
			__4 _57 293
			__5 _21 4_7
			2__ __4 ___
		`
		s, err := NewSudokuSquare(problem)
		if err != nil {
			t.Fatal("got unexpected error from valid input:", err)
		}

		apply(t, xWing, s)
		apply(t, hiddenSingle, s)

		expected, _ := FormatSudoku(expectedStep)
		actual, _ := FormatSudoku(s.String())
		assert.Equal(t, expected, actual)
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

func arbitraryTwoCols() (int, int) {
	c1, c2 := rand2N()
	log.Println("using cols:", c1, c2)
	return c1, c2
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
