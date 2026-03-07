package main

import (
	"fmt"
	"os"

	ui "gochess/ui"

	tea "charm.land/bubbletea/v2"
)

func main() {

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
