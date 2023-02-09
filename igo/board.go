package igo

import (
	"errors"
	"fmt"
)

type Board struct {
	Size  int
	Board [][]Stone
}

func NewBoard(size int) *Board {
	board := make([][]Stone, size, size)

	for i := 0; i < size; i++ {
		board[i] = make([]Stone, size, size)

		for j := 0; j < size; j++ {
			board[i][j] = Liberty
		}
	}

	return &Board{Size: size, Board: board}
}

func (b *Board) Place(stone Stone, x int, y int) error {
	if x < 0 || x >= b.Size || y < 0 || y > b.Size {
		return errors.New("Coord out of bounds.")
	}

	if b.Board[y][x] != Liberty {
		return errors.New("Coord already occupied by stone.")
	}

	b.Board[y][x] = stone

	return nil
}

func (b *Board) Remove(x int, y int) error {
	if x < 0 || x >= b.Size || y < 0 || y > b.Size {
		return errors.New("Coord out of bounds.")
	}

	if b.Board[y][x] == Liberty {
		return errors.New("Coord already empty.")
	}

	b.Board[y][x] = Liberty

	return nil
}

func (b *Board) coordsAround(x int, y int) [][2]int {
	around := make([][2]int, 0, 5)
	around = append(around, [2]int{x, y})
	if x > 0 {
		around = append(around, [2]int{x - 1, y})
	}
	if x < b.Size-1 {
		around = append(around, [2]int{x + 1, y})
	}
	if y > 0 {
		around = append(around, [2]int{x, y - 1})
	}
	if y < b.Size-1 {
		around = append(around, [2]int{x, y + 1})
	}

	return around
}

// Get's the color of the stone at coord and fill the group coord and counts
// the groups liberties. Does not form a group if the stone is a liberty
func (b *Board) stoneGroup(sg *StoneGroup, coord [2]int) bool {
	if sg.Stone == Liberty {
		return false
	}

	if sg.Include(coord) {
		return false
	}

	sg.Coords = append(sg.Coords, coord)

	around := b.coordsAround(coord[0], coord[1])
	for _, acoord := range around {
		if acoord == coord {
			continue
		}

		as := b.Board[acoord[1]][acoord[0]]

		if as == sg.Stone {
			_ = b.stoneGroup(sg, acoord)
		} else if as == Liberty {
			if !sg.IncludeLiberty(acoord) {
				sg.Liberties = append(sg.Liberties, acoord)
			}
		}
	}

	return true
}

func (b *Board) libGroup(sg *StoneGroup, coord [2]int) bool {
	if sg.Include(coord) {
		return false
	}

	sg.Coords = append(sg.Coords, coord)

	around := b.coordsAround(coord[0], coord[1])
	for _, acoord := range around {
		if acoord == coord {
			continue
		}

		as := b.Board[acoord[1]][acoord[0]]

		if as == Liberty {
			b.libGroup(sg, acoord)
		}
	}

	return true
}

func (b *Board) LibGroups() []StoneGroup {
	var groups []StoneGroup

	for x := 0; x < b.Size; x++ {
		for y := 0; y < b.Size; y++ {
			if b.Board[y][x] != Liberty {
				continue
			}

			coord := [2]int{x, y}
			found := false

			for _, sg := range groups {
				if sg.Include(coord) {
					found = true
					break
				}
			}

			if found {
				continue
			}

			sg := &StoneGroup{Stone: Liberty}
			if ok := b.libGroup(sg, coord); ok {
				groups = append(groups, *sg)
			}
		}
	}

	return groups
}

func (b *Board) Groups(x int, y int) ([]StoneGroup, error) {
	if x < 0 || x >= b.Size || y < 0 || y > b.Size {
		return nil, errors.New("Coord out of bounds.")
	}

	groups := make([]StoneGroup, 0, 5)
	around := b.coordsAround(x, y)
	for _, coord := range around {
		found := false

		for _, sg := range groups {
			if sg.Include(coord) {
				found = true
				break
			}
		}

		if found {
			continue
		}

		stone := b.Board[coord[1]][coord[0]]
		sg := &StoneGroup{Stone: stone}
		if ok := b.stoneGroup(sg, coord); ok {
			groups = append(groups, *sg)
		}
	}

	return groups, nil
}

func (b *Board) Format(f fmt.State, c rune) {
	for _, row := range b.Board {
		for _, stone := range row {
			stone.Format(f, c)
			f.Write([]byte(" "))
		}
		f.Write([]byte("\n"))
	}
}
