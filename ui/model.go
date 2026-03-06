package ui

import (
	"fmt"
	game "gochess/game"

	tea "charm.land/bubbletea/v2"
)

type model struct {
	game     *game.Game
	settings *TUISettings
	err      error
	cursor   tea.Cursor
}

func LoadTeaModelFromFen(fen game.FENCode) model {
	return model{
		// board
		game:     game.NewGame(),
		settings: &DEFAULT_SETTINGS,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		// exit the game
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) View() tea.View {
	// Header
	// s := "TEST BOARD, TYPE ctrl+c or q to quit!"
	s, err := m.formatCurrentBoard()
	if err != nil {
		m.err = err
		return tea.NewView(fmt.Sprintf("Opps! An error occured! Press q to quit! %s", err))
	}
	return tea.NewView(s)
}
