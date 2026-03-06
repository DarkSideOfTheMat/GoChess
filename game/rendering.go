package game

import "fmt"

type BoardRenderer interface {
	Render(b *Board) error
}

type TUIRenderer struct{}

// hollow symbols for white pieces, filled for black — shape alone distinguishes sides.
var pieceSymbols = map[Piece]string{
	WHITE | KING: "♔", WHITE | QUEEN: "♕", WHITE | ROOK: "♖",
	WHITE | BISHOP: "♗", WHITE | KNIGHT: "♘", WHITE | PAWN: "♙",
	BLACK | KING: "♚", BLACK | QUEEN: "♛", BLACK | ROOK: "♜",
	BLACK | BISHOP: "♝", BLACK | KNIGHT: "♞", BLACK | PAWN: "♟",
}

func (r TUIRenderer) Render(b *Board) {
	bgDark := "\x1b[48;5;130m"  // brown — dark squares
	bgLight := "\x1b[48;5;223m" // tan — light squares
	fg := "\x1b[38;5;232m"      // near-black — all pieces
	reset := "\x1b[0m"

	for rank := 7; rank >= 0; rank-- {
		fmt.Printf("%d ", rank+1)
		for file := range 8 {
			p := b.State[rank*8+file]

			bg := bgDark
			if (rank+file)%2 != 0 {
				bg = bgLight
			}

			if sym, ok := pieceSymbols[p]; ok {
				fmt.Printf("%s%s %s %s", bg, fg, sym, reset)
			} else {
				fmt.Printf("%s   %s", bg, reset)
			}
		}
		fmt.Printf("\n")
	}
	fmt.Println("  a b c d e f g h")
}
