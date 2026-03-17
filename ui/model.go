// Package ui handles the TUI settings for the gochess app
package ui

import (
	"fmt"
	"strings"

	"gochess/chess"
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
	errorMsg   string
	moveInput  textinput.Model
	width      int
	height     int
	state      programState
}

// NewModel creates a new TUI model connected to the given game session.
// The session is created and owned by the caller.
func NewModel(sess *session.GameSession) model {
	settings := &DEFAULT_SETTINGS

	ti := textinput.New()
	ti.CharLimit = 100
	ti.Prompt = "Type a Move: "
	ti.Styles().Focused.Placeholder.Width(30)
	ti.Placeholder = "Input a move like 'Qd1 h5' or 'e2 e4'"
	ti.Focus()

	// Get initial board state from the session
	initialBoard := sess.GetBoard()
	bm := newBoardModel(initialBoard.State, &settings.board)

	// Subscribe now so the channel is available in Init() and Update()
	// (Init() is a value receiver, so field assignments there are lost)
	eventsCh := sess.Subscribe()

	return model{
		session:    sess,
		eventsCh:   eventsCh,
		boardModel: bm,
		settings:   settings,
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
		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			playerInput := strings.TrimSpace(m.moveInput.Value())
			m.errorMsg = ""

			switch strings.ToLower(playerInput) {
			case "quit", "exit":
				return m, tea.Quit
			case "resign":
				m.moveInput.Reset()
				return m, tea.Batch(
					m.moveInput.Focus(),
					sendResign(m.session, m.session.GetBoard().ActiveColor),
				)
			}

			move := strings.Split(playerInput, " ")
			if len(move) != 2 {
				m.errorMsg = fmt.Sprintf("Invalid move: %s", playerInput)
				m.moveInput.Reset()
				return m, m.moveInput.Focus()
			}

			from, err := chess.ParseSquare(strings.TrimSpace(move[0]))
			if err != nil {
				m.errorMsg = fmt.Sprintf("Invalid square: %s", move[0])
				m.moveInput.Reset()
				return m, m.moveInput.Focus()
			}
			to, err := chess.ParseSquare(strings.TrimSpace(move[1]))
			if err != nil {
				m.errorMsg = fmt.Sprintf("Invalid square: %s", move[1])
				m.moveInput.Reset()
				return m, m.moveInput.Focus()
			}

			m.moveInput.Reset()
			return m, tea.Batch(
				m.moveInput.Focus(),
				sendMove(m.session, from, to),
			)
		default:
			var cmd tea.Cmd
			m.moveInput, cmd = m.moveInput.Update(msg)
			return m, cmd
		}

	// Handle game state updates from the session
	case protocol.GameStateEvent:
		m.boardModel.state = msg.Board.State
		return m, listenForGameEvents(m.eventsCh)

	case protocol.ErrorEvent:
		m.errorMsg = msg.Message
		return m, nil

	// Handle the screen resizing and set main proportions
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		// Header: 1 line + padding
		headerHeight := 1
		// Footer: 1 line for input
		footerHeight := 1
		// Error line if present
		errorHeight := 0
		if m.errorMsg != "" {
			errorHeight = 1
		}

		// Remaining space goes to the board
		boardHeight := m.height - headerHeight - footerHeight - errorHeight
		if boardHeight < 1 {
			boardHeight = 1
		}

		headerStyle.Width(m.width)
		m.moveInput.SetWidth(m.width)

		_, cmd := m.boardModel.Update(tea.WindowSizeMsg{Width: m.width, Height: boardHeight})

		return m, cmd
	}
	var cmd tea.Cmd
	m.moveInput, cmd = m.moveInput.Update(msg)
	return m, cmd
}

// Main View Function for BubbleTea
func (m model) View() tea.View {
	width := m.width
	height := m.height
	if width == 0 || height == 0 {
		return tea.NewView("")
	}

	// Header: top-center
	header := headerStyle.Width(width).Render("==== TEST GAME — type 'quit' to exit ====")
	headerHeight := lipgloss.Height(header)

	// Footer: bottom-left
	footer := m.moveInput.View()
	footerHeight := lipgloss.Height(footer)

	// Error display above footer
	var errorDisplay string
	errorHeight := 0
	if m.errorMsg != "" {
		errorDisplay = errorStyle.Render(m.errorMsg)
		errorHeight = lipgloss.Height(errorDisplay)
	}

	// Board: centered in remaining space
	boardAreaHeight := height - headerHeight - footerHeight - errorHeight
	if boardAreaHeight < 1 {
		boardAreaHeight = 1
	}

	boardContent := m.boardModel.Render()
	board := lipgloss.Place(width, boardAreaHeight, lipgloss.Center, lipgloss.Center, boardContent)

	// Assemble: header (top-center), board (centered), error, footer (bottom-left)
	parts := []string{header, board}
	if errorDisplay != "" {
		parts = append(parts, errorDisplay)
	}
	parts = append(parts, footer)

	s := lipgloss.JoinVertical(lipgloss.Left, parts...)
	return tea.NewView(s)
}
