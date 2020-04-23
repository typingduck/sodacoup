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
