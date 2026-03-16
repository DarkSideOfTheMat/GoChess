// Package ui handles the TUI settings for the gochess app
package ui

import (
	"fmt"
	"strings"

	game "gochess/game"
	protocol "gochess/game/protocol"
	session "gochess/game/session"

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
	session    *session.GameSession
	eventsCh   <-chan protocol.GameStateEvent
	boardModel boardModel
	settings   *TUISettings
	err        error
	moveInput  textinput.Model
	style      lipgloss.Style
	state      programState
}

func NewModel() model {
	// Initialize session
	sess := session.NewGameSession()

	// Set the default settings
	// TODO make this an IO step of init
	settings := &DEFAULT_SETTINGS

	// initialize text input model
	ti := textinput.New()
	ti.CharLimit = 100
	ti.Prompt = "Type a Move: "

	ti.Styles().Focused.Placeholder.Width(30)
	ti.Placeholder = "Input a move like 'Qd1 h5' or 'e2 e4'"
	ti.Focus()

	// initialize board model with starting position
	initialBoard := game.LoadFromFEN(game.FENCode("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"))
	board := NewBoardModel(initialBoard, &settings.board)

	// Subscribe now so the channel is available in Init() and Update()
	// (Init() is a value receiver, so field assignments there are lost)
	eventsCh := sess.Subscribe()

	return model{
		session:    &sess,
		eventsCh:   eventsCh,
		boardModel: board,
		settings:   &DEFAULT_SETTINGS,
		moveInput:  ti,
	}
}

func LoadTeaModelFromFen(fen game.FENCode) model {
	sess := session.NewGameSessionFromFen(fen)

	settings := &DEFAULT_SETTINGS

	ti := textinput.New()
	ti.CharLimit = 100
	ti.Prompt = "Type a Move: "
	ti.Styles().Focused.Placeholder.Width(30)
	ti.Placeholder = "Input a move like 'Qd1 h5' or 'e2 e4'"
	ti.Focus()

	// initialize board model from the FEN position
	initialBoard := game.LoadFromFEN(fen)
	board := NewBoardModel(initialBoard, &settings.board)

	eventsCh := sess.Subscribe()

	return model{
		session:    &sess,
		eventsCh:   eventsCh,
		boardModel: board,
		settings:   &DEFAULT_SETTINGS,
		moveInput:  ti,
	}
}

// Model initializing, this will take care of any initial I/O
// during startup.
func (m model) Init() tea.Cmd {
	return tea.Batch(
		tea.RequestWindowSize,
		listenForGameEvents(m.eventsCh),
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "q", "ctrl+z":
			return m, tea.Quit
		case "enter":
			playerInput := m.moveInput.Value()
			move := strings.Split(playerInput, " ")
			if len(move) != 2 {
				m.moveInput.Placeholder = fmt.Sprintf("Invalid move! %s", playerInput)
				m.moveInput.Reset()
				return m, m.moveInput.Focus()
			}

			from, err := parseSquare(strings.TrimSpace(move[0]))
			if err != nil {
				m.moveInput.Placeholder = fmt.Sprintf("Invalid square: %s", move[0])
				m.moveInput.Reset()
				return m, m.moveInput.Focus()
			}
			to, err := parseSquare(strings.TrimSpace(move[1]))
			if err != nil {
				m.moveInput.Placeholder = fmt.Sprintf("Invalid square: %s", move[1])
				m.moveInput.Reset()
				return m, m.moveInput.Focus()
			}

			m.moveInput.Reset()
			return m, tea.Batch(
				m.moveInput.Focus(),
				sendMove(m.session, from, to),
			)
		default:
			// Player is typing
			var cmd tea.Cmd
			m.moveInput, cmd = m.moveInput.Update(msg)
			return m, cmd
		}

	// Handle game state updates from the session
	case protocol.GameStateEvent:
		m.boardModel.board = msg.Board
		// Continue listening for the next event
		return m, listenForGameEvents(m.eventsCh)

	case protocol.ErrorEvent:
		m.moveInput.Placeholder = msg.Message
		return m, nil

	// Handle the screen resizing and set main proportions
	case tea.WindowSizeMsg:
		m.style.Height(msg.Height)
		m.style.Width(msg.Width)

		height, width := m.style.GetHeight(), m.style.GetWidth()

		// find the best board height to make the board model a square
		headerHeight := headerStyle.GetHeight()
		footerHeight := m.moveInput.Styles().Focused.Prompt.GetHeight()

		// moveInput settings

		// handle the prompt settings

		m.moveInput.SetWidth(width)

		promptWidth := len(m.moveInput.Prompt)
		placeholderMaxWidth := len(m.moveInput.Placeholder)

		m.moveInput.Styles().Focused.Prompt.Width(promptWidth)
		m.moveInput.Styles().Focused.Placeholder.MaxWidth(placeholderMaxWidth)
		m.moveInput.Styles().Focused.Placeholder.Width(width - promptWidth)

		boardHeight := height - headerHeight - footerHeight
		boardWidth := width

		m.moveInput.SetWidth(width)
		headerStyle.Width(width)

		// adjust component width and heights
		_, cmd := m.boardModel.Update(tea.WindowSizeMsg{Width: boardWidth, Height: boardHeight})

		return m, cmd
	}
	var cmd tea.Cmd
	m.moveInput, cmd = m.moveInput.Update(msg)
	return m, cmd
}

// Main View Function for BubbleTea
func (m model) View() tea.View {
	header := headerStyle.Render("==== TEST GAME, TYPE ctrl+c or q to quit! ====")

	board := m.boardModel.Render()
	footer := m.moveInput.View()

	s := lipgloss.JoinVertical(lipgloss.Center, header, board)
	s = lipgloss.JoinVertical(lipgloss.Left, s, footer)
	return tea.NewView(s)
}
