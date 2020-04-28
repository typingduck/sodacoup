package sodacouplib

import (
	"math/rand"
)

// Generates a problem that is solvable by the heuristic algorithms.
func GenerateProblem() (*SudokuSquare, error) {
	sud := randomFilledSudoku()
	// removing is more efficient than adding because of the way backtracking
	// works.
	err := removeCellsWhileSolvable(sud)
	return sud, err
}

func randomFilledSudoku() *SudokuSquare {
	s := newEmptySudoku()
	for row := 0; row < 9; row++ {
		for col := 0; col < 9; col++ {
			if e := s.setCell(row, col, findValueThatFitsCell(s, row, col)); e != nil {
				panic(e)
			} else if _, e := sanityCheck(s); e != nil {
				panic(e)
			}
		}
	}
	return s
}

func findValueThatFitsCell(s *SudokuSquare, row, col int) int {
	for {
		val := rand.Intn(9) + 1
		cell := &s.cells[row][col]
		if cell.hasCandidate(val) {
			tmp := *s
			if e := tmp.setCell(row, col, val); e != nil {
				panic(e)
			}
			_, e := backTrack(&tmp)
			isSolvable := e == nil
			if isSolvable {
				return val
			}
		}
	}
}

func removeCellsWhileSolvable(sud *SudokuSquare) error {
	var cells [9][9]byte
	copyFrom(*sud, &cells)
	err := redoWhileMakingChanges(func() (bool, error) {
		row, col := rand.Intn(9), rand.Intn(9)
		if !isSet(cells[row][col]) {
			return false, nil
		}
		removable, err := canRemove(cells, row, col)
		if err != nil {
			return false, err
		}
		if removable {
			cells[row][col] = 0
			return true, nil
		}
		return false, nil
	})
	if err != nil {
		return err
	}
	*sud = *newEmptySudoku()
	copyTo(cells, sud)
	return nil
}

func canRemove(cells [9][9]byte, row, col int) (bool, error) {
	cells[row][col] = 0
	sud := newEmptySudoku()
	copyTo(cells, sud)
	return trySolveWithHeuristics(sud)
}

// Keep doing `fn` as long as it's try and give up after
// a certain number of retries if it's false
func redoWhileMakingChanges(fn func() (bool, error)) error {
	retries := 100
	for i := 0; i < retries; i++ {
		b, e := fn()
		if e != nil {
			return e
		}
		if b {
			i = 0
		}
	}
	return nil
}
