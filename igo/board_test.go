package igo

import (
	"testing"
)

func TestGroups(t *testing.T) {
	board := NewBoard(9)
	_ = board.Place(Black, 0, 0)
	_ = board.Place(Black, 1, 0)
	_ = board.Place(Black, 0, 1)
	_ = board.Place(Black, 1, 1)

	groups, err := board.Groups(0, 0)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if len(groups) != 1 || len(groups[0].Coords) != 4 {
		t.Errorf("Expected to get a single group of four black stones")
	}

	// X X X X X 0 0 0 0
	// X X X X X 0 0 0 0
	// X X X X X 0 0 0 0
	// X X X X X 0 0 0 0
	// X X X X X 0 0 0 0
	// 0 0 0 0 0 X X X X
	// 0 0 0 0 0 X X X X
	// 0 0 0 0 0 X X X X
	// 0 0 0 0 0 X X X X

	board = NewBoard(9)
	for i := 0; i < 5; i++ {
		for j := 0; j < 5; j++ {
			_ = board.Place(Black, i, j)
		}
	}

	for i := 5; i < 9; i++ {
		for j := 0; j < 5; j++ {
			_ = board.Place(White, i, j)
		}
	}

	for i := 5; i < 9; i++ {
		for j := 5; j < 9; j++ {
			_ = board.Place(Black, i, j)
		}
	}

	for i := 0; i < 5; i++ {
		for j := 5; j < 9; j++ {
			_ = board.Place(White, i, j)
		}
	}

	groups, err = board.Groups(4, 4)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if len(groups) != 3 {
		t.Errorf("Expected to get 3 groups, got: %v", len(groups))
	}

	if groups[0].Stone != Black || len(groups[0].Coords) != 25 || !groups[0].Include(Coord{4, 4}) {
		t.Errorf("Expected first group to be the placed group of 25 black stones, got %T, %v, %v", groups[0].Stone, len(groups[0].Coords), groups[0].Coords)
	}

	board.Remove(4, 4)
	_ = board.Place(White, 4, 4)

	groups, err = board.Groups(4, 4)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if len(groups) != 2 {
		t.Errorf("Expected to get 2 groups, got: %v", len(groups))
	}

	if groups[0].Stone != White || len(groups[0].Coords) != 41 || !groups[0].Include(Coord{4, 4}) {
		t.Errorf("Expected first group to be the placed group of 41 white stones, got %T, %v, %v", groups[0].Stone, len(groups[0].Coords), groups[0].Coords)
	}

	board = NewBoard(9)
	_ = board.Place(Black, 4, 3)
	_ = board.Place(Black, 4, 5)
	_ = board.Place(Black, 3, 4)
	_ = board.Place(Black, 5, 4)

	_ = board.Place(White, 4, 4)

	groups, err = board.Groups(4, 4)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if len(groups) != 5 {
		t.Errorf("Expected to get 5 groups, got: %v", len(groups))
	}

	if groups[0].Stone != White || len(groups[0].Coords) != 1 || !groups[0].Include(Coord{4, 4}) {
		t.Errorf("Expected first group to be the placed group of 1 white stones, got %T, %v, %v", groups[0].Stone, len(groups[0].Coords), groups[0].Coords)
	}
}

func TestLibGroups(t *testing.T) {
	board := NewBoard(9)
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			_ = board.Place(Black, i, j)
		}
	}

	board.Remove(1, 1)
	board.Remove(1, 2)
	board.Remove(2, 1)
	board.Remove(2, 2)

	board.Remove(1, 7)
	board.Remove(1, 6)
	board.Remove(2, 7)
	board.Remove(2, 6)

	board.Remove(6, 1)
	board.Remove(6, 2)
	board.Remove(7, 1)
	board.Remove(7, 2)

	board.Remove(6, 7)
	board.Remove(6, 6)
	board.Remove(7, 7)
	board.Remove(7, 6)

	libs := board.LibGroups()

	if len(libs) != 4 {
		t.Errorf("Expected there to be 4 lib groups, got %v", len(libs))
	}

	for _, group := range libs {
		if len(group.Coords) != 4 {
			t.Errorf("Expected there to be 4 liberties in each group, got %v, %+v", len(group.Coords), group)
		}
	}

	board.Remove(4, 4)

	libs = board.LibGroups()
	counts := map[int]int{}

	if len(libs) != 5 {
		t.Errorf("Expected there to be 5 lib groups, got %v", len(libs))
	}

	for _, group := range libs {
		counts[len(group.Coords)]++
	}

	if counts[1] != 1 && counts[4] != 4 {
		t.Errorf("Expected there to be 4 groups with 4 liberties and 1 group with 1 liberty, got %v", len(libs))
	}
}
