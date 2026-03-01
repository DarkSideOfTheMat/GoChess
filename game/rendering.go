package game

import "fmt"

type BoardRenderer interface {
	Render(b *Board) error
}

type TUIRenderer struct{}

func (r TUIRenderer) Render(b *Board) {

	black_square := "\x1b[1;37;40m"
	white_square := "\x1b[1;30;47m"
	reset := "\x1b[0m"
	for rank := 7; rank >= 0; rank-- {
		fmt.Printf("%d ", rank+1)
		for file := range 8 {
			p := b.State[rank*8+file]

			color := black_square
			if (rank+file)%2 != 0 {
				color = white_square
			}

			if sym, ok := pieceSymbols[p]; ok {
				fmt.Printf("%s %s %s", color, sym, reset)
			} else {
				fmt.Printf("%s   %s", color, reset)
			}
		}
		fmt.Printf("\n")
	}
	fmt.Println("  a b c d e f g h")
}
