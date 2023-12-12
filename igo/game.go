package igo

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/disiqueira/gotree"
)

type Move struct {
	Number    int
	Board     Board
	Passed    bool
	BlackTurn bool
	Move      Coord
	KoStone   Coord
	Captures  map[string]int
	Resigned  bool
	RootMove  bool
	HeadMove  *Move
	Branches  []*Move
}

func gotreeFromBranches(m *Move, t gotree.Tree) {
	for _, b := range m.Branches {
		var stone Stone
		var move string

		if b.BlackTurn {
			stone = Black
		} else {
			stone = White
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

type Game struct {
	Captures     map[string]int
	BlackTurn    bool
	Passed       bool
	KoStoneCoord Coord
	RootMove     *Move
	LastMove     *Move
	Board        *Board
	Resigned     bool
	Over         bool
	ShowBranches bool
	Input        string
	ErrMsg       string
}

func NewGame(boardSize int) Game {
	board := NewBoard(boardSize)
	nilCoord := Coord{-1, -1}
	captures := map[string]int{}

	rootMove := &Move{
		Board:     *board.Dup(),
		Passed:    false,
		BlackTurn: false,
		Move:      nilCoord,
		KoStone:   nilCoord,
		Captures:  captures,
		RootMove:  true,
	}

	return Game{
		Captures:     captures,
		BlackTurn:    true,
		KoStoneCoord: nilCoord,
		RootMove:     rootMove,
		LastMove:     rootMove,
		Board:        board,
	}
}

func (game Game) Player() string {
	if game.BlackTurn {
		return "Black"
	} else {
		return "White"
	}
}

type AskForInputMsg struct{}

func askForInput() tea.Msg {
	return AskForInputMsg{}
}

func (game Game) Init() tea.Cmd {
	return askForInput
}

func (game Game) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case AskForInputMsg:
		game.Input = ""
	case tea.KeyMsg:
		if game.ShowBranches {
			game.ShowBranches = false
			return game, cmd
		}
		switch msg.String() {
		case "ctrl+c":
			return game, tea.Quit
		case "backspace":
			if len(game.Input) > 0 {
				game.Input = game.Input[0 : len(game.Input)-1]
			}
		case "enter":
			// play etc
			var ok bool
			if game, ok = game.Play(); ok {
				cmd = askForInput
			}
			if game.Over {
				cmd = tea.Quit
			}
		default:
			game.Input = game.Input + msg.String()
		}
	default:
		return game, cmd
	}

	return game, cmd
}

func (game Game) View() string {
	s := strings.Builder{}

	if game.ErrMsg != "" {
		s.WriteString(game.ErrMsg + "\n")
	}

	if game.ShowBranches {
		tree := gotree.New("Game")
		gotreeFromBranches(game.RootMove, tree)
		s.WriteString(tree.Print())

		return s.String()
	}

	s.WriteString(fmt.Sprint(game.Board) + "\n")

	if !game.Over {
		s.WriteString(fmt.Sprintf("Black Captures: %v\n", game.Captures["Black"]))
		s.WriteString(fmt.Sprintf("White Captures: %v\n", game.Captures["White"]))
		s.WriteString(fmt.Sprintf(
			"%s's turn. Enter coordintates (`x,y`), `pass`, `undo`, `resign`, `branches`: %s",
			game.Player(),
			game.Input,
		))

		return s.String()
	}

	if game.Resigned {
		s.WriteString(fmt.Sprintf("%s has resigned.\n", game.Player()))
		s.WriteString(fmt.Sprintln("The game is over."))

		return s.String()
	}

	if game.Passed {
		s.WriteString(fmt.Sprintf("%s has passed.\n", game.Player()))

		if game.Over {
			s.WriteString(fmt.Sprintln("The game is over."))
		}

		return s.String()
	}

	return s.String()
}

func Start() {
	var boardSize int

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

		boardSize = size
		break
	}

	game := NewGame(boardSize)

	if _, err := tea.NewProgram(game).Run(); err != nil {
		fmt.Printf("could not start program: %s\n", err)
		os.Exit(1)
	}
}

