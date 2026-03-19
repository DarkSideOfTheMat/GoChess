package ui

import (
	"gochess/chess"
	protocol "gochess/game/protocol"
	session "gochess/game/session"

	tea "charm.land/bubbletea/v2"
)

// listenForGameEvents returns a Bubble Tea command that blocks
// on the session's event channel and returns the next GameEvent.
// The concrete type (GameStateEvent or ErrorEvent) is returned as
// a tea.Msg so model.Update can switch on it directly.
func listenForGameEvents(ch <-chan protocol.GameEvent) tea.Cmd {
	return func() tea.Msg {
		return <-ch
	}
}

// sendMove sends a move to the session. Outcomes (success or illegal-move error)
// are delivered through the event stream and picked up by listenForGameEvents.
func sendMove(sess *session.GameSession, from chess.Square, to chess.Square) tea.Cmd {
	return func() tea.Msg {
		sess.Send(protocol.MoveMessage{From: from, To: to})
		return nil
	}
}

// sendResign sends a resignation to the session.
func sendResign(sess *session.GameSession, player chess.Color) tea.Cmd {
	return func() tea.Msg {
		sess.Send(protocol.ResignMessage{ResigningPlayer: player})
		return nil
	}
}
