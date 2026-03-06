package main

import (
	"fmt"
	"os"

	game "gochess/game"
	ui "gochess/ui"

	tea "charm.land/bubbletea/v2"
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
	// demoBoard("Starting Position",
	// 	"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
	// demoBoard("Mid-game",
	// 	"1B6/2n5/p1N1P2R/P1K3N1/4Pk2/1Q2p2p/6nP/1B4R1 w - - 0 1")

	// scanner := bufio.NewScanner(os.Stdin)
	// fmt.Println("\nGoChess - Type 'quit' to exit")
	// fmt.Print("> ")

	// for scanner.Scan() {
	// 	input := strings.TrimSpace(scanner.Text())
	// 	if input == "quit" {
	// 		fmt.Println("Goodbye!")
	// 		return
	// 	}
	// 	fmt.Printf("You entered: %s\n", input)
	// 	fmt.Print("> ")
	// }

	// For now we'll hardcode Bubble Tea
	p := tea.NewProgram(ui.LoadTeaModelFromFen("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"))

	if _, err := p.Run(); err != nil {
		fmt.Printf("Cannot start bubble tea! %v", err)
		os.Exit(1)
	}
}
