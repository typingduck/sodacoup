package sodacouplib

import (
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
)

func TestSampleGeneration(t *testing.T) {
	rand.Seed(0xC25)
	s, e := GenerateProblem()
	if e != nil {
		t.Fatal("failed to generate:", e)
	}
	expected, _ := FormatSudoku(`
		___ ___ 8__
		___ 2_8 9__
		___ 4_7 __5
		
		6_4 _3_ _2_
		8__ _9_ ___
		___ ___ _1_
		
		__1 ___ ___
		_42 __5 _3_
		___ 9__ ___
	`)
	result, _ := FormatSudoku(s.String())
	assert.Equal(t, expected, result)
}

func TestSolvable(t *testing.T) {
	s, e := GenerateProblem()
	if e != nil {
		t.Fatal("failed to generate:", e)
	}
	e = s.Solve()
	if e != nil {
		t.Fatal("failed to solve:", e)
	}
	assert.Equal(t, true, isSolved(s))
}
