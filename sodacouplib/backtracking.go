package sodacouplib

import "errors"

// backTrack Finds a solution SudokuSquare using a backtrackling algorithm.
// (Doesn't check the resulting solution is unique).
func backTrack(sud *SudokuSquare) (bool, error) {
	if backTrackRecursive(sud, 0, 0) {
		return false, nil
	}
	return false, errors.New("failed to converge")
}

func backTrackRecursive(sud *SudokuSquare, row, col int) bool {
	if col == 9 {
		col = 0
		row++
	}
	if row == 9 {
		return true
	}
	cell := sud.cells[row][col]
	isSet := cell != 0
	if isSet {
		return backTrackRecursive(sud, row, col+1)
	}
	for n := 1; n <= 9; n++ {
		if isValidMove(sud, row, col, n) {
			sud.cells[row][col] = byte(n)
			success := backTrackRecursive(sud, row, col+1)
			if success {
				return true
			}
			sud.cells[row][col] = 0
		}
	}
	return false
}

func isValidMove(sud *SudokuSquare, row, col, val int) bool {
	n := byte(val)
	cell := sud.cells[row][col]
	if cell != 0 {
		return false
	}
	for r := 0; r < 9; r++ {
		if sud.cells[r][col] == n {
			return false
		}
	}
	for c := 0; c < 9; c++ {
		if sud.cells[row][c] == n {
			return false
		}
	}
	blkRow := row - row%3
	blkCol := col - col%3
	for r := blkRow; r < blkRow+3; r++ {
		for c := blkCol; c < blkCol+3; c++ {
			if sud.cells[r][c] == n {
				return false
			}
		}
	}
	return true
}