func (game Game) Play() (Game, bool) {
	var stone Stone
	var capturedCount int
	var captured []StoneGroup

	game.ErrMsg = ""

	if game.BlackTurn {
		stone = Black
	} else {
		stone = White
	}

	move := Coord{-1, -1}

	input := strings.TrimSpace(game.Input)

	if input == "resign" {
		game.Resigned = true
		game.Over = true
		game.LastMove = &Move{
			Number:    game.LastMove.Number + 1,
			Board:     *game.Board.Dup(),
			Passed:    true,
			BlackTurn: !game.BlackTurn,
			Move:      move,
			KoStone:   Coord{-1, -1},
			Captures:  game.Captures,
			Resigned:  true,
			HeadMove:  game.LastMove,
		}
		appendAsBranch(game.LastMove)

		return game, false
	}

	if input == "pass" {
		if game.Passed {
			game.Over = true
		}

		game.Passed = true
		game.BlackTurn = !game.BlackTurn
		game.LastMove = &Move{
			Number:    game.LastMove.Number + 1,
			Board:     *game.Board.Dup(),
			Passed:    true,
			BlackTurn: !game.BlackTurn,
			Move:      move,
			KoStone:   Coord{-1, -1},
			Captures:  game.Captures,
			HeadMove:  game.LastMove,
		}
		appendAsBranch(game.LastMove)

		return game, !game.Over
	}

	if input == "undo" {
		if game.LastMove.RootMove {
			game.ErrMsg = "Already at the beginning of the game."

			return game, true
		}

		game.LastMove = game.LastMove.HeadMove
		game.BlackTurn = !game.BlackTurn
		game.Board = game.LastMove.Board.Dup()
		game.Passed = game.LastMove.Passed
		game.KoStoneCoord = game.LastMove.KoStone
		game.Captures = game.LastMove.Captures

		return game, true
	}

	if input == "branches" {
		game.ShowBranches = true

		return game, true
	}

	game.Passed = false

	coords := strings.Split(input, ",")
	if len(coords) != 2 {
		game.ErrMsg = fmt.Sprintf("Invalid input: `%s`\n", input)

		return game, true
	}

	x, xerr := strconv.Atoi(coords[0])
	y, yerr := strconv.Atoi(coords[1])
	if xerr != nil || yerr != nil {
		game.ErrMsg = fmt.Sprintf("Invalid input: `%s`\n", input)

		return game, true
	}
	move = Coord{x, y}

	err := game.Board.Place(stone, x, y)
	if err != nil {
		game.ErrMsg = fmt.Sprintf("Cannot place %s stone at `%s`!\n", game.Player(), input)

		return game, true
	}
	groups, _ := game.Board.Groups(x, y)
	placedGroup := groups[0]
	otherGroups := groups[1:]

	capturedCount = 0
	captured = make([]StoneGroup, 0, len(otherGroups))
	for _, sg := range otherGroups {
		if len(sg.Liberties) == 0 {
			captured = append(captured, sg)
			capturedCount += len(sg.Coords)
		}
	}

	// Check suicide rule and revert if fails
	if capturedCount == 0 && len(placedGroup.Liberties) == 0 {
		err := game.Board.Remove(x, y)
		if err != nil {
			panic("This cannot happen: tried to remove a stone that doesn't exist")
		}

		game.ErrMsg = fmt.Sprintf("Cannot play suicide at `%s`!\n", input)

		return game, true
	}

	// Check ko rule and revert if fails. Also keep track of ko stones
	if capturedCount == 1 {
		if move == game.KoStoneCoord {
			err := game.Board.Remove(x, y)
			if err != nil {
				panic("This cannot happen: tried to remove a stone that doesn't exist")
			}

			game.ErrMsg = fmt.Sprintf("Cannot play ko at `%s`!\n", input)
			return game, true
		}

		game.KoStoneCoord = captured[0].Coords[0]
	} else {
		game.KoStoneCoord = Coord{-1, -1}
	}

	// Commit the captures, since there will be no more reverts
	game.Captures[game.Player()] += capturedCount

	for _, sg := range captured {
		for _, sc := range sg.Coords {
			err = game.Board.Remove(sc[0], sc[1])
			if err != nil {
				panic("This cannot happen: tried to remove a stone that doesn't exist")
			}
		}
	}

	game.BlackTurn = !game.BlackTurn
	game.LastMove = &Move{
		Number:    game.LastMove.Number + 1,
		Board:     *game.Board.Dup(),
		BlackTurn: !game.BlackTurn,
		Move:      move,
		KoStone:   game.KoStoneCoord,
		Captures:  game.Captures,
		HeadMove:  game.LastMove,
	}
	appendAsBranch(game.LastMove)

	return game, true
}
