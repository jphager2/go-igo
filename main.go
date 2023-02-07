package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Stone interface {
	Format(fmt.State, rune)
}

type LibertyStone int
type BlackStone int
type WhiteStone int

const (
	Liberty = LibertyStone(0)
	Black   = BlackStone(1)
	White   = WhiteStone(2)
)

func (s BlackStone) Format(f fmt.State, c rune) {
	f.Write([]byte("X"))
}

func (s WhiteStone) Format(f fmt.State, c rune) {
	f.Write([]byte("0"))
}

func (s LibertyStone) Format(f fmt.State, c rune) {
	f.Write([]byte("+"))
}

type Board struct {
	Size  int
	Board [][]Stone
}

type StoneGroup struct {
	Stone     Stone
	Liberties [][2]int
	Coords    [][2]int
}

func (sg *StoneGroup) Include(coord [2]int) bool {
	for _, sc := range sg.Coords {
		if sc == coord {
			return true
		}
	}

	return false
}

func (sg *StoneGroup) IncludeLiberty(coord [2]int) bool {
	for _, sc := range sg.Liberties {
		if sc == coord {
			return true
		}
	}

	return false
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

func (b *Board) CoordsAround(x int, y int) [][2]int {
	around := make([][2]int, 0, 5)
	around = append(around, [2]int{x, y})
	if x > 0 {
		around = append(around, [2]int{x - 1, y})
	}
	if x < b.Size-2 {
		around = append(around, [2]int{x + 1, y})
	}
	if y > 0 {
		around = append(around, [2]int{x, y - 1})
	}
	if y < b.Size-2 {
		around = append(around, [2]int{x, y + 1})
	}

	return around
}

// Get's the color of the stone at coord and fill the group coord and counts
// the groups liberties. Does not form a group if the stone is a liberty
func (b *Board) MustGroup(sg *StoneGroup, coord [2]int) bool {
	if sg.Stone == Liberty {
		return false
	}

	if sg.Include(coord) {
		return false
	}

	sg.Coords = append(sg.Coords, coord)

	around := b.CoordsAround(coord[0], coord[1])
	for _, acoord := range around {
		if acoord == coord {
			continue
		}

		as := b.Board[acoord[1]][acoord[0]]

		if as == sg.Stone {
			_ = b.MustGroup(sg, acoord)
		} else if as == Liberty {
			if !sg.IncludeLiberty(acoord) {
				sg.Liberties = append(sg.Liberties, acoord)
			}
		}
	}

	return true
}

func (b *Board) Groups(x int, y int) ([]StoneGroup, error) {
	if x < 0 || x >= b.Size || y < 0 || y > b.Size {
		return nil, errors.New("Coord out of bounds.")
	}

	groups := make([]StoneGroup, 0, 5)
	around := b.CoordsAround(x, y)
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
		if ok := b.MustGroup(sg, coord); ok {
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

func PlayGame() {
	captures := map[string]int{}
	passed := false
	blackTurn := true
	board := NewBoard(19)

	// Play Loop
	for {
		var player string
		var stone Stone

		if blackTurn {
			player = "Black"
			stone = Black
		} else {
			player = "White"
			stone = White
		}

		fmt.Println(board)
		fmt.Printf("Black Captures: %v\n", captures["Black"])
		fmt.Printf("White Captures: %v\n", captures["White"])
		fmt.Printf("%s's turn. Enter coordintates (`x,y`), `pass`, `resign`: ", player)

		var input string
		fmt.Scanln(&input)

		input = strings.TrimSpace(input)

		if input == "resign" {
			fmt.Printf("%s has resigned.\n", player)
			fmt.Println("The game is over.")
			break
		}

		if input == "pass" {
			fmt.Printf("%s has passed.\n", player)

			if passed {
				fmt.Println("The game is over.")
				break
			} else {
				blackTurn = !blackTurn
				passed = true
				continue
			}
		}

		passed = false

		coords := strings.Split(input, ",")
		if len(coords) != 2 {
			fmt.Printf("Invalid input: `%s`\n", input)
			continue
		}

		x, xerr := strconv.Atoi(coords[0])
		y, yerr := strconv.Atoi(coords[1])
		if xerr != nil || yerr != nil {
			fmt.Printf("Invalid input: `%s`\n", input)
			continue
		}

		err := board.Place(stone, x, y)
		if err != nil {
			fmt.Printf("Cannot place %s stone at `%s`!\n", player, input)
			continue
		}
		groups, _ := board.Groups(x, y)
		placedGroup := groups[0]
		otherGroups := groups[1:]

		captured := false
		for _, sg := range otherGroups {
			if len(sg.Liberties) == 0 {
				captured = true
				for _, sc := range sg.Coords {
					err = board.Remove(sc[0], sc[1])
					if err != nil {
						panic("This cannot happen: tried to remove a stone that doesn't exist")
					}
					captures[player]++
				}
			}
		}

		if !captured && len(placedGroup.Liberties) == 0 {
			err := board.Remove(x, y)
			if err != nil {
				panic("This cannot happen: tried to remove a stone that doesn't exist")
			}

			fmt.Printf("Cannot play suicide at `%s`!\n", input)
			continue
		}

		blackTurn = !blackTurn

		fmt.Println()
	}
}

func main() {
	PlayGame()
}
