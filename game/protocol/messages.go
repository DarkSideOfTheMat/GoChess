// Package protocol defines the shared messages and events
// sent between the client and the server
package protocol

import (
	"time"

	"gochess/chess"
)

type SessionStatus uint

const (
	SessionUnknown SessionStatus = iota
	SessionInitializing
	SessionDisconnected
	SessionInProgress
	SessionEnded
)

// MoveMessage a move to send to the chess engine
type MoveMessage struct {
	From  chess.Square
	To    chess.Square
	Promo chess.Piece
}

// ResignMessage is a message to end the game by resignation
type ResignMessage struct {
	ResigningPlayer chess.Color
}

// DrawOfferMessage is sent when a player wants to request a draw
type DrawOfferMessage struct {
	OfferingPlayer chess.Color
}

// DrawResponseMessage is a response to a draw offer
type DrawResponseMessage struct {
	AcceptingPlayer chess.Color
	Accept          bool
}

// GameStateEvent is the current game state according to the Engine
type GameStateEvent struct {
	Board      chess.Board
	WhiteTime  time.Duration
	BlackTime  time.Duration
	LastMove   *chess.Ply
	LegalMoves []chess.Ply
	Status     chess.GameStatus
}

// ErrorEvent is sent in the case of an engine error
type ErrorEvent struct {
	Message string
	Err     error
}

// GameEvent is implemented by all events sent from the session to the UI.
// This allows BubbleTea to process both success and error outcomes as typed messages.
type GameEvent interface {
	isGameEvent()
}

func (GameStateEvent) isGameEvent() {}
func (ErrorEvent) isGameEvent()     {}
