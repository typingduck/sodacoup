package sodacouplib

import "errors"

// backTrack Finds a solution SudokuSquare using a backtrackling algorithm.
// (Doesn't check the resulting solution is unique).
func backTrack(sud *SudokuSquare) (bool, error) {
	var cells [9][9]byte
	copyFrom(*sud, &cells)
	if backTrackRecursive(&cells, 0, 0) {
		copyTo(cells, sud)
		return false, nil
	}
	return false, errors.New("failed to converge")
}

func backTrackRecursive(cells *[9][9]byte, row, col int) bool {
	if col == 9 {
		col = 0
		row++
	}
	if row == 9 {
		return true
	}
	cell := cells[row][col]
	if isSet(cell) {
		return backTrackRecursive(cells, row, col+1)
	}
	for n := 1; n <= 9; n++ {
		if isValidMove(cells, row, col, n) {
			cells[row][col] = byte(n)
			success := backTrackRecursive(cells, row, col+1)
			if success {
				return true
			}
			cells[row][col] = 0
		}
	}
	return false
}

func isValidMove(cells *[9][9]byte, row, col, val int) bool {
	n := byte(val)
	if isSet(cells[row][col]) {
		return false
	}
	for r := 0; r < 9; r++ {
		if cells[r][col] == n {
			return false
		}
	}
	for c := 0; c < 9; c++ {
		if cells[row][c] == n {
			return false
		}
	}
	blkRow := row - row%3
	blkCol := col - col%3
	for r := blkRow; r < blkRow+3; r++ {
		for c := blkCol; c < blkCol+3; c++ {
			if cells[r][c] == n {
				return false
			}
		}
	}
	return true
}

func isSet(cell byte) bool {
	return cell != 0
}

func copyFrom(sud SudokuSquare, cells *[9][9]byte) {
	for r := 0; r < 9; r++ {
		for c := 0; c < 9; c++ {
			if sud.cells[r][c].isSet {
				cells[r][c] = sud.cells[r][c].value
			}
		}
	}
}

func copyTo(cells [9][9]byte, sud *SudokuSquare) {
	for r := 0; r < 9; r++ {
		for c := 0; c < 9; c++ {
			if isSet(cells[r][c]) && !sud.cells[r][c].isSet {
				e := sud.setCell(r, c, int(cells[r][c]))
				if e != nil {
					panic(e) // something completely wrong if we end up here
				}
			}
		}
	}
}
