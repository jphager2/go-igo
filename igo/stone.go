package igo

import "fmt"

type Stone interface {
	Format(fmt.State, rune)
}

type LibertyStone string
type BlackStone string
type WhiteStone string
type MarkStone string

const (
	Liberty = LibertyStone("L")
	Black   = BlackStone("B")
	White   = WhiteStone("W")
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

func (s MarkStone) Format(f fmt.State, c rune) {
	f.Write([]byte(s))
}
