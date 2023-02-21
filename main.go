package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/disiqueira/gotree"
	"github.com/jphager2/go-igo/igo"
)

type Move struct {
	Number    int
	Board     igo.Board
	Passed    bool
	BlackTurn bool
	Move      igo.Coord
	KoStone   igo.Coord
	Captures  map[string]int
	Resigned  bool
	RootMove  bool
	HeadMove  *Move
	Branches  []*Move
}

func gotreeFromBranches(m *Move, t gotree.Tree) {
	for _, b := range m.Branches {
		var stone igo.Stone
		var move string

		if b.BlackTurn {
			stone = igo.Black
		} else {
			stone = igo.White
		}

		if b.Resigned {
			move = "Resign"
		} else if b.Passed {
			move = "Pass"
		} else {
			move = fmt.Sprintf("%d,%d", b.Move[0], b.Move[1])
		}

		tree := t.Add(fmt.Sprintf("%v(%d) [%v]", stone, b.Number, move))
		gotreeFromBranches(b, tree)
	}
}

func appendAsBranch(m *Move) {
	found := false

	for _, b := range m.HeadMove.Branches {
		found = b.Move == m.Move && b.Resigned == m.Resigned && b.Passed == m.Passed
	}

	if !found {
		m.HeadMove.Branches = append(m.HeadMove.Branches, m)
	}
}

func main() {
	var board *igo.Board

	captures := map[string]int{}
	blackTurn := true
	passed := false
	koStoneCoord := igo.Coord{-1, -1}

	for {
		fmt.Printf("Let's Go! What board size (9, 13, 19): ")

		var input string
		fmt.Scanln(&input)

		input = strings.TrimSpace(input)

		size, err := strconv.Atoi(input)

		if err != nil {
			fmt.Printf("Invalid input: `%s`\n", input)
			continue
		}

		board = igo.NewBoard(size)
		break
	}

	rootMove := &Move{
		Board:     *board.Dup(),
		Passed:    passed,
		BlackTurn: false,
		Move:      igo.Coord{-1, -1},
		KoStone:   koStoneCoord,
		Captures:  captures,
		RootMove:  true,
	}
	lastMove := rootMove

	// Play Loop
	for {
		var player string
		var stone igo.Stone
		var capturedCount int
		var captured []igo.StoneGroup

		if blackTurn {
			player = "Black"
			stone = igo.Black
		} else {
			player = "White"
			stone = igo.White
		}

		fmt.Println(board)
		fmt.Printf("Black Captures: %v\n", captures["Black"])
		fmt.Printf("White Captures: %v\n", captures["White"])
		fmt.Printf(
			"%s's turn. Enter coordintates (`x,y`), `pass`, `undo`, `resign`, `branches`: ",
			player,
		)

		move := igo.Coord{-1, -1}
		var input string
		fmt.Scanln(&input)

		input = strings.TrimSpace(input)

		if input == "resign" {
			fmt.Printf("%s has resigned.\n", player)
			fmt.Println("The game is over.")
			lastMove = &Move{
				Number:    lastMove.Number + 1,
				Board:     *board.Dup(),
				Passed:    true,
				BlackTurn: !blackTurn,
				Move:      move,
				KoStone:   igo.Coord{-1, -1},
				Captures:  captures,
				Resigned:  true,
				HeadMove:  lastMove,
			}
			appendAsBranch(lastMove)
			break
		}

		if input == "pass" {
			fmt.Printf("%s has passed.\n", player)

			if passed {
				fmt.Println("The game is over.")
				break
			}

			passed = true
			blackTurn = !blackTurn
			lastMove = &Move{
				Number:    lastMove.Number + 1,
				Board:     *board.Dup(),
				Passed:    true,
				BlackTurn: !blackTurn,
				Move:      move,
				KoStone:   igo.Coord{-1, -1},
				Captures:  captures,
				HeadMove:  lastMove,
			}
			appendAsBranch(lastMove)
			continue
		}

		if input == "undo" {
			if lastMove.RootMove {
				fmt.Println("Already at the beginning of the game.")
				continue
			}

			lastMove = lastMove.HeadMove
			blackTurn = !blackTurn
			board = lastMove.Board.Dup()
			passed = lastMove.Passed
			koStoneCoord = lastMove.KoStone
			captures = lastMove.Captures
			continue
		}

		if input == "branches" {
			tree := gotree.New("Game")
			gotreeFromBranches(rootMove, tree)
			fmt.Println(tree.Print())

			continue
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
		move = igo.Coord{x, y}

		err := board.Place(stone, x, y)
		if err != nil {
			fmt.Printf("Cannot place %s stone at `%s`!\n", player, input)
			continue
		}
		groups, _ := board.Groups(x, y)
		placedGroup := groups[0]
		otherGroups := groups[1:]

		capturedCount = 0
		captured = make([]igo.StoneGroup, 0, len(otherGroups))
		for _, sg := range otherGroups {
			if len(sg.Liberties) == 0 {
				captured = append(captured, sg)
				capturedCount += len(sg.Coords)
			}
		}

		// Check suicide rule and revert if fails
		if capturedCount == 0 && len(placedGroup.Liberties) == 0 {
			err := board.Remove(x, y)
			if err != nil {
				panic("This cannot happen: tried to remove a stone that doesn't exist")
			}

			fmt.Printf("Cannot play suicide at `%s`!\n", input)
			continue
		}

		// Check ko rule and revert if fails. Also keep track of ko stones
		if capturedCount == 1 {
			fmt.Printf(
				"CaptureCount: %v, CapturedGroups: %v\n",
				capturedCount,
				len(captured),
			)
			fmt.Printf("Coord: %v, KoStoneCoord: %v\n", move, koStoneCoord)

			if move == koStoneCoord {
				err := board.Remove(x, y)
				if err != nil {
					panic("This cannot happen: tried to remove a stone that doesn't exist")
				}

				fmt.Printf("Cannot play ko at `%s`!\n", input)
				continue
			}

			koStoneCoord = captured[0].Coords[0]
		} else {
			koStoneCoord = igo.Coord{-1, -1}
		}

		// Commit the captures, since there will be no more reverts
		captures[player] += capturedCount

		for _, sg := range captured {
			for _, sc := range sg.Coords {
				err = board.Remove(sc[0], sc[1])
				if err != nil {
					panic("This cannot happen: tried to remove a stone that doesn't exist")
				}
			}
		}

		blackTurn = !blackTurn
		lastMove = &Move{
			Number:    lastMove.Number + 1,
			Board:     *board.Dup(),
			BlackTurn: !blackTurn,
			Move:      move,
			KoStone:   koStoneCoord,
			Captures:  captures,
			HeadMove:  lastMove,
		}
		appendAsBranch(lastMove)

		fmt.Println()
	}
}
