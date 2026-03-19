// Package session handles the current game connection
// and processes messages
package session

import (
	"gochess/chess"
	game "gochess/game"
	protocol "gochess/game/protocol"
)

type Command interface{}

type Session interface {
	Send(cmd Command) error
	State() protocol.SessionStatus
	Subscribe() <-chan protocol.GameStateEvent
}

type GameSession struct {
	game          *game.Game
	turn          chess.Color
	status        protocol.SessionStatus
	eventsStream  chan protocol.GameEvent
	pendingDraw   bool
	drawOfferedBy chess.Color
}

func NewGameSession() GameSession {
	return GameSession{
		game:         game.NewGame(),
		turn:         chess.WHITE,
		status:       protocol.SessionInitializing,
		eventsStream: make(chan protocol.GameEvent, 1),
	}
}

func NewGameSessionFromFen(fen game.FENCode) GameSession {
	return GameSession{
		game:         game.LoadGameFromFen(fen),
		turn:         chess.WHITE,
		status:       protocol.SessionInitializing,
		eventsStream: make(chan protocol.GameEvent, 1),
	}
}

// GetBoard returns the current board state
func (gs *GameSession) GetBoard() chess.Board {
	return *gs.game.Board
}

func (gs *GameSession) Send(cmd Command) error {
	switch cmd := cmd.(type) {
	case protocol.MoveMessage:
		gs.invalidateDrawOffer()
		gs.eventsStream <- gs.HandleMove(cmd)
	case protocol.DrawOfferMessage:
		gs.handleDrawOffer(cmd)
		gs.eventsStream <- gs.getGameStateEvent()
	case protocol.DrawResponseMessage:
		gs.handleDrawResponse(cmd)
		gs.eventsStream <- gs.getGameStateEvent()
	case protocol.ResignMessage:
		gs.invalidateDrawOffer()
		gs.eventsStream <- gs.HandleResignMsg(cmd)
	default:
		gs.invalidateDrawOffer()
	}
	return nil
}

func (gs *GameSession) Subscribe() <-chan protocol.GameEvent {
	if gs.status == protocol.SessionInitializing {
		gs.status = protocol.SessionInProgress
		// Send initial game state to the subscriber
		go func() {
			gs.eventsStream <- gs.getGameStateEvent()
		}()
	}
	return gs.eventsStream
}

// HandleMove is the main interaction point during a chess game
func (gs *GameSession) HandleMove(moveMsg protocol.MoveMessage) protocol.GameEvent {
	err := gs.game.MakeMove(moveMsg.From, moveMsg.To, moveMsg.Promo)
	if err != nil {
		return protocol.ErrorEvent{Message: err.Error(), Err: err}
	}
	return gs.getGameStateEvent()
}

// HandleResignMsg will process a resign move and return an ended game.
func (gs *GameSession) HandleResignMsg(resignMsg protocol.ResignMessage) protocol.GameEvent {
	err := gs.game.PlayerResigns(resignMsg.ResigningPlayer)
	if err != nil {
		return protocol.ErrorEvent{Message: err.Error(), Err: err}
	}
	return gs.getGameStateEvent()
}

func (gs *GameSession) handleDrawOffer(msg protocol.DrawOfferMessage) {
	gs.pendingDraw = true
	gs.drawOfferedBy = msg.OfferingPlayer
}

func (gs *GameSession) handleDrawResponse(msg protocol.DrawResponseMessage) {
	if !gs.pendingDraw {
		return
	}
	if msg.Accept {
		gs.game.Status = chess.Draw
	}
	gs.invalidateDrawOffer()
}

func (gs *GameSession) invalidateDrawOffer() {
	gs.pendingDraw = false
}

func (gs *GameSession) getGameStateEvent() protocol.GameStateEvent {
	return protocol.GameStateEvent{
		Board:      *gs.game.Board,
		LastMove:   gs.game.GetLastMove(),
		LegalMoves: gs.game.GetLegalMoves(),
		Status:     gs.game.Status,
	}
}
