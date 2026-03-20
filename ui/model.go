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
	session        *session.GameSession
	eventsCh       <-chan protocol.GameEvent
	boardModel     boardModel
	settings       *TUISettings
	err            error
	errorMsg       string
	moveInput      textinput.Model
	width          int
	height         int
	state          programState
	selectedSquare *chess.Square
	legalMoves     []chess.Ply
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

	case tea.MouseClickMsg:
		if msg.Button == tea.MouseLeft {
			if sq, ok := m.squareFromClick(msg.X, msg.Y); ok {
				if m.selectedSquare == nil {
					m.selectedSquare = &sq
					m.boardModel.selectedSquare = m.selectedSquare
					m.boardModel.legalMoveSquares = legalTargetsFrom(sq, m.legalMoves)
				} else {
					from := *m.selectedSquare
					m.selectedSquare = nil
					m.boardModel.selectedSquare = nil
					m.boardModel.legalMoveSquares = make(map[chess.Square]bool)
					m.errorMsg = ""
					return m, tea.Batch(
						m.moveInput.Focus(),
						sendMove(m.session, from, sq),
					)
				}
			}
		}
		return m, nil

	// Handle game state updates from the session
	case protocol.GameStateEvent:
		m.boardModel.state = msg.Board.State
		m.legalMoves = msg.LegalMoves
		return m, listenForGameEvents(m.eventsCh)

	case protocol.ErrorEvent:
		m.errorMsg = msg.Message
		return m, listenForGameEvents(m.eventsCh)

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

// legalTargetsFrom returns a set of destination squares for all legal moves
// starting from the given square.
func legalTargetsFrom(from chess.Square, moves []chess.Ply) map[chess.Square]bool {
	targets := make(map[chess.Square]bool)
	for _, m := range moves {
		if m.StartIdx == from {
			targets[m.EndIdx] = true
		}
	}
	return targets
}

// squareFromClick maps a terminal coordinate to a board square.
// Returns the square and true if the click landed on the board grid,
// or (0, false) if it was outside.
//
// Board layout (fixed dimensions):
//   boardTotalWidth  = 46  (42 content + 2 border + 2 margin)
//   boardTotalHeight = 13  ( 9 content + 2 border + 2 margin)
//   gridLeft offset  = +4  (1 margin + 1 border + 2 rank label)
//   gridTop  offset  = +2  (1 margin + 1 border)
func (m model) squareFromClick(x, y int) (chess.Square, bool) {
	errorHeight := 0
	if m.errorMsg != "" {
		errorHeight = 1
	}
	boardAreaHeight := m.height - 1 - 1 - errorHeight // 1 header + 1 footer
	const boardTotalWidth, boardTotalHeight = 46, 13
	boardTotalLeft := (m.width - boardTotalWidth) / 2
	boardTotalTop := 1 + (boardAreaHeight-boardTotalHeight)/2 // 1 = header row
	gridLeft := boardTotalLeft + 4
	gridTop := boardTotalTop + 2
	relX := x - gridLeft
	relY := y - gridTop
	if relX < 0 || relX >= 40 || relY < 0 || relY >= 8 {
		return 0, false
	}
	file := relX / 5
	rank := 7 - relY
	return chess.Square(rank*8 + file), true
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
	v := tea.NewView(s)
	v.MouseMode = tea.MouseModeCellMotion
	return v
}
