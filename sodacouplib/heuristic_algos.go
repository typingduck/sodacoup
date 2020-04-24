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
