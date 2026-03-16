// Package protocol defines the shared messages and events
// sent between the client and the server
package protocol

import (
	"gochess/chess"
	game "gochess/game"
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

// A DrawOfferMessage is sent when the player wants to request a draw
type DrawOfferMessage struct{}

// GameStateEvent is the current game state according to the Engine
type GameStateEvent struct {
	Board      game.Board
	Clock      game.Clock
	LastMove   *chess.Ply
	LegalMoves []chess.Ply
	Status     chess.GameStatus
}

// ErrorEvent is sent in the case of an engine error
type ErrorEvent struct {
	Message string
	err     error
}
