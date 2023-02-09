package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/jphager2/go-igo/igo"
)

func main() {
	captures := map[string]int{}
	passed := false
	blackTurn := true
	board := igo.NewBoard(19)

	// Play Loop
	for {
		var player string
		var stone igo.Stone

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
