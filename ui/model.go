package ui

import (
	game "gochess/game"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
)

type model struct {
	game      *game.Game
	settings  *TUISettings
	err       error
	moveInput textinput.Model
}

func LoadTeaModelFromFen(fen game.FENCode) model {
	ti := textinput.New()
	ti.Placeholder = "e2e4"
	ti.Focus()
	return model{
		game:      game.NewGame(),
		settings:  &DEFAULT_SETTINGS,
		moveInput: ti,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		}
	}
	var cmd tea.Cmd
	m.moveInput, cmd = m.moveInput.Update(msg)
	return m, cmd
}
