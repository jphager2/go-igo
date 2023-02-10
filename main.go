package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/jphager2/go-igo/igo"
)

func main() {
	var player string
	var stone igo.Stone
	var board *igo.Board
	var capturedCount int
	var captured []igo.StoneGroup

	captures := map[string]int{}
	passed := false
	blackTurn := true
	koStoneCoord := [2]int{-1, -1}

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

	// Play Loop
	for {
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

		capturedCount = 0
		captured = make([]igo.StoneGroup, 0, len(otherGroups))
		for _, sg := range otherGroups {
			if len(sg.Liberties) == 0 {
				captured = append(captured, sg)
				capturedCount += len(sg.Coords)
			}
		}

		// Check suicide rule and undo if fails
		if capturedCount == 0 && len(placedGroup.Liberties) == 0 {
			err := board.Remove(x, y)
			if err != nil {
				panic("This cannot happen: tried to remove a stone that doesn't exist")
			}

			fmt.Printf("Cannot play suicide at `%s`!\n", input)
			continue
		}

		// Check ko rule and undo if fails. Also keep track of ko stones
		if capturedCount == 1 {
			fmt.Printf("CaptureCount: %v, CapturedGroups: %v\n", capturedCount, len(captured))
			fmt.Printf("Coord: %v, KoStoneCoord: %v\n", [2]int{x, y}, koStoneCoord)

			if [2]int{x, y} == koStoneCoord {
				err := board.Remove(x, y)
				if err != nil {
					panic("This cannot happen: tried to remove a stone that doesn't exist")
				}

				fmt.Printf("Cannot play ko at `%s`!\n", input)
				continue
			}

			koStoneCoord = captured[0].Coords[0]
		} else {
			koStoneCoord = [2]int{-1, -1}
		}

		// Commit the captures, since there will be no more undos
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

		fmt.Println()
	}
}
