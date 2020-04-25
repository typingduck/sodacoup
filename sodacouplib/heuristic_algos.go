package sodacouplib

import (
	"log"
	"math"
)

type sudokuAlgo func(*SudokuSquare) (bool, error)

func nakedSingle(sud *SudokuSquare) (bool, error) {
	return applyToCells(sud, func(cell *SudokuCell) (bool, error) {
		changes := false

		// candidates is a mask between 2^1 -> 2^9
		hasOneAvailableValue := (cell.candidates & (cell.candidates - 1)) == 0
		if !cell.isSet && hasOneAvailableValue {
			value := int(math.Log2(float64(cell.candidates)))
			if e := sud.setCell(cell.row, cell.col, value); e != nil {
				return false, e
			}
			log.Printf("Cell %d,%d has single value %d", cell.row, cell.col, value)
			changes = true
		}
		return changes, nil
	})
}

// where a value is only fits into 1 of the 9 cells of a row/col/block
func hiddenSingle(sud *SudokuSquare) (bool, error) {
	return applyToNonagons(sud, func(nona nonagon) (bool, error) {
		changes := false
		for val := 1; val <= 9; val++ {
			availablePlaces := 0
			idx := -1
			for i, cell := range nona.cells {
				if cell.hasCandidate(val) {
					availablePlaces++
					idx = i
					if availablePlaces >= 2 {
						break
					}
				}
			}
			if availablePlaces == 1 {
				cell := nona.cells[idx]
				if e := sud.setCell(cell.row, cell.col, val); e != nil {
					return false, e
				}
				changes = true
				log.Printf("%s has hidden single at %d,%d for value %d",
					nona.name, cell.row, cell.col, val)
			}
		}
		return changes, nil
	})
}

// if a candidate has only 2 or 3 available cells in a block that are along a line, then
// that's a pointing pair that removes that candiate as an option along that line outside
// the block
func pointingPair(sud *SudokuSquare) (bool, error) {
	return applyToBlocks(
		sud,
		func(r1, r2, c1, c2 int) (bool, error) {
			return pointingPairBlockFn(sud, r1, r2, c1, c2)
		})
}

// for each value, check each block row or column to see if it has a pointing pair for that value
func pointingPairBlockFn(sud *SudokuSquare, blkRowStart, blkRowEnd, blkColStart, blkColEnd int) (bool, error) {
	changes := false
	for val := 1; val <= 9; val++ {
		alignedPlaces := 0
		pointingPairRow := -1
		pointingPairCol := -1

		for row := blkRowStart; row <= blkRowEnd; row++ {
			for col := blkColStart; col <= blkColEnd; col++ {

				cell := sud.cells[row][col]
				if cell.hasCandidate(val) {
					alignedPlaces++
					if pointingPairRow == -1 {
						pointingPairRow = row
					} else if pointingPairRow != row {
						// fail, not same row
						pointingPairRow = -2
					}
					if pointingPairCol == -1 {
						pointingPairCol = col
					} else if pointingPairCol != col {
						// fail, not same col
						pointingPairCol = -2
					}
				}
			}
		}

		hasHorizontalPointingPair := pointingPairRow >= 0 && alignedPlaces > 1
		if hasHorizontalPointingPair {
			impacting := false
			for col := 0; col < 9; col++ {
				if col < blkColStart || col > blkColEnd {
					cell := &sud.cells[pointingPairRow][col]
					if cell.hasCandidate(val) {
						cell.removeCandidate(val)
						impacting = true
					}
				}
			}
			if impacting {
				log.Printf("  row  %d pointing pair for value %d", pointingPairRow, val)
			}
			changes = changes || impacting
		}

		hasVerticalPointingPair := pointingPairCol >= 0 && alignedPlaces > 1
		if hasVerticalPointingPair {
			impacting := false
			for row := 0; row < 9; row++ {
				if row < blkRowStart || row > blkRowEnd {
					cell := &sud.cells[row][pointingPairCol]
					if cell.hasCandidate(val) {
						cell.removeCandidate(val)
						impacting = true
					}
				}
			}
			if impacting {
				log.Printf("column %d pointing pair for value %d", pointingPairCol, val)
			}
			changes = changes || impacting
		}
	}

	return changes, nil
}

// When the only candidate cells of a value in a row or column occur in the same
// block then that value can be removed as a candidate from all the other cells
// of that block.
// Is called a 'claimingPair' to distinguish it from a hiddenSingle but the algorithm
// below doesn't care if it's a single/pair/triple.
func claimingPair(sud *SudokuSquare) (bool, error) {
	changes := false
	for val := 1; val <= 9; val++ {
		// for each block
		for blockStartRow := 0; blockStartRow < 9; blockStartRow += 3 {
			for blockStartCol := 0; blockStartCol < 9; blockStartCol += 3 {
				blockEndRow, blockEndCol := blockStartRow+2, blockStartCol+2

				// horizontal
				for testRow := blockStartRow; testRow <= blockEndRow; testRow++ {

					availableOnRowInsideBlock := false
					for col := blockStartCol; col <= blockEndCol; col++ {
						cell := sud.cells[testRow][col]
						if cell.hasCandidate(val) {
							availableOnRowInsideBlock = true
						}
					}
					availableOnRowOutsideBlock := false
					for col := 0; col < 9; col++ {
						if col < blockStartCol || col > blockEndCol {
							cell := sud.cells[testRow][col]
							if cell.hasCandidate(val) {
								availableOnRowOutsideBlock = true
							}
						}
					}
					if availableOnRowInsideBlock && !availableOnRowOutsideBlock {
						impacting := false
						for row := blockStartRow; row <= blockEndRow; row++ {
							if row != testRow {
								for col := blockStartCol; col <= blockEndCol; col++ {
									cell := &sud.cells[row][col]
									if cell.hasCandidate(val) {
										cell.removeCandidate(val)
										impacting = true
									}
								}
							}
						}
						if impacting {
							log.Printf("  row  %d claiming pair for value %d", testRow, val)
						}
						changes = changes || impacting
					}
				}

				// vertical
				for testCol := blockStartCol; testCol <= blockEndCol; testCol++ {
					availableOnColInsideBlock := false
					for row := blockStartRow; row <= blockEndRow; row++ {
						cell := sud.cells[row][testCol]
						if cell.hasCandidate(val) {
							availableOnColInsideBlock = true
						}
					}
					availableOnColOutsideBlock := false
					for row := 0; row < 9; row++ {
						if row < blockStartRow || row > blockEndRow {
							cell := sud.cells[row][testCol]
							if cell.hasCandidate(val) {
								availableOnColOutsideBlock = true
							}
						}
					}
					if availableOnColInsideBlock && !availableOnColOutsideBlock {
						impacting := false
						for col := blockStartCol; col <= blockEndCol; col++ {
							if col != testCol {
								for row := blockStartRow; row <= blockEndRow; row++ {
									cell := &sud.cells[row][col]
									if cell.hasCandidate(val) {
										cell.removeCandidate(val)
										impacting = true
									}
								}
							}
						}
						if impacting {
							log.Printf("column %d claiming pair for value %d", testCol, val)
						}
						changes = changes || impacting
					}
				}
			}
		}
	}
	return changes, nil
}
