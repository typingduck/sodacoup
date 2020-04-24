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
