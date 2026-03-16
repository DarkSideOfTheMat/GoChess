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
	// State() GameStateEvent
	State() protocol.SessionStatus
	Subscribe() <-chan protocol.GameStateEvent
}

type GameSession struct {
	game         *game.Game
	turn         chess.Color
	status       protocol.SessionStatus
	eventsStream chan protocol.GameStateEvent
}

func NewGameSession() GameSession {
	return GameSession{
		game:         game.NewGame(),
		turn:         chess.WHITE,
		status:       protocol.SessionInitializing,
		eventsStream: make(chan protocol.GameStateEvent, 1),
	}
}

func NewGameSessionFromFen(fen game.FENCode) GameSession {
	return GameSession{
		game:         game.LoadGameFromFen(fen),
		turn:         chess.WHITE,
		status:       protocol.SessionInitializing,
		eventsStream: make(chan protocol.GameStateEvent, 1),
	}
}

func (gs *GameSession) Send(cmd Command) error {
	var event protocol.GameStateEvent
	switch cmd := cmd.(type) {
	case protocol.MoveMessage:
		event = gs.HandleMove(cmd)
	// case protocol.DrawOfferMessage:
	// event = gs.HandleDrawOffer(cmd)
	// case protocol.DrawResponseMessage:
	// event = gs.HandelDrawResponse(cmd)
	case protocol.ResignMessage:
		event = gs.HandleResignMsg(cmd)
	}
	gs.eventsStream <- event
	return nil
}

func (gs *GameSession) Subscribe() <-chan protocol.GameStateEvent {
	if gs.status == protocol.SessionInitializing {
		gs.status = protocol.SessionInProgress
		// Send initial game state to the subscriber
		go func() {
			gs.eventsStream <- gs.getGameStateEvent()
		}()
	}
	return gs.eventsStream
}

// GameSession.HandleMove is the main interaction point during a chess game
func (gs *GameSession) HandleMove(moveMsg protocol.MoveMessage) protocol.GameStateEvent {
	// do some work to update the board
	from := moveMsg.From
	to := moveMsg.To
	promo := moveMsg.Promo
	gs.game.MakeMove(from, to, promo)

	return gs.getGameStateEvent()
}

// GameSession.HandleResignMsg will process a resign move and return an ended game.
func (gs *GameSession) HandleResignMsg(resignMsg protocol.ResignMessage) protocol.GameStateEvent {
	err := gs.game.PlayerResigns(resignMsg.ResigningPlayer)
	if err != nil {
		// TODO: handle this error
	}
	return gs.getGameStateEvent()
}

func (gs *GameSession) getGameStateEvent() protocol.GameStateEvent {
	return protocol.GameStateEvent{
		Board:      *gs.game.Board,
		Clock:      *gs.game.Clock,
		LastMove:   gs.game.GetLastMove(),
		LegalMoves: gs.game.GetLegalMoves(),
		Status:     gs.game.Status,
	}
}
