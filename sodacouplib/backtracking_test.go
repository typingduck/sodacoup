package sodacouplib

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBacktracking_unsolvableProblem(t *testing.T) {
	unsolvableProblem := `
		1_3 456 729
		426 789 1_3
		789 123 456

		214 365 897
		365 897 214
		897 214 365

		531 642 978
		642 978 531
		978 531 642
	`
	s, err := NewSudokuSquare(unsolvableProblem)
	if err != nil {
		t.Fatal("got unexpected error from valid input:", err)
	}
	_, err = backTrack(s)
	if err == nil {
		t.Fatal("got unexpected success (!) on bad problem")
	}
	assert.Contains(t, err.Error(), "failed to converge")
}

func TestBacktracking_edgeCases(t *testing.T) {
	// algorithm doesn't check input is a valid sudoku so just
	// document here how it behaves for these cases.
	t.Run("all empty square", func(t *testing.T) {
		problem := `
			___ ___ ___
			___ ___ ___
			___ ___ ___

			___ ___ ___
			___ ___ ___
			___ ___ ___

			___ ___ ___
			___ ___ ___
			___ ___ ___
		`
		solutionItShouldConvergeTo := `
			123 456 789
			456 789 123
			789 123 456

			214 365 897
			365 897 214
			897 214 365

			531 642 978
			642 978 531
			978 531 642
		`
		testRunBackTrack(t, problem, solutionItShouldConvergeTo)
	})
	t.Run("simple element", func(t *testing.T) {
		// algorithm doesn't guarantee uniqueness. So just
		// document how it behaves.
		problem := `
			___ ___ ___
			___ ___ ___
			___ ___ ___

			___ ___ ___
			___ _1_ ___
			___ ___ ___

			___ ___ ___
			___ ___ ___
			___ ___ ___
		`
		solutionItShouldConvergeTo := `
			123 456 789
			456 789 123
			789 123 456

			214 365 897
			367 918 245
			598 247 361

			631 892 574
			845 671 932
			972 534 618
		`
		testRunBackTrack(t, problem, solutionItShouldConvergeTo)
	})
}

func TestBacktracking_validProblems(t *testing.T) {
	sampleProblems := []struct{ name, problem, solution string }{{
		"simple missing 8s",
		`
		123 456 7_9
		456 7_9 123
		7_9 123 456

		234 567 _91
		567 _91 234
		_91 234 567

		345 67_ 912
		67_ 912 345
		912 345 67_
		`, `
		123 456 789
		456 789 123
		789 123 456

		234 567 891
		567 891 234
		891 234 567

		345 678 912
		678 912 345
		912 345 678
	`}, {
		"problem 1",
		`
		__5 __2 __4
		___ 5__ ___
		_9_ _7_ 8_1

		___ 3__ ___
		5__ 81_ 2_3
		__6 ___ __7

		_39 64_ ___
		___ ___ ___
		__7 __5 _2_

		`, `

		185 962 374
		743 581 962
		692 473 851

		928 357 146
		574 816 293
		316 294 587

		239 648 715
		451 729 638
		867 135 429
	`}, {
		"problem 2",
		`
		__8 7_4 ___
		45_ 82_ _36
		2_3 6__ 9__

		_12 _87 ___
		_9_ 2_3 _5_
		___ 14_ 89_

		__7 __6 3_4
		64_ _78 _21
		___ 4_2 6__

		`, `
		968 734 215
		451 829 736
		273 651 948

		512 987 463
		894 263 157
		736 145 892

		127 596 384
		649 378 521
		385 412 679
	`}, {
		"problem 3",
		`
		2__ 4__ 6__
		_13 _28 __7
		_76 _5_ 8__

		9__ ___ _6_
		__5 ___ 3__
		_3_ ___ __9

		__4 _3_ 92_
		3__ 74_ 51_
		__8 __5 __3
		`, `
		289 473 651
		513 628 497
		476 159 832

		947 382 165
		625 917 384
		831 564 279

		754 831 926
		392 746 518
		168 295 743
	`}, {
		"problem 4",
		`
		___ ___ _92
		___ 2_6 8_7
		___ _71 5__

		_64 3__ ___
		95_ ___ _86
		___ __8 45_

		__6 48_ ___
		4_8 6_5 ___
		72_ ___ ___
		`, `
		687 543 192
		541 296 837
		239 871 564

		864 352 719
		953 714 286
		172 968 453

		316 487 925
		498 625 371
		725 139 648
	`}, {
		"problem 5",
		`
		___ __9 ___
		_9_ ___ _65
		8__ 3__ ___

		__3 ___ __6
		___ 7__ 82_
		__1 ___ 34_

		__5 8__ ___
		___ _37 ___
		62_ 1__ __9
		`, `
		162 579 483
		397 284 165
		854 361 972

		273 418 596
		946 753 821
		581 926 347

		735 892 614
		419 637 258
		628 145 739
		`},
	}
	for _, tc := range sampleProblems {
		tc := tc // for parallel
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			testRunBackTrack(t, tc.problem, tc.solution)
		})
	}
}

func testRunBackTrack(t *testing.T, problem, expected string) {
	if expectedF, err := FormatSudoku(expected); err != nil {
		t.Error("got unexpected error from valid input:", err)
	} else if s, err := NewSudokuSquare(problem); err != nil {
		t.Error("got unexpected error from valid input:", err)
	} else if _, err := backTrack(s); err != nil {
		t.Error("got unexpected error from solving:", err)
	} else if result, err := FormatSudoku(s.String()); err != nil {
		t.Error("got unexpected error from formatting result:", err)
	} else if result != expectedF {
		t.Errorf("unmatched!\nwanted:\n%s\ngot:\n%s", expectedF, result)
	}
}
