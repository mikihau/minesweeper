package main

import (
	"math/rand"
	"testing"
)

func TestReveal(t *testing.T) {
	/*
		initialized board:
			# 3 # # 1 0 0 0
			# 4 2 3 2 1 0 0
			# 3 1 2 # 1 0 0
			1 3 # 3 1 1 0 0
			0 2 # 2 0 0 0 0
			0 1 1 1 1 1 1 0
			0 0 0 1 2 # 1 0
			0 0 0 1 # 2 1 0
	*/
	cases := []struct {
		x                   int
		y                   int
		expectedNumRevealed int
	}{
		{7, 0, 16},
		{2, 3, 1},
		{0, 0, 1},
	}
	for _, testCase := range cases {
		rand.Seed(5)
		board := makeBoard(8, 8, 10)
		board.updateOnReveal(testCase.x, testCase.y)
		if board.numRevealed != testCase.expectedNumRevealed {
			t.Errorf("Revealing coord (%v, %v) should reveal %v cells, but got %v: \n%v", testCase.x, testCase.y, testCase.expectedNumRevealed, board.numRevealed, board)
		}
	}
}

func TestFlag(t *testing.T) {
	rand.Seed(5)
	board := makeBoard(8, 8, 10)
	board.updateOnFlag(0, 0)
	if board.numFlagged != 1 {
		t.Errorf("Flagging coord (%v, %v) should flag %v cells, but got %v: \n%v", 0, 0, 1, board.numFlagged, board)
	}
	board.updateOnReveal(7, 0)
	board.updateOnFlag(7, 0)
	if board.numFlagged != 1 {
		t.Errorf("Flagging coord (%v, %v) should flag %v cells, but got %v: \n%v", 0, 0, 1, board.numFlagged, board)
	}
}
