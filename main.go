package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"gochess/game"
)

func demoBoard(label string, fen game.FENCode) {
	renderer := game.TUIRenderer{}

	fmt.Printf("\n=== %s ===\n", label)
	fmt.Printf("FEN: %s\n\n", fen)
	board := game.LoadFromFEN(fen)
	renderer.Render(&board)
	if board.Validate() {
		fmt.Println("Board: valid")
	} else {
		fmt.Println("Board: INVALID")
	}
}

func main() {
	demoBoard("Starting Position",
		"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
	demoBoard("Mid-game",
		"1B6/2n5/p1N1P2R/P1K3N1/4Pk2/1Q2p2p/6nP/1B4R1 w - - 0 1")

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("\nGoChess - Type 'quit' to exit")
	fmt.Print("> ")

	for scanner.Scan() {
		input := strings.TrimSpace(scanner.Text())
		if input == "quit" {
			fmt.Println("Goodbye!")
			return
		}
		fmt.Printf("You entered: %s\n", input)
		fmt.Print("> ")
	}
}
