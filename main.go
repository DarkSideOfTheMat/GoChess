package main

import (
	"fmt"
	"os"

	session "gochess/game/session"
	ui "gochess/ui"

	tea "charm.land/bubbletea/v2"
)

func main() {
	sess := session.NewGameSession()

	p := tea.NewProgram(ui.NewModel(&sess))

	if _, err := p.Run(); err != nil {
		fmt.Printf("Cannot start bubble tea! %v", err)
		os.Exit(1)
	}
}
