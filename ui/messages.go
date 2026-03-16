package ui

import (
	"fmt"

	"gochess/chess"
	protocol "gochess/game/protocol"
	session "gochess/game/session"

	tea "charm.land/bubbletea/v2"
)

// listenForGameEvents returns a Bubble Tea command that blocks
// on the session's event channel and returns the next GameStateEvent
func listenForGameEvents(ch <-chan protocol.GameStateEvent) tea.Cmd {
	return func() tea.Msg {
		return <-ch
	}
}

// sendMove sends a move to the session and returns a command
// that reports any error back as a tea.Msg
func sendMove(sess *session.GameSession, from chess.Square, to chess.Square) tea.Cmd {
	return func() tea.Msg {
		msg := protocol.MoveMessage{From: from, To: to}
		err := sess.Send(msg)
		if err != nil {
			return protocol.ErrorEvent{Message: fmt.Sprintf("move error: %v", err)}
		}
		return nil
	}
}

// parseSquare converts a string like "e2" to a chess.Square index
func parseSquare(s string) (chess.Square, error) {
	if len(s) != 2 {
		return 0, fmt.Errorf("invalid square: %s", s)
	}
	file := int(s[0] - 'a')
	rank := int(s[1] - '1')
	if file < 0 || file > 7 || rank < 0 || rank > 7 {
		return 0, fmt.Errorf("invalid square: %s", s)
	}
	return chess.Square(rank*8 + file), nil
}
