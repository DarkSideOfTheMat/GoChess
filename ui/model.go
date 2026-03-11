// Package ui handles the TUI settings for the gochess app
package ui

import (
	"fmt"
	"strings"

	game "gochess/game"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type programState int

const (
	initializing programState = iota
	ready
)

// model is the main model which handles other submodels
// of the BubbleTea TUI
//
// It is responsible for handling user input, screen state
// and passing commands between other models.
//
// It also handles screen size / resizing
type model struct {
	game       *game.Game
	boardModel boardModel
	settings   *TUISettings
	err        error
	moveInput  textinput.Model
	width      int
	height     int
	state      programState
}

func LoadTeaModelFromFen(fen game.FENCode) model {
	// initialize text input model
	ti := textinput.New()
	ti.CharLimit = 20
	ti.Prompt = "Move: Type a move like 'e2 e4'"
	ti.Placeholder = "Type a move like 'e2 e4'"
	ti.Focus()

	// initialize board model
	board := boardModel.New()
	return model{
		game:       game.NewGame(),
		boardModel: board,
		settings:   &DEFAULT_SETTINGS,
		moveInput:  ti,
	}
}

// Model initializing, this will take care of any initial I/O
// during startup.
func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "q", "ctrl+z":
			return m, tea.Quit
		case "enter":
			// clear the text buffer and check the command
			// TODO implement...
			playerInput := m.moveInput.Value()
			// validate the input ...
			// for now it will just be a square code... in the future it should be better
			move := strings.Split(playerInput, " ")
			if len(move) != 2 {
				m.moveInput.Placeholder = fmt.Sprintf("Invalid move! %s", playerInput)
			}

			m.moveInput.Reset()
			return m, m.moveInput.Focus()
		default:
			// Player is typing
			var cmd tea.Cmd
			m.moveInput, cmd = m.moveInput.Update(msg)
			return m, cmd

		}
	// Handle the screen resizing and set main perportions
	// ...
	case tea.WindowSizeMsg:
		m.height, m.width = msg.Height, msg.Height
		// adjust component width and heights

	}
	var cmd tea.Cmd
	m.moveInput, cmd = m.moveInput.Update(msg)
	return m, cmd
}

// Main View Function for BubbleTea
func (m model) View() tea.View {
	header := headerStyle.Render("==== TEST GAME, TYPE ctrl+c or q to quit! ====")
	//board, err := m.formatCurrentBoard()
	//if err != nil {
	//	board = fmt.Sprintf("Error rendering board: %s", err)
	//}
	//

	board := m.boardModel.Render()
	footer := footerStyle.Render("Move: " + m.moveInput.View())
	return tea.NewView(lipgloss.JoinVertical(lipgloss.Left, header, board, footer))
}
